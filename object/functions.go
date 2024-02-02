package object

import (
	"encoding/binary"
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

type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Name string
	Fn   BuiltinFunction
}

func (b *Builtin) Type() Type      { return BuiltinType }
func (b *Builtin) Inspect() string { return fmt.Sprintf("builtin function: %s", b.Name) }

type CompiledFunction struct {
	Instructions code.Instructions
	NumLocals    int
	NumParams    int
}

func (f *CompiledFunction) Type() Type      { return CompiledFunctionType }
func (f *CompiledFunction) Inspect() string { return fmt.Sprintf("CompiledFunction[%p]", f) }

func (f *CompiledFunction) UnmarshalBytes(data []byte) (int, error) {
	if t := Type(data[0]); t != f.Type() {
		return 0, fmt.Errorf("invalid type: got %s - want %s", t, f.Type())
	}

	f.NumLocals = int(binary.BigEndian.Uint64(data[1:])) // this is probably less than uint64
	f.NumParams = int(binary.BigEndian.Uint64(data[9:])) // this is probably less than uint64

	lenIns, err := f.Instructions.UnmarshalBytes(data[17:])
	if err != nil {
		return 0, err
	}

	return lenIns + 17, nil
}

func (f *CompiledFunction) MarshalBytes() ([]byte, error) {
	out := make([]byte, 17)

	out[0] = byte(f.Type())
	binary.BigEndian.PutUint64(out[1:], uint64(f.NumLocals))
	binary.BigEndian.PutUint64(out[9:], uint64(f.NumParams))

	ins, err := f.Instructions.MarshalBytes()
	if err != nil {
		return nil, err
	}

	return append(out, ins...), nil
}

type Closure struct {
	Fn   *CompiledFunction
	Free []Object
}

func (c *Closure) Type() Type      { return ClosureType }
func (c *Closure) Inspect() string { return fmt.Sprintf("Closure[%p]", c) }
