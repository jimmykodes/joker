package evaluator

import (
	"errors"
	"fmt"

	"github.com/jimmykodes/jk/ast"
	"github.com/jimmykodes/jk/object"
)

var (
	Null = &object.Null{}
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
	case *ast.ReassignStatement:
		_, ok := env.Get(n.Name.Value)
		if !ok {
			return newError("cannot assign to uninitialized variable: %s", n.Name.Value)
		}
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
	case *ast.ContinueStatement:
		return &object.Continue{}
	case *ast.BreakStatement:
		return &object.Break{}
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
	case *ast.WhileExpression:
		return evalWhile(n, env)
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
	case *ast.IndexExpression:
		return evalIndex(n, env)
	case *ast.Identifier:
		return evalIdent(n, env)
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
	case *ast.ArrayLiteral:
		elems := make([]object.Object, len(n.Elements))
		for i, element := range n.Elements {
			elems[i] = Eval(element, env)
		}
		return &object.Array{Elements: elems}
	default:
		return newError("invalid node type: %T", node)
	}
	return Null
}

func toBoolObject(b bool) object.Object {
	if b {
		return object.True
	}
	return object.False
}

func newError(format string, a ...any) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(o object.Object) bool {
	return o != nil && o.Type() == object.ErrorType
}

func applyFunc(fn object.Object, args []object.Object) object.Object {
	switch f := fn.(type) {
	case *object.Builtin:
		return f.Fn(args...)
	case *object.Function:
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
	default:
		return newError("cannot call a non-function: %s", fn.Type())
	}
}

func evalIndex(index *ast.IndexExpression, env *object.Environment) object.Object {
	left := Eval(index.Left, env)
	if isError(left) {
		return left
	}
	i := Eval(index.Index, env)
	if isError(i) {
		return i
	}
	o, err := left.Idx(i)
	if errors.Is(err, object.ErrUnsupportedType) {
		return newError("cannot index %s with type %s", left.Type(), i.Type())
	} else if errors.Is(err, object.ErrUnsupportedOperation) {
		return newError("cannot index object of type %s", left.Type())
	} else if err != nil {
		return newError(err.Error())
	}
	return o
}

func evalIdent(ident *ast.Identifier, env *object.Environment) object.Object {
	if o, ok := env.Get(ident.Value); ok {
		return o
	}
	if b, ok := builtins[ident.Value]; ok {
		return b
	}
	return newError("identifier not found: %s", ident.Value)
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
		if res.Type() == object.ReturnType || res.Type() == object.ErrorType || res.Type() == object.BreakType {
			return res
		}
		if res.Type() == object.ContinueType {
			return Null
		}
	}
	return res
}

func evalPrefix(operator string, right object.Object) object.Object {
	var (
		o   object.Object
		err error
	)
	switch operator {
	case "!":
		o, err = right.Bang()
	case "-":
		o, err = right.Negative()
	default:
		return newError("unknown operator %s%s", operator, right.Type())
	}
	if errors.Is(err, object.ErrUnsupportedOperation) {
		return newError("unsupported operation (%s) on type %s", operator, right.Type())
	} else if err != nil {
		return newError(err.Error())
	}
	return o
}

func evalInfix(operator string, left, right object.Object) object.Object {
	var (
		r   object.Object
		err error
	)
	switch operator {
	case "+":
		r, err = left.Add(right)
	case "-":
		r, err = left.Minus(right)
	case "*":
		r, err = left.Mult(right)
	case "/":
		r, err = left.Div(right)
	case "%":
		r, err = left.Mod(right)
	case "<":
		r, err = left.LT(right)
	case ">":
		r, err = left.GT(right)
	case "<=":
		r, err = left.LTE(right)
	case ">=":
		r, err = left.GTE(right)
	case "==":
		r, err = left.EQ(right)
	case "!=":
		r, err = left.NEQ(right)
	default:
		return newError("unknown operator %s %s %s", left.Type(), operator, right.Type())
	}
	if errors.Is(err, object.ErrUnsupportedOperation) {
		return newError("unsupported operation (%s) on %s", operator, left.Type())
	} else if errors.Is(err, object.ErrUnsupportedType) {
		return newError("unsupported operation (%s) between %s and %s", operator, left.Type(), right.Type())
	} else if err != nil {
		return newError(err.Error())
	}
	return r

}

func evalIf(n *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(n.Condition, env)
	if isError(condition) {
		return condition
	}

	b, err := condition.Bool()
	if errors.Is(err, object.ErrUnsupportedOperation) {
		return newError("cannot implicitly convert %s to bool", condition.Type())
	} else if err != nil {
		return newError(err.Error())
	}
	if b {
		return Eval(n.Consequence, env)
	}
	if n.Alternative != nil {
		return Eval(n.Alternative, env)
	}
	return Null
}

func evalWhile(n *ast.WhileExpression, env *object.Environment) object.Object {
	var res object.Object
	for {
		condition := Eval(n.Condition, env)
		if isError(condition) {
			return condition
		}
		b, err := condition.Bool()
		if errors.Is(err, object.ErrUnsupportedOperation) {
			return newError("cannot implicitly convert %s to bool", condition.Type())
		} else if err != nil {
			return newError(err.Error())
		}
		if !b {
			return res
		}
		res = Eval(n.Body, env)
		if res.Type() == object.ReturnType || res.Type() == object.BreakType {
			return res
		}
	}
}
