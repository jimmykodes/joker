package object

import (
	"fmt"
	"strconv"
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
