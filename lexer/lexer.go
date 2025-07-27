package lexer

import (
	"fmt"
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
	TokenComment = "//"
	TokenLParen  = "LPAREN"
	TokenRParen  = "RPAREN"
	TokenNewline = "NEWLINE"
)

type Token struct {
	Type  TokenType
	Value string
}

type Error struct {
	Message    string
	Line       int
	Column     int
	LineSource string
}

func (e *Error) Error() string {
	pointer := ""
	for i := 0; i < e.Column-1; i++ {
		pointer += " "
	}
	pointer += "^"
	return fmt.Sprintf(`%s at Line %d, column %d: 
%s
%s`, e.Message, e.Line, e.Column, e.LineSource, pointer)
}

type Lexer struct {
	input     []rune
	pos       int
	Line      int
	Col       int
	lineStart int // index of the start of the current Line
}

func NewLexer(input string) *Lexer {
	return &Lexer{input: []rune(input), Line: 1, Col: 1, lineStart: 0}
}

func (l *Lexer) CurrentLineSource() string {
	start := l.lineStart
	end := l.pos
	for end < len(l.input) && l.input[end] != '\n' {
		end++
	}
	return string(l.input[start:end])
}

// Next returns the next token, or an error if an unknown or invalid token is encountered.
func (l *Lexer) Next() (Token, error) {
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch == '\n' {
			l.Line++
			l.Col = 1
			l.lineStart = l.pos + 1
			l.pos++
			return Token{Type: TokenNewline, Value: "\n"}, nil
		}
		if unicode.IsSpace(ch) {
			l.Col++
			l.pos++
			continue
		}
		break
	}
	if l.pos >= len(l.input) {
		return Token{Type: TokenEOF}, nil
	}

	ch := l.input[l.pos]
	startCol := l.Col

	switch {
	case unicode.IsDigit(ch):
		start := l.pos
		for l.pos < len(l.input) && unicode.IsDigit(l.input[l.pos]) {
			l.pos++
			l.Col++
		}
		return Token{Type: TokenNumber, Value: string(l.input[start:l.pos])}, nil
	case unicode.IsLetter(ch):
		start := l.pos
		for l.pos < len(l.input) && unicode.IsLetter(l.input[l.pos]) {
			l.pos++
			l.Col++
		}
		word := string(l.input[start:l.pos])
		if word == "print" {
			return Token{Type: TokenPrint, Value: word}, nil
		}
		return Token{Type: TokenIdent, Value: word}, nil
	case ch == '+':
		l.pos++
		l.Col++
		return Token{Type: TokenPlus, Value: "+"}, nil
	case ch == '*':
		l.pos++
		l.Col++
		return Token{Type: TokenStar, Value: "*"}, nil
	case ch == '=':
		l.pos++
		l.Col++
		return Token{Type: TokenEqual, Value: "="}, nil
	case ch == '"':
		l.pos++ // skip opening quote
		l.Col++
		start := l.pos
		for l.pos < len(l.input) && l.input[l.pos] != '"' {
			if l.input[l.pos] == '\n' {
				break
			}
			l.pos++
			l.Col++
		}
		if l.pos >= len(l.input) || l.input[l.pos] != '"' {
			return Token{}, &Error{
				Message:    "unterminated string",
				Line:       l.Line,
				Column:     startCol,
				LineSource: l.CurrentLineSource(),
			}
		}
		str := string(l.input[start:l.pos])
		l.pos++ // skip closing quote
		l.Col++
		return Token{Type: TokenString, Value: str}, nil
	case ch == '(': // support left paren
		l.pos++
		l.Col++
		return Token{Type: TokenLParen, Value: "("}, nil
	case ch == ')': // support right paren
		l.pos++
		l.Col++
		return Token{Type: TokenRParen, Value: ")"}, nil
	default:
		unknownChar := l.input[l.pos]
		err := &Error{
			Message:    fmt.Sprintf("unknown character '%c'", unknownChar),
			Line:       l.Line,
			Column:     startCol,
			LineSource: l.CurrentLineSource(),
		}
		return Token{}, err
	}
}
