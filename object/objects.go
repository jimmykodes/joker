package object

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jimmykodes/jk/ast"
)

var (
	ErrUnsupportedOperation = errors.New("unsupported operation")
	ErrUnsupportedType      = errors.New("unsupported type for operation")
)

type Null struct{ baseObject }

func (n *Null) Type() Type      { return NullType }
func (n *Null) Inspect() string { return "null" }

func (n *Null) Bool() (bool, error) {
	return false, nil
}

type Continue struct {
	baseObject
}

func (c *Continue) Type() Type      { return ContinueType }
func (c *Continue) Inspect() string { return "continue" }

type Break struct {
	baseObject
}

func (b *Break) Type() Type      { return BreakType }
func (b *Break) Inspect() string { return "break" }

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
