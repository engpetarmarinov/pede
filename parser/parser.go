package parser

import (
	"fmt"
	"strconv"

	"github.com/engpetarmarinov/pede/ast"
	"github.com/engpetarmarinov/pede/lexer"
)

type Parser struct {
	lx        *lexer.Lexer
	cur       lexer.Token
	curLine   int
	curColumn int
	curSource string
}

func NewParser(lx *lexer.Lexer) *Parser {
	p := &Parser{lx: lx}
	p.next()
	return p
}

func (p *Parser) next() error {
	tok, err := p.lx.Next()
	if err != nil {
		return err
	}
	p.cur = tok
	if p.lx != nil {
		p.curLine = p.lx.Line
		p.curColumn = p.lx.Col
		p.curSource = p.lx.CurrentLineSource()
	}
	return nil
}

// Parse parses a program (sequence of statements)
func (p *Parser) Parse() (*ast.Program, error) {
	stmts := []ast.Stmt{}
	for {
		// Skip any NEWLINE tokens before parsing a statement
		for p.cur.Type == lexer.TokenNewline {
			if err := p.next(); err != nil {
				return nil, err
			}
		}
		if p.cur.Type == lexer.TokenEOF {
			break
		}
		stmt, err := p.parseStmt()
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, stmt)
	}
	return &ast.Program{Stmts: stmts}, nil
}

// parseStmt parses a single statement
func (p *Parser) parseStmt() (ast.Stmt, error) {
	// Skip over any NEWLINE tokens
	for p.cur.Type == lexer.TokenNewline {
		if err := p.next(); err != nil {
			return nil, err
		}
	}
	if p.cur.Type == lexer.TokenPrint {
		return p.parsePrint()
	}
	if p.cur.Type == lexer.TokenIdent {
		name := p.cur.Value
		if err := p.next(); err != nil {
			return nil, err
		}
		if p.cur.Type != lexer.TokenEqual {
			return nil, &lexer.Error{
				Message:    "parser: expected '=' after identifier",
				Line:       p.curLine,
				Column:     p.curColumn,
				LineSource: p.curSource,
			}
		}
		if err := p.next(); err != nil {
			return nil, err
		}
		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		return &ast.Assignment{Name: name, Expr: expr}, nil
	}
	return nil, &lexer.Error{
		Message:    fmt.Sprintf("parser: unexpected token: %v", p.cur),
		Line:       p.curLine,
		Column:     p.curColumn,
		LineSource: p.curSource,
	}
}

// parsePrint parses a print statement: print(expr)
func (p *Parser) parsePrint() (ast.Stmt, error) {
	if err := p.next(); err != nil {
		return nil, err
	}
	if p.cur.Type != lexer.TokenLParen {
		return nil, &lexer.Error{
			Message:    "parser: expected '(' after print",
			Line:       p.curLine,
			Column:     p.curColumn,
			LineSource: p.curSource,
		}
	}
	if err := p.next(); err != nil {
		return nil, err
	}
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if p.cur.Type != lexer.TokenRParen {
		return nil, &lexer.Error{
			Message:    "parser: expected ')' after print expression",
			Line:       p.curLine,
			Column:     p.curColumn,
			LineSource: p.curSource,
		}
	}
	if err := p.next(); err != nil {
		return nil, err
	}
	return &ast.PrintStmt{Expr: expr}, nil
}

func (p *Parser) parseExpr() (ast.Expr, error) {
	left, err := p.parseTerm()
	if err != nil {
		return nil, err
	}
	for p.cur.Type == lexer.TokenPlus {
		op := p.cur.Value
		if err := p.next(); err != nil {
			return nil, err
		}
		right, err := p.parseTerm()
		if err != nil {
			return nil, err
		}
		left = &ast.Binary{Op: op, Left: left, Right: right}
	}
	return left, nil
}

func (p *Parser) parseTerm() (ast.Expr, error) {
	left, err := p.parseFactor()
	if err != nil {
		return nil, err
	}
	for p.cur.Type == lexer.TokenStar {
		op := p.cur.Value
		if err := p.next(); err != nil {
			return nil, err
		}
		right, err := p.parseFactor()
		if err != nil {
			return nil, err
		}
		left = &ast.Binary{Op: op, Left: left, Right: right}
	}
	return left, nil
}

func (p *Parser) parseFactor() (ast.Expr, error) {
	switch p.cur.Type {
	case lexer.TokenNumber:
		val, _ := strconv.ParseFloat(p.cur.Value, 64)
		if err := p.next(); err != nil {
			return nil, err
		}
		return &ast.Number{Value: val}, nil
	case lexer.TokenIdent:
		name := p.cur.Value
		if err := p.next(); err != nil {
			return nil, err
		}
		return &ast.Variable{Name: name}, nil
	case lexer.TokenString:
		str := p.cur.Value
		if err := p.next(); err != nil {
			return nil, err
		}
		return &ast.String{Value: str}, nil
	default:
		return nil, &lexer.Error{
			Message:    fmt.Sprintf("parser: unexpected token in expression: %v", p.cur),
			Line:       p.curLine,
			Column:     p.curColumn,
			LineSource: p.curSource,
		}
	}
}
