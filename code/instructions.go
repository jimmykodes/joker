package code

import (
	"fmt"
	"strings"
)

type Instructions []byte

func (ins Instructions) String() string {
	var sb strings.Builder
	i := 0
	for i < len(ins) {
		def, err := Lookup(ins[i])
		if err != nil {
			fmt.Fprintf(&sb, "Error: %s", err)
			continue
		}
		operands, read := ReadOperands(def, ins[i+1:])
		fmt.Fprintf(&sb, "%04d %s\n", i, ins.fmtInstruction(def, operands))
		i += 1 + read
	}
	return sb.String()
}

func (ins Instructions) fmtInstruction(def *Definition, operands []int) string {
	opCount := len(def.OperandWidths)
	if len(operands) != opCount {
		return fmt.Sprintf("Error: operand len %d does not match defined %d\n", len(operands), opCount)
	}
	switch opCount {
	case 0:
		return def.Name
	case 1:
		return fmt.Sprintf("%s %d", def.Name, operands[0])
	}
	return fmt.Sprintf("Error: unhandled operand count (%d) for %s", opCount, def.Name)
}
