package code

import (
	"encoding/binary"
	"fmt"
	"strings"
)

type Instructions []byte

func (ins Instructions) String() string {
	var sb strings.Builder
	i := 0
	for i < len(ins) {
		widths, err := OpWidths(ins[i])
		if err != nil {
			fmt.Fprintf(&sb, "Error: %s", err)
			continue
		}
		operands, read := ReadOperands(widths, ins[i+1:])
		fmt.Fprintf(&sb, "%04d %s\n", i, ins.fmtInstruction(Opcode(ins[i]), widths, operands))
		i += 1 + read
	}
	return sb.String()
}

func (ins Instructions) fmtInstruction(op Opcode, widths []int, operands []int) string {
	opCount := len(widths)
	if len(operands) != opCount {
		return fmt.Sprintf("Error: operand len %d does not match defined %d\n", len(operands), opCount)
	}
	switch opCount {
	case 0:
		return op.String()
	case 1:
		return fmt.Sprintf("%s %d", op.String(), operands[0])
	}
	return fmt.Sprintf("Error: unhandled operand count (%d) for %s", opCount, op.String())
}

func Instruction(op Opcode, operands ...int) []byte {
	if op >= lastOpcode {
		return []byte{}
	}

	widths := opWidths[op]
	instructionLen := 1
	for _, w := range widths {
		instructionLen += w
	}

	instruction := make([]byte, instructionLen)
	instruction[0] = byte(op)
	offset := 1

	for i, o := range operands {
		width := widths[i]
		switch width {
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(o))
		case 1:
			instruction[offset] = byte(o)
		}
		offset += width
	}

	return instruction
}
