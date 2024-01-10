package code

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

	lastOpcode
)
