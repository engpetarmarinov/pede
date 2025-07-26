package lexer

import (
	"unicode"
)

type TokenType string

const (
	TokenEOF     = "EOF"
	TokenIdent   = "IDENT"
	TokenNumber  = "NUMBER"
	TokenPlus    = "+"
	TokenStar    = "*"
	TokenEqual   = "="
	TokenPrint   = "PRINT"
	TokenString  = "STRING"
	TokenUnknown = "UNKNOWN"
)

type Token struct {
	Type  TokenType
	Value string
}

type Lexer struct {
	input []rune
	pos   int
}

func NewLexer(input string) *Lexer {
	return &Lexer{input: []rune(input)}
}

func (l *Lexer) Next() Token {
	for l.pos < len(l.input) && unicode.IsSpace(l.input[l.pos]) {
		l.pos++
	}
	if l.pos >= len(l.input) {
		return Token{Type: TokenEOF}
	}

	ch := l.input[l.pos]

	switch {
	case unicode.IsDigit(ch):
		start := l.pos
		for l.pos < len(l.input) && unicode.IsDigit(l.input[l.pos]) {
			l.pos++
		}
		return Token{Type: TokenNumber, Value: string(l.input[start:l.pos])}
	case unicode.IsLetter(ch):
		start := l.pos
		for l.pos < len(l.input) && unicode.IsLetter(l.input[l.pos]) {
			l.pos++
		}
		word := string(l.input[start:l.pos])
		if word == "print" {
			return Token{Type: TokenPrint, Value: word}
		}
		return Token{Type: TokenIdent, Value: word}
	case ch == '+':
		l.pos++
		return Token{Type: TokenPlus, Value: "+"}
	case ch == '*':
		l.pos++
		return Token{Type: TokenStar, Value: "*"}
	case ch == '=':
		l.pos++
		return Token{Type: TokenEqual, Value: "="}
	case ch == '"':
		l.pos++ // skip opening quote
		start := l.pos
		for l.pos < len(l.input) && l.input[l.pos] != '"' {
			l.pos++
		}
		if l.pos >= len(l.input) {
			return Token{Type: TokenUnknown, Value: "unterminated string"}
		}
		str := string(l.input[start:l.pos])
		l.pos++ // skip closing quote
		// Ensure we do not return the closing quote as a separate token
		return Token{Type: TokenString, Value: str}
	default:
		l.pos++
		return Token{Type: TokenUnknown, Value: string(ch)}
	}
}
