package lexer

import (
	"bufio"
	"fmt"
	"io"
	"unicode"

	"gopl/internal/token"
)

// Lexer converts input text into tokens.
type Lexer struct {
	r      *bufio.Reader
	line   int
	column int
}

// New creates a lexer for the provided reader.
func New(reader io.Reader) *Lexer {
	return &Lexer{r: bufio.NewReader(reader), line: 1, column: 0}
}

// Next returns the next token from the input stream.
func (l *Lexer) Next() token.Token {
	l.skipWhitespaceAndComments()

	startLine, startCol := l.line, l.column+1
	ch, ok := l.peek()
	if !ok {
		return token.Token{Kind: token.EOF, Lexeme: "end-of-stream", Line: startLine, Col: startCol}
	}

	if kind, found := token.SingleCharTokens[ch]; found {
		l.read()
		return l.makeToken(kind, string(ch), startLine, startCol)
	}

	if ch == '<' || ch == '>' || ch == '=' || ch == '!' {
		return l.readOperator(ch, startLine, startCol)
	}

	switch ch {
	case '\'':
		return l.readChar(startLine, startCol)
	case '"':
		return l.readString(startLine, startCol)
	}

	if unicode.IsDigit(rune(ch)) {
		return l.readNumber(startLine, startCol)
	}
	if unicode.IsLetter(rune(ch)) || ch == '_' {
		return l.readWord(startLine, startCol)
	}
	l.read()
	return l.illegal(fmt.Sprintf("unexpected character %q", ch), startLine, startCol)
}

func (l *Lexer) skipWhitespaceAndComments() {
	for {
		ch, ok := l.peek()
		if !ok {
			return
		}
		if unicode.IsSpace(rune(ch)) {
			l.read()
			continue
		}
		if ch == '#' {
			for {
				ch, ok = l.read()
				if !ok || ch == '\n' {
					break
				}
			}
			continue
		}
		return
	}
}

func (l *Lexer) readOperator(ch byte, startLine int, startCol int) token.Token {
	l.read()
	if next, ok := l.peek(); ok {
		twoChar := string([]byte{ch, next})
		if kind, found := token.TwoCharOperators[twoChar]; found {
			l.read()
			return l.makeToken(kind, twoChar, startLine, startCol)
		}
	}

	switch ch {
	case '<':
		return l.makeToken(token.Less, string(ch), startLine, startCol)
	case '>':
		return l.makeToken(token.Greater, string(ch), startLine, startCol)
	case '=':
		return l.makeToken(token.Assign, string(ch), startLine, startCol)
	case '!':
		return l.illegal("unexpected '!'", startLine, startCol)
	default:
		return l.illegal(fmt.Sprintf("unexpected operator %q", ch), startLine, startCol)
	}
}

func (l *Lexer) readWord(startLine, startCol int) token.Token {
	var lexeme []byte
	for {
		ch, ok := l.peek()
		if !ok || !isIdentPart(ch) {
			break
		}
		l.read()
		lexeme = append(lexeme, ch)
	}
	word := string(lexeme)
	if kind, found := token.KeywordTokens[word]; found {
		return l.makeToken(kind, word, startLine, startCol)
	}
	return l.makeToken(token.Ident, word, startLine, startCol)
}

func (l *Lexer) readNumber(startLine, startCol int) token.Token {
	var lexeme []byte
	isFloat := false
	for {
		ch, ok := l.peek()
		if !ok {
			break
		}
		if unicode.IsDigit(rune(ch)) {
			l.read()
			lexeme = append(lexeme, ch)
			continue
		}
		if ch == '.' {
			if isFloat {
				break
			}
			isFloat = true
			l.read()
			lexeme = append(lexeme, ch)
			continue
		}
		break
	}
	value := string(lexeme)
	if isFloat {
		return l.makeToken(token.FloatVal, value, startLine, startCol)
	}
	return l.makeToken(token.IntVal, value, startLine, startCol)
}

func (l *Lexer) readString(startLine, startCol int) token.Token {
	l.read() // opening quote
	var lexeme []byte
	for {
		ch, ok := l.read()
		if !ok {
			return l.illegal("found end-of-file in string", startLine, startCol)
		}
		if ch == '\n' {
			return l.illegal("found end-of-line in string", startLine, startCol)
		}
		if ch == '"' {
			return l.makeToken(token.StringVal, string(lexeme), startLine, startCol)
		}
		if ch == '\\' {
			next, ok := l.read()
			if !ok {
				return l.illegal("found end-of-file in string", startLine, startCol)
			}
			switch next {
			case 'n':
				lexeme = append(lexeme, '\n')
			case 't':
				lexeme = append(lexeme, '\t')
			case '\\':
				lexeme = append(lexeme, '\\')
			case '"':
				lexeme = append(lexeme, '"')
			default:
				lexeme = append(lexeme, '\\', next)
			}
			continue
		}
		lexeme = append(lexeme, ch)
	}
}

func (l *Lexer) readChar(startLine, startCol int) token.Token {
	l.read() // opening quote
	ch, ok := l.read()
	if !ok {
		return l.illegal("found end-of-file in character", startLine, startCol)
	}
	if ch == '\n' || ch == '\t' {
		return l.illegal("found invalid character literal", startLine, startCol)
	}
	if ch == '\\' {
		next, ok := l.read()
		if !ok {
			return l.illegal("found end-of-file in character", startLine, startCol)
		}
		ch = translateEscape(next)
	}
	end, ok := l.read()
	if !ok || end != '\'' {
		return l.illegal("expecting closing quote for character", startLine, startCol)
	}
	return l.makeToken(token.CharVal, string(ch), startLine, startCol)
}

func translateEscape(ch byte) byte {
	switch ch {
	case 'n':
		return '\n'
	case 't':
		return '\t'
	case '\\':
		return '\\'
	case '\'':
		return '\''
	default:
		return ch
	}
}

func (l *Lexer) makeToken(kind token.Kind, lexeme string, line, col int) token.Token {
	return token.Token{Kind: kind, Lexeme: lexeme, Line: line, Col: col}
}

func (l *Lexer) illegal(msg string, line, col int) token.Token {
	return token.Token{Kind: token.Illegal, Lexeme: msg, Line: line, Col: col}
}

func (l *Lexer) read() (byte, bool) {
	ch, err := l.r.ReadByte()
	if err != nil {
		return 0, false
	}
	if ch == '\n' {
		l.line++
		l.column = 0
	} else {
		l.column++
	}
	return ch, true
}

func (l *Lexer) peek() (byte, bool) {
	b, err := l.r.Peek(1)
	if err != nil || len(b) == 0 {
		return 0, false
	}
	return b[0], true
}

func isIdentPart(ch byte) bool {
	return unicode.IsLetter(rune(ch)) || unicode.IsDigit(rune(ch)) || ch == '_'
}
