package evaluator

import (
	"fmt"

	"github.com/jimmykodes/jk/ast"
	"github.com/jimmykodes/jk/object"
)

var (
	True  = &object.Boolean{Value: true}
	False = &object.Boolean{Value: false}
	Null  = &object.Null{}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch n := node.(type) {
	case *ast.Program:
		return evalProgram(n, env)
	case *ast.LetStatement:
		r := Eval(n.Value, env)
		if isError(r) {
			return r
		}
		env.Set(n.Name.Value, r)
	case *ast.BlockStatement:
		return evalBlockStatements(n, env)
	case *ast.ReturnStatement:
		r := Eval(n.Value, env)
		if isError(r) {
			return r
		}
		return &object.Return{Value: r}
	case *ast.ExpressionStatement:
		return Eval(n.Expression, env)
	case *ast.PrefixExpression:
		r := Eval(n.Right, env)
		if isError(r) {
			return r
		}
		return evalPrefix(n.Operator, r)
	case *ast.InfixExpression:
		l := Eval(n.Left, env)
		if isError(l) {
			return l
		}
		r := Eval(n.Right, env)
		if isError(r) {
			return r
		}
		return evalInfix(n.Operator, l, r)
	case *ast.IfExpression:
		return evalIf(n, env)
	case *ast.CallExpression:
		f := Eval(n.Function, env)
		if isError(f) {
			return f
		}
		args, err := evalExpressions(n.Arguments, env)
		if isError(err) {
			return err
		}
		return applyFunc(f, args)
	case *ast.Identifier:
		o, ok := env.Get(n.Value)
		if !ok {
			return newError("identifier not found: %s", n.Value)
		}
		return o
	case *ast.FunctionLiteral:
		return &object.Function{Parameters: n.Parameters, Body: n.Body, Env: env}
	case *ast.IntegerLiteral:
		return &object.Integer{Value: n.Value}
	case *ast.FloatLiteral:
		return &object.Float{Value: n.Value}
	case *ast.BooleanLiteral:
		return toBoolObject(n.Value)
	case *ast.StringLiteral:
		return &object.String{Value: n.Value}
	default:
		return newError("invalid node type: %T", node)
	}
	return Null
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

func newError(format string, a ...any) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(o object.Object) bool {
	return o != nil && o.Type() == object.ErrorType
}

func applyFunc(fn object.Object, args []object.Object) object.Object {
	f, ok := fn.(*object.Function)
	if !ok {
		return newError("cannot call a non-function: %s", fn.Type())
	}
	if len(args) != len(f.Parameters) {
		return newError("invalid number of args. got %d - want %d", len(args), len(f.Parameters))
	}
	wrappedEnv := object.NewEnvironment(object.EncloseOuterOption(f.Env))
	for i, parameter := range f.Parameters {
		wrappedEnv.Set(parameter.Value, args[i])
	}
	ret := Eval(f.Body, wrappedEnv)
	if r, ok := ret.(*object.Return); ok {
		return r.Value
	}
	return ret
}

func evalExpressions(exp []ast.Expression, env *object.Environment) ([]object.Object, object.Object) {
	res := make([]object.Object, len(exp))
	for i, expression := range exp {
		r := Eval(expression, env)
		if isError(r) {
			return nil, r
		}
		res[i] = r
	}
	return res, Null
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var res object.Object
	for _, stmt := range program.Statements {
		res = Eval(stmt, env)
		switch r := res.(type) {
		case *object.Return:
			return r.Value
		case *object.Error:
			return r
		}
	}
	return res
}

func evalBlockStatements(block *ast.BlockStatement, env *object.Environment) object.Object {
	var res object.Object
	for _, statement := range block.Statements {
		res = Eval(statement, env)
		if res.Type() == object.ReturnType || res.Type() == object.ErrorType {
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
		return newError("unknown operator %s%s", operator, right.Type())
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
	return newError("unknown operator %s %s %s", left.Type(), operator, right.Type())
}

func evalIf(n *ast.IfExpression, env *object.Environment) object.Object {
	switch Eval(n.Condition, env) {
	case True:
		return Eval(n.Consequence, env)
	case False:
		if n.Alternative != nil {
			return Eval(n.Alternative, env)
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
	default:
		return newError("invalid type for !: %s", right.Type())
	}
}

func evalNegative(right object.Object) object.Object {
	switch r := right.(type) {
	case *object.Integer:
		return &object.Integer{Value: -r.Value}
	case *object.Float:
		return &object.Float{Value: -r.Value}
	default:
		return newError("invalid type for !: %s", right.Type())
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
	return newError("unknown operator %s %s %s", left.Type(), operator, right.Type())
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
	return newError("unknown operator %s %s %s", left.Type(), operator, right.Type())
}
