package vm

import (
	"github.com/jimmykodes/joker/code"
)

type Frame struct {
	instructions code.Instructions
	ip           int
	basePointer  int
}

func NewFrame(inst code.Instructions, basePointer int) *Frame {
	return &Frame{instructions: inst, ip: -1, basePointer: basePointer}
}
