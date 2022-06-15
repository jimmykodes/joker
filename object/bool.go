package object

import (
	"strconv"
)

var (
	True  = &Boolean{Value: true}
	False = &Boolean{Value: false}
)

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() Type      { return BoolType }
func (b *Boolean) Inspect() string { return strconv.FormatBool(b.Value) }

func (b *Boolean) Bool() *Boolean {
	return b
}

func (b *Boolean) Invert() *Boolean {
	if b == True {
		return False
	}
	return True
}

func (b *Boolean) HashKey() HashKey {
	h := HashKey{Type: BoolType}
	if b.Value {
		h.Value = 1
	}
	return h
}
