package lexer

import (
	"strings"
	"testing"

	"gopl/internal/token"
)

func TestNextTokenizesProgram(t *testing.T) {
	lexer := New(strings.NewReader("int main() { return 42 }"))
	want := []token.Kind{token.IntType, token.Ident, token.LParen, token.RParen, token.LBrace, token.Return, token.IntVal, token.RBrace, token.EOF}
	for i, expected := range want {
		if got := lexer.Next().Kind; got != expected {
			t.Fatalf("token %d = %q, want %q", i, got, expected)
		}
	}
}

func TestStringEscapes(t *testing.T) {
	lexer := New(strings.NewReader(`"line\nnext"`))
	tok := lexer.Next()
	if tok.Kind != token.StringVal || tok.Lexeme != "line\nnext" {
		t.Fatalf("got %#v, want string with newline escape", tok)
	}
}

func TestIllegalCharacter(t *testing.T) {
	tok := New(strings.NewReader("@")).Next()
	if tok.Kind != token.Illegal {
		t.Fatalf("kind = %q, want illegal", tok.Kind)
	}
}
