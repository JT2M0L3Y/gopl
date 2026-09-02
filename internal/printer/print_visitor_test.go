package printer

import (
	"bytes"
	"strings"
	"testing"

	"gopl/internal/lexer"
	"gopl/internal/parser"
)

func TestPrintVisitorProducesProgram(t *testing.T) {
	program, err := parser.New(lexer.New(strings.NewReader(`void main() { print(2 + 3) }`))).Parse()
	if err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	if err := program.Accept(NewPrintVisitor(&output)); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), "void main()") || !strings.Contains(output.String(), "print(2 + 3)") {
		t.Fatalf("unexpected output: %q", output.String())
	}
}
