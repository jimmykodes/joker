package code

import "fmt"

type Definition struct {
	Name          string
	OperandWidths []int
}

var opWidths = [lastOpcode][]int{
	OpConstant:      {2},
	OpJump:          {2},
	OpJumpNotTruthy: {2},
}

var definitions = [lastOpcode]*Definition{}

func init() {
	for i := Opcode(0); i < lastOpcode; i++ {
		definitions[i] = &Definition{i.String(), opWidths[i]}
	}
}

// Lookup will return the Definition of the provided opcode, or an error
// if the provided byte is undefined
func Lookup(op byte) (*Definition, error) {
	if op >= byte(lastOpcode) {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}
	return definitions[op], nil
}
