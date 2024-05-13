package ast

type Stmt interface {
	Node
	stmt()
}

//go:generate go run ../tools/stmt Expr Expr:Expr
//go:generate go run ../tools/stmt Block Stmts:[]Stmt
//go:generate go run ../tools/stmt Let Name:*IdentExpr Value:Expr
