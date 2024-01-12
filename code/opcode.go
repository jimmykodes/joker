package code

import "fmt"

type Opcode byte

//go:generate stringer -type Opcode
const (
	// stack manipulation
	OpConstant Opcode = iota
	OpPop

	// arithmetic
	OpAdd
	OpSub
	OpMult
	OpDiv
	OpMod

	// bool
	OpTrue
	OpFalse
	OpNull

	// Comparison
	OpEQ
	OpNEQ
	OpGT
	OpGTE

	// prefix
	OpMinus
	OpBang

	// jump
	OpJump
	OpJumpNotTruthy

	// variables
	OpSetGlobal
	OpGetGlobal

	lastOpcode
)

var opWidths = [lastOpcode][]int{
	OpConstant:      {2},
	OpJump:          {2},
	OpJumpNotTruthy: {2},
	OpSetGlobal:     {2},
	OpGetGlobal:     {2},
}

func OpWidths(op byte) ([]int, error) {
	if op >= byte(lastOpcode) {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}
	return opWidths[op], nil
}
