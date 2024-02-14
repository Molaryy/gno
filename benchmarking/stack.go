package benchmarking

const initStackSize int = 64

var (
	measurementStack []*measurement
	stackSize        int
)

func initStack() {
	measurementStack = make([]*measurement, initStackSize)
}

func StartMeasurement(opCode OpCode) {
	if stackSize != 0 {
		measurementStack[stackSize-1].pause()
	}

	if stackSize == len(measurementStack) {
		newStack := make([]*measurement, stackSize*2)
		copy(newStack, measurementStack)
		measurementStack = newStack
	}

	measurementStack[stackSize] = startNewMeasurement(opCode)
	stackSize++
}

// Pause pauses current measurement on the stack
func Pause() {
	if stackSize != 0 {
		measurementStack[stackSize-1].pause()
	}
}

// Resume resumes current measurement on the stack
func Resume() {
	if stackSize != 0 {
		measurementStack[stackSize-1].resume()
	}
}

// StopMeasurement ends the current measurement and resumes the previous one
// if one exists. It accepts the number of bytes that were read/written to/from
// the store. This value is zero if the operation is not a read or write.
func StopMeasurement(size uint32) {
	if stackSize == 0 {
		return
	}

	measurementStack[stackSize-1].end(size)

	stackSize--

	if stackSize != 0 {
		measurementStack[stackSize-1].resume()
	}
}

func RecordAllocation(size uint32) {
	if stackSize == 0 {
		return
	}

	measurementStack[stackSize-1].allocation += size
}
