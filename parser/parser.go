package parser

import (
	"fmt"
	"strconv"

	"github.com/engpetarmarinov/pede/ast"
	"github.com/engpetarmarinov/pede/lexer"
)

type Parser struct {
	lx  *lexer.Lexer
	cur lexer.Token
}

func NewParser(lx *lexer.Lexer) *Parser {
	p := &Parser{lx: lx}
	p.next()
	return p
}

func (p *Parser) next() {
	p.cur = p.lx.Next()
}

// Parse parses a program (sequence of statements)
func (p *Parser) Parse() *ast.Program {
	stmts := []ast.Stmt{}
	for p.cur.Type != lexer.TokenEOF {
		// Skip empty/unknown tokens (e.g., newlines if present)
		if p.cur.Type == lexer.TokenUnknown && (p.cur.Value == "\n" || p.cur.Value == ";") {
			p.next()
			continue
		}
		stmts = append(stmts, p.parseStmt())
	}
	return &ast.Program{Stmts: stmts}
}

// parseStmt parses a single statement
func (p *Parser) parseStmt() ast.Stmt {
	if p.cur.Type == lexer.TokenPrint {
		return p.parsePrint()
	}
	if p.cur.Type == lexer.TokenIdent {
		name := p.cur.Value
		p.next()
		if p.cur.Type != lexer.TokenEqual {
			panic("expected '='")
		}
		p.next()
		expr := p.parseExpr()
		return &ast.Assignment{Name: name, Expr: expr}
	}
	panic(fmt.Sprintf("unexpected token: %v", p.cur))
}

// parsePrint parses a print statement: print(expr)
func (p *Parser) parsePrint() ast.Stmt {
	p.next() // consume 'print'
	if p.cur.Type != lexer.TokenUnknown || p.cur.Value != "(" {
		panic("expected '('")
	}
	p.next() // consume '('
	expr := p.parseExpr()
	if p.cur.Type != lexer.TokenUnknown || p.cur.Value != ")" {
		panic("expected ')'")
	}
	p.next() // consume ')'
	return &ast.PrintStmt{Expr: expr}
}

func (p *Parser) parseExpr() ast.Expr {
	left := p.parseTerm()
	for p.cur.Type == lexer.TokenPlus {
		op := p.cur.Value
		p.next()
		right := p.parseTerm()
		left = &ast.Binary{Op: op, Left: left, Right: right}
	}

	return left
}

func (p *Parser) parseTerm() ast.Expr {
	left := p.parseFactor()
	for p.cur.Type == lexer.TokenStar {
		op := p.cur.Value
		p.next()
		right := p.parseFactor()
		left = &ast.Binary{Op: op, Left: left, Right: right}
	}

	return left
}

func (p *Parser) parseFactor() ast.Expr {
	switch p.cur.Type {
	case lexer.TokenNumber:
		val, _ := strconv.ParseFloat(p.cur.Value, 64)
		p.next()
		return &ast.Number{Value: val}
	case lexer.TokenIdent:
		name := p.cur.Value
		p.next()
		return &ast.Variable{Name: name}
	case lexer.TokenString:
		str := p.cur.Value
		p.next()
		return &ast.String{Value: str}
	default:
		panic(fmt.Sprintf("unexpected token: %v", p.cur))
	}
}
