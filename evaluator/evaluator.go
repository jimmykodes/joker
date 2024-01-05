package evaluator

import (
	"fmt"

	"github.com/jimmykodes/joker/ast"
	"github.com/jimmykodes/joker/object"
)

var Null = &object.Null{}

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch n := node.(type) {
	case *ast.Program:
		return evalProgram(n, env)
	case *ast.LetStatement:
		// TODO: check if redefining already existing var
		r := Eval(n.Value, env)
		if isError(r) {
			return r
		}
		env.Set(n.Name.Value, r)
	case *ast.DefineStatement:
		// TODO: check if redefining already existing var
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
	case *ast.FuncStatement:
		f := Eval(n.Fn, env)
		if isError(f) {
			return f
		}
		env.Set(n.Name.Value, f)
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
		return applyFunc(f, args, env)
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
			elem := Eval(element, env)
			if isError(elem) {
				return elem
			}
			elems[i] = elem
		}
		return &object.Array{Elements: elems}
	case *ast.MapLiteral:
		return evalMap(n, env)
	case *ast.CommentLiteral:
		return Null
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

func applyFunc(fn object.Object, args []object.Object, env *object.Environment) object.Object {
	switch f := fn.(type) {
	case *object.Builtin:
		return f.Fn(env, args...)
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
	l, ok := left.(object.Indexer)
	if !ok {
		return object.ErrUnsupportedType
	}
	return l.Idx(i)
}

func evalMap(m *ast.MapLiteral, env *object.Environment) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)
	for k, v := range m.Pairs {
		kv := Eval(k, env)
		if isError(kv) {
			return kv
		}
		hashable, ok := kv.(object.Hashable)
		if !ok {
			return newError("cannot use %s as map key", kv.Type())
		}
		vv := Eval(v, env)
		if isError(vv) {
			return vv
		}
		pairs[hashable.HashKey()] = object.HashPair{Key: kv, Value: vv}
	}
	return &object.Map{Pairs: pairs}
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
		if isError(res) {
			return res
		}
		if res.Type() == object.ReturnType || res.Type() == object.BreakType {
			return res
		}
		if res.Type() == object.ContinueType {
			return Null
		}
	}
	return res
}

func evalPrefix(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBang(right)
	case "-":
		r, ok := right.(object.Negater)
		if !ok {
			return newError("unknown operator (%s) on %s", operator, right.Type())
		}
		return r.Negative()
	default:
		return newError("unknown operator %s%s", operator, right.Type())
	}
}

func evalInfix(operator string, left, right object.Object) object.Object {
	switch operator {
	case "+":
		l, ok := left.(object.Adder)
		if !ok {
			return newError("unsupported operation (%s) on %s", operator, left.Type())
		}
		return l.Add(right)
	case "-":
		l, ok := left.(object.Subber)
		if !ok {
			return newError("unsupported operation (%s) on %s", operator, left.Type())
		}
		return l.Sub(right)
	case "*":
		l, ok := left.(object.MultDiver)
		if !ok {
			return newError("unsupported operation (%s) on %s", operator, left.Type())
		}
		return l.Mult(right)
	case "/":
		l, ok := left.(object.MultDiver)
		if !ok {
			return newError("unsupported operation (%s) on %s", operator, left.Type())
		}
		return l.Div(right)
	case "%":
		l, ok := left.(object.Modder)
		if !ok {
			return newError("unsupported operation (%s) on %s", operator, left.Type())
		}
		return l.Mod(right)
	case "<":
		l, ok := left.(object.Inequality)
		if !ok {
			return newError("unsupported operation (%s) on %s", operator, left.Type())
		}
		return l.LT(right)
	case ">":
		l, ok := left.(object.Inequality)
		if !ok {
			return newError("unsupported operation (%s) on %s", operator, left.Type())
		}
		return l.GT(right)
	case "<=":
		l, ok := left.(object.Inequality)
		if !ok {
			return newError("unsupported operation (%s) on %s", operator, left.Type())
		}
		return l.LTE(right)
	case ">=":
		l, ok := left.(object.Inequality)
		if !ok {
			return newError("unsupported operation (%s) on %s", operator, left.Type())
		}
		return l.GTE(right)
	case "==":
		l, ok := left.(object.Equal)
		if !ok {
			return newError("unsupported operation (%s) on %s", operator, left.Type())
		}
		return l.EQ(right)
	case "!=":
		l, ok := left.(object.Equal)
		if !ok {
			return newError("unsupported operation (%s) on %s", operator, left.Type())
		}
		return l.NEQ(right)
	default:
		return newError("unknown operator %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalBang(obj object.Object) object.Object {
	b, ok := obj.(object.Booler)
	if !ok {
		return newError("unknown operator (!) on %s", obj.Type())
	}
	return b.Bool().Invert()
}

func evalIf(n *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(n.Condition, env)
	if isError(condition) {
		return condition
	}
	b, ok := condition.(object.Booler)
	if !ok {
		return newError("cannot implicitly convert %s to bool", condition.Type())
	}
	if b.Bool() == object.True {
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
		b, ok := condition.(object.Booler)
		if !ok {
			return newError("cannot implicitly convert %s to bool", condition.Type())
		}
		if b.Bool() == object.False {
			return res
		}
		loopRes := Eval(n.Body, env)
		if isError(loopRes) {
			return loopRes
		}
		if loopRes.Type() == object.ReturnType {
			return loopRes
		}
		if loopRes.Type() == object.BreakType {
			return res
		}
		res = loopRes
	}
}
