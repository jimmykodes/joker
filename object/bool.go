package object

import (
	"strconv"
)

var (
	True  = &Boolean{Value: true}
	False = &Boolean{Value: false}
)

type Boolean struct {
	baseObject
	Value bool
}

func (b *Boolean) Type() Type      { return BoolType }
func (b *Boolean) Inspect() string { return strconv.FormatBool(b.Value) }

func (b *Boolean) Bool() (*Boolean, error) {
	if b.Value {
		return True, nil
	}
	return False, nil
}

func (b *Boolean) HashKey() (*HashKey, error) {
	h := &HashKey{Type: BoolType}
	if b.Value {
		h.Value = 1
	}
	return h, nil
}
