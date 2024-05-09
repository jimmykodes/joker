package ast

type Stmt interface {
	Node
	stmt()
}

//go:generate go run ../tools/stmt Expr Expr:Expr
//go:generate go run ../tools/stmt Program Stmts:[]Stmt
