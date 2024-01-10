package code

import "fmt"

type Definition struct {
	Name          string
	OperandWidths []int
}

var definitions = [lastOpcode]*Definition{
	// stack manipulation
	OpConstant: {OpConstant.String(), []int{2}},
	OpPop:      {OpPop.String(), []int{}},

	// arithmetic
	OpAdd:  {OpAdd.String(), []int{}},
	OpSub:  {OpSub.String(), []int{}},
	OpMult: {OpMult.String(), []int{}},
	OpDiv:  {OpDiv.String(), []int{}},
	OpMod:  {OpMod.String(), []int{}},

	// bools
	OpTrue:  {OpTrue.String(), []int{}},
	OpFalse: {OpFalse.String(), []int{}},
}

// Lookup will return the Definition of the provided opcode, or an error
// if the provided byte is undefined
func Lookup(op byte) (*Definition, error) {
	if op >= byte(lastOpcode) {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}
	return definitions[op], nil
}
