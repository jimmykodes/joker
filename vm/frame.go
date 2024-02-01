package vm

import (
	"github.com/jimmykodes/joker/code"
	"github.com/jimmykodes/joker/object"
)

func NewFrame(cl *object.Closure, basePointer int) *Frame {
	return &Frame{cl: cl, ip: -1, basePointer: basePointer}
}

type Frame struct {
	cl          *object.Closure
	ip          int
	basePointer int
}

func (f Frame) Instructions() code.Instructions {
	return f.cl.Fn.Instructions
}
