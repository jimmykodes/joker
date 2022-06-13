package object

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/jimmykodes/jk/ast"
)

type Object interface {
	Type() Type
	Inspect() string
	Add(Object) (Object, error)
}

type baseObject struct{}

func (b *baseObject) Add(_ Object) (Object, error) {
	return nil, ErrUnsupportedOperation
}

var (
	ErrUnsupportedOperation = errors.New("unsupported operation")
	ErrUnsupportedType      = errors.New("unsupported type for operation")
)

type Null struct{ baseObject }

func (n *Null) Type() Type      { return NullType }
func (n *Null) Inspect() string { return "null" }

type Integer struct {
	baseObject
	Value int64
}

func (i *Integer) Type() Type      { return IntegerType }
func (i *Integer) Inspect() string { return strconv.FormatInt(i.Value, 10) }
func (i *Integer) Add(obj Object) (Object, error) {
	switch o := obj.(type) {
	case *Integer:
		return &Integer{Value: i.Value + o.Value}, nil
	case *Float:
		return &Float{Value: float64(i.Value) + o.Value}, nil
	default:
		return nil, ErrUnsupportedType
	}
}

type Float struct {
	baseObject
	Value float64
}

func (f *Float) Type() Type      { return FloatType }
func (f *Float) Inspect() string { return fmt.Sprintf("%f", f.Value) }
func (f *Float) Add(obj Object) (Object, error) {
	switch o := obj.(type) {
	case *Integer:
		return &Float{Value: f.Value + float64(o.Value)}, nil
	case *Float:
		return &Float{Value: f.Value + o.Value}, nil
	default:
		return nil, ErrUnsupportedType
	}
}

type Boolean struct {
	baseObject
	Value bool
}

func (b *Boolean) Type() Type      { return BoolType }
func (b *Boolean) Inspect() string { return strconv.FormatBool(b.Value) }

type String struct {
	baseObject
	Value string
}

func (s *String) Type() Type      { return StringType }
func (s *String) Inspect() string { return `"` + s.Value + `"` }
func (s *String) Add(obj Object) (Object, error) {
	if o, ok := obj.(*String); ok {
		return &String{Value: s.Value + o.Value}, nil
	}
	return nil, ErrUnsupportedType
}

type Return struct {
	baseObject
	Value Object
}

func (r *Return) Type() Type      { return ReturnType }
func (r *Return) Inspect() string { return r.Value.Inspect() }

type Error struct {
	baseObject
	Message string
}

func (e *Error) Type() Type      { return ErrorType }
func (e *Error) Inspect() string { return e.Message }

type Function struct {
	baseObject
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() Type { return FunctionType }
func (f *Function) Inspect() string {
	var sb strings.Builder
	params := make([]string, len(f.Parameters))
	for i, p := range f.Parameters {
		params[i] = p.String()
	}
	fmt.Fprintf(&sb, "fn(%s) {\n%s\n}", strings.Join(params, ", "), f.Body.String())
	return sb.String()
}

type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	baseObject
	Fn BuiltinFunction
}

func (b *Builtin) Type() Type      { return BuiltinType }
func (b *Builtin) Inspect() string { return "builtin function" }

type Array struct {
	baseObject
	Elements []Object
}

func (a *Array) Type() Type { return ArrayType }

func (a *Array) Inspect() string {
	elements := make([]string, len(a.Elements))
	for i, element := range a.Elements {
		elements[i] = element.Inspect()
	}
	return fmt.Sprintf("[%s]", strings.Join(elements, ", "))
}
