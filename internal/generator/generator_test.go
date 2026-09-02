package generator

import (
	"strings"
	"testing"

	"gopl/internal/lexer"
	"gopl/internal/parser"
	"gopl/internal/vm"
)

func TestGenerateCreatesFunctionFrames(t *testing.T) {
	program, err := parser.New(lexer.New(strings.NewReader(`int add(int a, int b) { return a + b }
void main() { print(add(1, 2)) }`))).Parse()
	if err != nil {
		t.Fatal(err)
	}
	runtime := vm.New()
	if err := New(runtime).Generate(program); err != nil {
		t.Fatal(err)
	}
	if err := runtime.Execute(false); err != nil {
		t.Fatal(err)
	}
}
