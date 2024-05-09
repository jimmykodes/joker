package eval

import (
	"fmt"

	"github.com/jimmykodes/joker/ast"
	"github.com/jimmykodes/joker/token"
)

var symbolTable = make(map[string]Object)

func Eval(node ast.Node) Object {
	switch node := node.(type) {
	case *ast.ProgramStmt:
		for _, stmt := range node.Stmts {
			obj := Eval(stmt)
			if isErr(obj) {
				fmt.Println("error:", obj)
				return obj
			}
			fmt.Println(obj)
		}
		return nil
	case *ast.LetStmt:
		v := Eval(node.Value)
		if isErr(v) {
			return v
		}
		symbolTable[node.Name] = v
		return nil
	case *ast.ExprStmt:
		return Eval(node.Expr)
	case *ast.BinaryExpr:
		left := Eval(node.Left)
		if isErr(left) {
			return left
		}
		right := Eval(node.Right)
		if isErr(right) {
			return right
		}
		switch node.Operator.Type {
		case token.Plus:
			l, ok := left.(Addable)
			if !ok {
				return Errf("%s does not implement addition", left.Type())
			}
			return l.Add(right)
		default:
			return Errf("unsupported operator: %q", node.Operator.Type)
		}
	case *ast.IntLitExpr:
		return IntObject{Value: node.Value}
	case *ast.FloatLitExpr:
		return FloatObject{Value: node.Value}
	case *ast.StringLitExpr:
		return StringObject{Value: node.Value}

	default:
		return Errf("unsupported node type: %T", node)
	}
}

func isErr(obj Object) bool {
	return obj != nil && obj.Type() == ErrType
}
