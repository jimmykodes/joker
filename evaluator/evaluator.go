package evaluator

import (
	"github.com/jimmykodes/jk/ast"
	"github.com/jimmykodes/jk/object"
)

var (
	True  = &object.Boolean{Value: true}
	False = &object.Boolean{Value: false}
	Null  = &object.Null{}
)

func Eval(node ast.Node) object.Object {
	switch n := node.(type) {
	case *ast.Program:
		return evalProgram(n)
	case *ast.BlockStatement:
		return evalBlockStatements(n)
	case *ast.ReturnStatement:
		return &object.Return{Value: Eval(n.Value)}
	case *ast.ExpressionStatement:
		return Eval(n.Expression)
	case *ast.PrefixExpression:
		return evalPrefix(n.Operator, Eval(n.Right))
	case *ast.InfixExpression:
		return evalInfix(n.Operator, Eval(n.Left), Eval(n.Right))
	case *ast.IfExpression:
		return evalIf(n)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: n.Value}
	case *ast.FloatLiteral:
		return &object.Float{Value: n.Value}
	case *ast.BooleanLiteral:
		return toBoolObject(n.Value)
	case *ast.StringLiteral:
		return &object.String{Value: n.Value}
	}
	return nil
}

func toBoolObject(b bool) object.Object {
	if b {
		return True
	}
	return False
}

func intToFloat(obj object.Object) object.Object {
	return &object.Float{Value: float64(obj.(*object.Integer).Value)}
}

func evalProgram(program *ast.Program) object.Object {
	var res object.Object
	for _, stmt := range program.Statements {
		res = Eval(stmt)
		if r, ok := res.(*object.Return); ok {
			return r.Value
		}
	}
	return res
}

func evalBlockStatements(block *ast.BlockStatement) object.Object {
	var res object.Object
	for _, statement := range block.Statements {
		res = Eval(statement)
		if res.Type() == object.ReturnType {
			return res
		}
	}
	return res
}

func evalPrefix(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBang(right)
	case "-":
		return evalNegative(right)
	default:
		return Null
	}
}

func evalInfix(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.IntegerType && right.Type() == object.IntegerType:
		return evalIntInfix(operator, left, right)
	case left.Type() == object.IntegerType && right.Type() == object.FloatType:
		return evalFloatInfix(operator, intToFloat(left), right)
	case left.Type() == object.FloatType && right.Type() == object.IntegerType:
		return evalFloatInfix(operator, left, intToFloat(right))
	case left.Type() == object.FloatType && right.Type() == object.FloatType:
		return evalFloatInfix(operator, left, right)
	}
	return Null
}

func evalIf(n *ast.IfExpression) object.Object {
	switch Eval(n.Condition) {
	case True:
		return Eval(n.Consequence)
	case False:
		if n.Alternative != nil {
			return Eval(n.Alternative)
		}
		fallthrough
	default:
		return Null
	}
}

func evalBang(right object.Object) object.Object {
	switch right {
	case True:
		return False
	case False:
		return True
	case Null:
		return True
	default:
		return False
	}
}

func evalNegative(right object.Object) object.Object {
	switch r := right.(type) {
	case *object.Integer:
		return &object.Integer{Value: r.Value}
	case *object.Float:
		return &object.Float{Value: r.Value}
	default:
		return Null
	}
}

func evalIntInfix(operator string, left, right object.Object) object.Object {
	lv := left.(*object.Integer).Value
	rv := right.(*object.Integer).Value
	switch operator {
	case "+":
		return &object.Integer{Value: lv + rv}
	case "-":
		return &object.Integer{Value: lv - rv}
	case "*":
		return &object.Integer{Value: lv * rv}
	case "/":
		return &object.Integer{Value: lv / rv}
	case "%":
		return &object.Integer{Value: lv % rv}
	case "<":
		return toBoolObject(lv < rv)
	case ">":
		return toBoolObject(lv > rv)
	case "<=":
		return toBoolObject(lv <= rv)
	case ">=":
		return toBoolObject(lv >= rv)
	case "==":
		return toBoolObject(lv == rv)
	case "!=":
		return toBoolObject(lv != rv)
	}
	return Null
}

func evalFloatInfix(operator string, left, right object.Object) object.Object {
	lv := left.(*object.Float).Value
	rv := right.(*object.Float).Value
	switch operator {
	case "+":
		return &object.Float{Value: lv + rv}
	case "-":
		return &object.Float{Value: lv - rv}
	case "*":
		return &object.Float{Value: lv * rv}
	case "/":
		return &object.Float{Value: lv / rv}
	case "<":
		return toBoolObject(lv < rv)
	case ">":
		return toBoolObject(lv > rv)
	case "<=":
		return toBoolObject(lv <= rv)
	case ">=":
		return toBoolObject(lv >= rv)
	case "==":
		return toBoolObject(lv == rv)
	case "!=":
		return toBoolObject(lv != rv)
	}
	return Null
}
