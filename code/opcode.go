package code

type Opcode byte

//go:generate stringer -type Opcode
const (
	OpConstant Opcode = iota
	lastOpcode
)
