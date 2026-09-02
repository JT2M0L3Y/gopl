package vm

type FrameInfo struct {
	Name         string
	ArgCount     int
	Instructions []Instr
}

type Frame struct {
	Info      FrameInfo
	ProgCount int
	Variables []Value
	OpStack   []Value
}

func (f *Frame) New() {
	f.ProgCount = 0
	f.Variables = nil
	f.OpStack = nil
}
