// This file is autogenerated from the genstd tool (@/misc/stdgen); do not edit.
// To regenerate it, run `go generate` from @/stdlibs.

package stdlibs

import (
	"reflect"

	gno "github.com/gnolang/gno/gnovm/pkg/gnolang"
	lib0 "github.com/gnolang/gno/gnovm/tests/stdlibs/std"
)

type nativeFunc struct {
	gnoPkg  string
	gnoFunc gno.Name
	params  []gno.FieldTypeExpr
	results []gno.FieldTypeExpr
	f       func(m *gno.Machine)
}

var nativeFuncs = [...]nativeFunc{
	{
		"std",
		"AssertOriginCall",

		[]gno.FieldTypeExpr{},
		[]gno.FieldTypeExpr{},
		func(m *gno.Machine) {
			lib0.AssertOriginCall(
				m,
			)
		},
	},
	{
		"std",
		"IsOriginCall",

		[]gno.FieldTypeExpr{},
		[]gno.FieldTypeExpr{
			{Name: gno.N("r0"), Type: gno.X("bool")},
		},
		func(m *gno.Machine) {
			r0 := lib0.IsOriginCall(
				m,
			)

			m.PushValue(gno.Go2GnoValue(
				m.Alloc,
				m.Store,
				reflect.ValueOf(&r0).Elem(),
			))
		},
	},
	{
		"std",
		"TestCurrentRealm",

		[]gno.FieldTypeExpr{},
		[]gno.FieldTypeExpr{
			{Name: gno.N("r0"), Type: gno.X("string")},
		},
		func(m *gno.Machine) {
			r0 := lib0.TestCurrentRealm(
				m,
			)

			m.PushValue(gno.Go2GnoValue(
				m.Alloc,
				m.Store,
				reflect.ValueOf(&r0).Elem(),
			))
		},
	},
	{
		"std",
		"TestSkipHeights",

		[]gno.FieldTypeExpr{
			{Name: gno.N("p0"), Type: gno.X("int64")},
		},
		[]gno.FieldTypeExpr{},
		func(m *gno.Machine) {
			b := m.LastBlock()
			var (
				p0  int64
				rp0 = reflect.ValueOf(&p0).Elem()
			)

			gno.Gno2GoValue(b.GetPointerTo(nil, gno.NewValuePathBlock(1, 0, "")).TV, rp0)

			lib0.TestSkipHeights(p0)
		},
	},
	{
		"std",
		"ClearStoreCache",

		[]gno.FieldTypeExpr{},
		[]gno.FieldTypeExpr{},
		func(m *gno.Machine) {
			lib0.ClearStoreCache(
				m,
			)
		},
	},
}
