package ast

type Expr interface{}

type Variable struct {
	Name string
}

type Number struct {
	Value float64
}

type String struct {
	Value string
}

type Binary struct {
	Op    string
	Left  Expr
	Right Expr
}

type Stmt interface{}

type Assignment struct {
	Name string
	Expr Expr
}

var _ Stmt = (*Assignment)(nil)

type PrintStmt struct {
	Expr Expr
}

type Program struct {
	Stmts []Stmt
}
