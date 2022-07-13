package object

import (
	"fmt"
	"strings"

	"github.com/jimmykodes/joker/ast"
)

var (
	ErrUnsupportedType = &Error{Message: "unsupported type for operation"}
)

type Object interface {
	Type() Type
	Inspect() string
}

type Hashable interface {
	HashKey() HashKey
}

type Adder interface {
	Add(Object) Object
}

type Subber interface {
	Sub(Object) Object
}

type MultDiver interface {
	Mult(Object) Object
	Div(Object) Object
}

type Modder interface {
	Mod(Object) Object
}

type Lenner interface {
	Len() *Integer
}

type Inequality interface {
	LT(Object) Object
	LTE(Object) Object
	GT(Object) Object
	GTE(Object) Object
}

type Equal interface {
	EQ(Object) Object
	NEQ(Object) Object
}

type Indexer interface {
	Idx(Object) Object
}

type Booler interface {
	Bool() *Boolean
}

type Negater interface {
	Negative() Object
}

type Accessor interface {
	Access() *Environment
}

type Null struct{}

func (n *Null) Type() Type      { return NullType }
func (n *Null) Inspect() string { return "null" }

func (n *Null) Bool() (*Boolean, error) {
	return False, nil
}

type Continue struct{}

func (c *Continue) Type() Type      { return ContinueType }
func (c *Continue) Inspect() string { return "continue" }

type Break struct{}

func (b *Break) Type() Type      { return BreakType }
func (b *Break) Inspect() string { return "break" }

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
func (e *Error) Error() string   { return e.Message }

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

type BuiltinFunction func(env *Environment, args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() Type      { return BuiltinType }
func (b *Builtin) Inspect() string { return "builtin function" }
