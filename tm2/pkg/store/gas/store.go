package gas

import (
	"github.com/gnolang/gno/telemetry"
	"github.com/gnolang/gno/telemetry/traces"
	"github.com/gnolang/gno/tm2/pkg/store/types"
	"github.com/gnolang/overflow"
	"go.opentelemetry.io/otel/attribute"
)

var _ types.Store = &Store{}

// Store applies gas tracking to an underlying Store. It implements the
// Store interface.
type Store struct {
	gasMeter  types.GasMeter
	gasConfig types.GasConfig
	parent    types.Store
}

// New returns a reference to a new GasStore.
func New(parent types.Store, gasMeter types.GasMeter, gasConfig types.GasConfig) *Store {
	kvs := &Store{
		gasMeter:  gasMeter,
		gasConfig: gasConfig,
		parent:    parent,
	}
	return kvs
}

// Implements Store.
func (gs *Store) Get(key []byte) (value []byte) {
	var gas int64
	// telemetry  start
	var span *traces.Span
	if telemetry.IsEnabled() && traces.IsTraceStore() {
		span = traces.StartSpan(
			"Store.Get",
		)
		defer func() {
			span.SetAttributes(
				attribute.Int64(types.GasReadCostFlatDesc, gs.gasConfig.ReadCostFlat),
				attribute.Int64(types.GasReadPerByteDesc, gas),
			)
			span.End()
		}()
	}
	// telemetry end

	gs.gasMeter.ConsumeGas(gs.gasConfig.ReadCostFlat, types.GasReadCostFlatDesc)
	value = gs.parent.Get(key)

	gas = overflow.Mul64p(gs.gasConfig.ReadCostPerByte, types.Gas(len(value)))
	gs.gasMeter.ConsumeGas(gas, types.GasReadPerByteDesc)

	return value
}

// Implements Store.
func (gs *Store) Set(key []byte, value []byte) {
	var gas int64
	// telemetry code start
	var span *traces.Span
	if telemetry.IsEnabled() && traces.IsTraceStore() {
		span = traces.StartSpan(
			"Store.Set",
		)
		defer func() {
			span.SetAttributes(
				attribute.Int64(types.GasWriteCostFlatDesc, gs.gasConfig.WriteCostFlat),
				attribute.Int64(types.GasWritePerByteDesc, gas),
			)
			span.End()
		}()
	}
	// telemetry code end

	types.AssertValidValue(value)
	gs.gasMeter.ConsumeGas(gs.gasConfig.WriteCostFlat, types.GasWriteCostFlatDesc)

	gas = overflow.Mul64p(gs.gasConfig.WriteCostPerByte, types.Gas(len(value)))
	gs.gasMeter.ConsumeGas(gas, types.GasWritePerByteDesc)
	gs.parent.Set(key, value)
}

// Implements Store.
func (gs *Store) Has(key []byte) bool {
	gs.gasMeter.ConsumeGas(gs.gasConfig.HasCost, types.GasHasDesc)
	return gs.parent.Has(key)
}

// Implements Store.
func (gs *Store) Delete(key []byte) {
	// telemetry  start
	var span *traces.Span
	if telemetry.IsEnabled() && traces.IsTraceStore() {
		span = traces.StartSpan(
			"Store.Delete",
		)
		defer func() {
			span.SetAttributes(
				attribute.Int64(types.GasDeleteDesc, gs.gasConfig.DeleteCost),
			)
			span.End()
		}()
	}
	// telemetry end

	// charge gas to prevent certain attack vectors even though space is being freed
	gs.gasMeter.ConsumeGas(gs.gasConfig.DeleteCost, types.GasDeleteDesc)
	gs.parent.Delete(key)
}

// Iterator implements the Store interface. It returns an iterator which
// incurs a flat gas cost for seeking to the first key/value pair and a variable
// gas cost based on the current value's length if the iterator is valid.
func (gs *Store) Iterator(start, end []byte) types.Iterator {
	return gs.iterator(start, end, true)
}

// ReverseIterator implements the Store interface. It returns a reverse
// iterator which incurs a flat gas cost for seeking to the first key/value pair
// and a variable gas cost based on the current value's length if the iterator
// is valid.
func (gs *Store) ReverseIterator(start, end []byte) types.Iterator {
	return gs.iterator(start, end, false)
}

// Implements Store.
func (gs *Store) CacheWrap() types.Store {
	panic("cannot CacheWrap a gas.Store")
}

// Implements Store.
func (gs *Store) Write() {
	gs.parent.Write()
}

func (gs *Store) iterator(start, end []byte, ascending bool) types.Iterator {
	var parent types.Iterator
	if ascending {
		parent = gs.parent.Iterator(start, end)
	} else {
		parent = gs.parent.ReverseIterator(start, end)
	}

	gi := newGasIterator(gs.gasMeter, gs.gasConfig, parent)
	if gi.Valid() {
		gi.(*gasIterator).consumeSeekGas()
	}

	return gi
}

type gasIterator struct {
	gasMeter  types.GasMeter
	gasConfig types.GasConfig
	parent    types.Iterator
}

func newGasIterator(gasMeter types.GasMeter, gasConfig types.GasConfig, parent types.Iterator) types.Iterator {
	return &gasIterator{
		gasMeter:  gasMeter,
		gasConfig: gasConfig,
		parent:    parent,
	}
}

// Implements Iterator.
func (gi *gasIterator) Domain() (start []byte, end []byte) {
	return gi.parent.Domain()
}

// Implements Iterator.
func (gi *gasIterator) Valid() bool {
	return gi.parent.Valid()
}

// Next implements the Iterator interface. It seeks to the next key/value pair
// in the iterator. It incurs a flat gas cost for seeking and a variable gas
// cost based on the current value's length if the iterator is valid.
func (gi *gasIterator) Next() {
	if gi.Valid() {
		gi.consumeSeekGas()
	}

	gi.parent.Next()
}

// Key implements the Iterator interface. It returns the current key and it does
// not incur any gas cost.
func (gi *gasIterator) Key() (key []byte) {
	key = gi.parent.Key()
	return key
}

// Value implements the Iterator interface. It returns the current value and it
// does not incur any gas cost.
func (gi *gasIterator) Value() (value []byte) {
	value = gi.parent.Value()
	return value
}

// Implements Iterator.
func (gi *gasIterator) Close() {
	gi.parent.Close()
}

// consumeSeekGas consumes a flat gas cost for seeking and a variable gas cost
// based on the current value's length.
func (gi *gasIterator) consumeSeekGas() {
	var gas int64
	// telemetry start
	var span *traces.Span
	if telemetry.IsEnabled() && traces.IsTraceStore() {
		span = traces.StartSpan(
			"Store.Seek",
		)
		defer func() {
			span.SetAttributes(
				attribute.Int64(types.GasIterNextCostFlatDesc, gi.gasConfig.IterNextCostFlat),
				attribute.Int64(types.GasValuePerByteDesc, gas),
			)
			span.End()
		}()
	}
	// telemetry  end

	value := gi.Value()
	gas = overflow.Mul64p(gi.gasConfig.ReadCostPerByte, types.Gas(len(value)))
	gi.gasMeter.ConsumeGas(gi.gasConfig.IterNextCostFlat, types.GasIterNextCostFlatDesc)
	gi.gasMeter.ConsumeGas(gas, types.GasValuePerByteDesc)
}
