package semantic

import (
	"strings"
	"testing"

	"gopl/internal/lexer"
	"gopl/internal/parser"
)

func checkSource(t *testing.T, source string) error {
	t.Helper()
	program, err := parser.New(lexer.New(strings.NewReader(source))).Parse()
	if err != nil {
		return err
	}
	return NewSemanticChecker().Check(program)
}

func TestCheckValidProgram(t *testing.T) {
	err := checkSource(t, `int add(int a, int b) { return a + b }
void main() { print(add(1, 2)) }`)
	if err != nil {
		t.Fatalf("valid program rejected: %v", err)
	}
}

func TestCheckRequiresMain(t *testing.T) {
	err := checkSource(t, "int add(int a) { return a }")
	if err == nil || !strings.Contains(err.Error(), "missing main") {
		t.Fatalf("error = %v, want missing main", err)
	}
}

func TestCheckRejectsTypeMismatch(t *testing.T) {
	err := checkSource(t, `void main() { int value = true }`)
	if err == nil || !strings.Contains(err.Error(), "type mismatch") {
		t.Fatalf("error = %v, want type mismatch", err)
	}
}
