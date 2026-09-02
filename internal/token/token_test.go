package token

import "testing"

func TestKeywordTokens(t *testing.T) {
	for word, want := range map[string]Kind{"int": IntType, "void": VoidType, "while": While, "true": BoolVal} {
		if got := KeywordTokens[word]; got != want {
			t.Fatalf("keyword %q = %q, want %q", word, got, want)
		}
	}
}

func TestTokenStringIncludesLocation(t *testing.T) {
	tok := Token{Kind: Ident, Lexeme: "name", Line: 3, Col: 7}
	if got := tok.String(); got != `ident("name" @ 3:7)` {
		t.Fatalf("String() = %q", got)
	}
}
