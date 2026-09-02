package parser

import (
	"strings"
	"testing"

	"gopl/internal/ast"
	"gopl/internal/lexer"
)

func parseTestProgram(t *testing.T, source string) *ast.Program {
	t.Helper()
	program, err := New(lexer.New(strings.NewReader(source))).Parse()
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	return program
}

func TestParseProgramStructure(t *testing.T) {
	program := parseTestProgram(t, `struct Point { int x, int y }
int add(int a, int b) { return a + b }
void main() { print(add(1, 2)) }`)
	if len(program.StructDefs) != 1 || len(program.FunDefs) != 2 {
		t.Fatalf("got %d structs and %d functions", len(program.StructDefs), len(program.FunDefs))
	}
	if program.FunDefs[0].FunName.Lexeme != "add" || len(program.FunDefs[0].Params) != 2 {
		t.Fatalf("unexpected function: %#v", program.FunDefs[0])
	}
}

func TestParseReportsSyntaxError(t *testing.T) {
	_, err := New(lexer.New(strings.NewReader("void main( {"))).Parse()
	if err == nil {
		t.Fatal("expected syntax error")
	}
}
