package code

type (
	Instructions []byte
	Opcode       byte
)

//go:generate stringer -type Opcode
const (
	OpConstant Opcode = iota
	lastOpcode
)
