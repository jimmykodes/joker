package vm

import (
	"github.com/jimmykodes/joker/code"
)

type Frame struct {
	instructions code.Instructions
	ip           int
}

func NewFrame(inst code.Instructions) *Frame {
	return &Frame{instructions: inst, ip: -1}
}
