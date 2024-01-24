package code

import (
	"encoding/binary"
	"fmt"
	"strings"
)

type Instructions []byte

func (ins *Instructions) UnmarshalBytes(data []byte) (int, error) {
	lenIns := int(binary.BigEndian.Uint64(data))
	*ins = data[8 : 8+lenIns]
	return lenIns + 8, nil
}

func (ins Instructions) MarshalBytes() ([]byte, error) {
	lenIns := len(ins)
	out := make([]byte, 8, lenIns+8)
	binary.BigEndian.PutUint64(out, uint64(lenIns))
	out = append(out, ins...)
	return out, nil
}

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

	var sb strings.Builder
	sb.WriteString(op.String())
	for _, operand := range operands {
		fmt.Fprintf(&sb, " %d", operand)
	}
	return sb.String()
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
