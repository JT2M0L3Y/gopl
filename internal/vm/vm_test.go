package vm

import (
	"bytes"
	"strings"
	"testing"
)

func TestExecuteArithmeticAndOutput(t *testing.T) {
	var output bytes.Buffer
	runtime := New()
	runtime.SetOutput(&output)
	runtime.Add(&FrameInfo{Name: "main", Instructions: []Instr{
		NewInstr(PUSH, 2), NewInstr(PUSH, 3), NewInstr(ADD, nil), NewInstr(WRITE, nil),
		NewInstr(PUSH, nil), NewInstr(RET, nil),
	}})
	if err := runtime.Execute(false); err != nil {
		t.Fatal(err)
	}
	if output.String() != "5" {
		t.Fatalf("output = %q, want %q", output.String(), "5")
	}
}

func TestExecuteFunctionCall(t *testing.T) {
	var output bytes.Buffer
	runtime := New()
	runtime.SetOutput(&output)
	runtime.Add(&FrameInfo{Name: "add", ArgCount: 2, Instructions: []Instr{
		NewInstr(STORE, 0), NewInstr(STORE, 1), NewInstr(LOAD, 0), NewInstr(LOAD, 1), NewInstr(ADD, nil), NewInstr(RET, nil),
	}})
	runtime.Add(&FrameInfo{Name: "main", Instructions: []Instr{
		NewInstr(PUSH, 4), NewInstr(PUSH, 6), NewInstr(CALL, "add"), NewInstr(WRITE, nil), NewInstr(PUSH, nil), NewInstr(RET, nil),
	}})
	if err := runtime.Execute(false); err != nil {
		t.Fatal(err)
	}
	if output.String() != "10" {
		t.Fatalf("output = %q, want %q", output.String(), "10")
	}
}

func TestExecuteReadsInput(t *testing.T) {
	var output bytes.Buffer
	runtime := New()
	runtime.SetInput(strings.NewReader("hello\n"))
	runtime.SetOutput(&output)
	runtime.Add(&FrameInfo{Name: "main", Instructions: []Instr{
		NewInstr(READ, nil), NewInstr(WRITE, nil), NewInstr(PUSH, nil), NewInstr(RET, nil),
	}})
	if err := runtime.Execute(false); err != nil {
		t.Fatal(err)
	}
	if output.String() != "hello" {
		t.Fatalf("output = %q, want %q", output.String(), "hello")
	}
}
