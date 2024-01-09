package code

import "fmt"

type Definition struct {
	Name          string
	OperandWidths []int
}

var definitions = [lastOpcode]*Definition{
	OpConstant: {OpConstant.String(), []int{2}},
	OpAdd:      {OpAdd.String(), []int{}},
}

// Lookup will return the Definition of the provided opcode, or an error
// if the provided byte is undefined
func Lookup(op byte) (*Definition, error) {
	if op >= byte(lastOpcode) {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}
	return definitions[op], nil
}
