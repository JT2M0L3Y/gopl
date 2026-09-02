package token

import "fmt"

var SingleCharTokens = map[byte]Kind{
	'.': Dot,
	',': Comma,
	'(': LParen,
	')': RParen,
	'[': LBracket,
	']': RBracket,
	';': Semicolon,
	'{': LBrace,
	'}': RBrace,
	'+': Plus,
	'-': Minus,
	'*': Times,
	'/': Divide,
	'=': Assign,
	'<': Less,
	'>': Greater,
}

var TwoCharOperators = map[string]Kind{
	"<=": LessEq,
	">=": GreaterEq,
	"==": Equal,
	"!=": NotEqual,
}

var KeywordTokens = map[string]Kind{
	"true":   BoolVal,
	"false":  BoolVal,
	"null":   NullVal,
	"int":    IntType,
	"double": FloatType,
	"bool":   BoolType,
	"string": StringType,
	"char":   CharType,
	"void":   VoidType,
	"struct": Struct,
	"array":  Array,
	"for":    For,
	"while":  While,
	"if":     If,
	"elseif": ElseIf,
	"else":   Else,
	"and":    And,
	"or":     Or,
	"not":    Not,
	"new":    New,
	"return": Return,
}

type Kind string

const (
	Illegal Kind = "illegal"
	EOF     Kind = "eof"
	Ident   Kind = "ident"

	Dot       Kind = "dot"
	Comma     Kind = "comma"
	LParen    Kind = "lparen"
	RParen    Kind = "rparen"
	LBracket  Kind = "lbracket"
	RBracket  Kind = "rbracket"
	Semicolon Kind = "semicolon"
	LBrace    Kind = "lbrace"
	RBrace    Kind = "rbrace"
	Plus      Kind = "plus"
	Minus     Kind = "minus"
	Times     Kind = "times"
	Divide    Kind = "divide"
	Assign    Kind = "assign"
	Less      Kind = "less"
	Greater   Kind = "greater"
	LessEq    Kind = "less_eq"
	GreaterEq Kind = "greater_eq"
	Equal     Kind = "equal"
	NotEqual  Kind = "not_equal"

	IntVal    Kind = "int_val"
	FloatVal  Kind = "float_val"
	CharVal   Kind = "char_val"
	StringVal Kind = "string_val"
	BoolVal   Kind = "bool_val"
	NullVal   Kind = "null_val"

	IntType    Kind = "int_type"
	FloatType  Kind = "float_type"
	BoolType   Kind = "bool_type"
	StringType Kind = "string_type"
	CharType   Kind = "char_type"
	VoidType   Kind = "void_type"

	Struct Kind = "struct"
	Array  Kind = "array"
	For    Kind = "for"
	While  Kind = "while"
	If     Kind = "if"
	ElseIf Kind = "elseif"
	Else   Kind = "else"
	And    Kind = "and"
	Or     Kind = "or"
	Not    Kind = "not"
	New    Kind = "new"
	Return Kind = "return"
)

// Token captures a lexeme and its source location.
type Token struct {
	Kind   Kind
	Lexeme string
	Line   int
	Col    int
}

func (t Token) String() string {
	return fmt.Sprintf("%s(%q @ %d:%d)", t.Kind, t.Lexeme, t.Line, t.Col)
}
