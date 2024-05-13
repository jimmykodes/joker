package ast

type Expr interface {
	Node
	expr()
}

//go:generate go run ../tools/expr Binary Left:Expr Right:Expr Operator:token.Token
//go:generate go run ../tools/expr Unary Right:Expr Operator:token.Token
//go:generate go run ../tools/expr Grouping Expr:Expr Token:token.Token
//go:generate go run ../tools/expr Ident Token:token.Token Name:string
//go:generate go run ../tools/expr CommentLit Token:token.Token Value:string
//go:generate go run ../tools/expr StringLit Token:token.Token Value:string
//go:generate go run ../tools/expr IntLit Token:token.Token Value:int64
//go:generate go run ../tools/expr FloatLit Token:token.Token Value:float64
//go:generate go run ../tools/expr BoolLit Token:token.Token Value:bool
//go:generate go run ../tools/expr NilLit Token:token.Token
//go:generate go run ../tools/expr FuncLit Token:token.Token Params:[]*IdentExpr Body:*BlockStmt
