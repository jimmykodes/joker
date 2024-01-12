package object

import (
	"fmt"
	"strings"

	"github.com/jimmykodes/joker/ast"
	"github.com/jimmykodes/joker/code"
)

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

type CompiledFunction struct {
	Instructions code.Instructions
}

func (f *CompiledFunction) Type() Type { return CompiledFunctionType }
func (f *CompiledFunction) Inspect() string {
	return fmt.Sprintf("CompiledFunction[%p]", f)
}
