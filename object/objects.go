package object

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jimmykodes/jk/ast"
)

type Object interface {
	Type() Type
	Inspect() string
}

type Null struct{}

func (n *Null) Type() Type      { return NullType }
func (n *Null) Inspect() string { return "null" }

type Integer struct {
	Value int64
}

func (i *Integer) Type() Type      { return IntegerType }
func (i *Integer) Inspect() string { return strconv.FormatInt(i.Value, 10) }

type Float struct {
	Value float64
}

func (f *Float) Type() Type      { return FloatType }
func (f *Float) Inspect() string { return fmt.Sprintf("%f", f.Value) }

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() Type      { return BoolType }
func (b *Boolean) Inspect() string { return strconv.FormatBool(b.Value) }

type String struct {
	Value string
}

func (s *String) Type() Type      { return StringType }
func (s *String) Inspect() string { return `"` + s.Value + `"` }

type Return struct {
	Value Object
}

func (r *Return) Type() Type      { return ReturnType }
func (r *Return) Inspect() string { return r.Value.Inspect() }

type Error struct {
	Message string
}

func (e *Error) Type() Type      { return ErrorType }
func (e *Error) Inspect() string { return e.Message }

type Function struct {
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
	Fn BuiltinFunction
}

func (b *Builtin) Type() Type      { return BuiltinType }
func (b *Builtin) Inspect() string { return "builtin function" }
