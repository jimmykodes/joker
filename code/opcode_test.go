package code

import (
	"testing"
)

func TestReadOperands(t *testing.T) {
	tests := []struct {
		op        Opcode
		operands  []int
		bytesRead int
	}{
		// stack mod
		{OpConstant, []int{65535}, 2},
		{OpPop, []int{}, 0},

		// arithmetic
		{OpAdd, []int{}, 0},
		{OpSub, []int{}, 0},
		{OpMult, []int{}, 0},
		{OpDiv, []int{}, 0},
		{OpMod, []int{}, 0},

		// bool
		{OpTrue, []int{}, 0},
		{OpFalse, []int{}, 0},
		{OpNull, []int{}, 0},

		// comparison
		{OpEQ, []int{}, 0},
		{OpNEQ, []int{}, 0},
		{OpGT, []int{}, 0},
		{OpGTE, []int{}, 0},

		// prefix
		{OpMinus, []int{}, 0},
		{OpBang, []int{}, 0},

		// jump
		{OpJump, []int{12}, 2},
		{OpJumpNotTruthy, []int{22}, 2},

		// variables
		{OpSetGlobal, []int{65535}, 2},
		{OpGetGlobal, []int{65535}, 2},

		// Array
		{OpArray, []int{65535}, 2},
	}
	for _, tt := range tests {
		instruction := Instruction(tt.op, tt.operands...)
		widths, err := OpWidths(byte(tt.op))
		if err != nil {
			t.Fatalf("definition not found: %s", err)
		}
		operandsRead, n := ReadOperands(widths, instruction[1:])
		if n != tt.bytesRead {
			t.Fatalf("invalid bytes read: got %d - want %d", n, tt.bytesRead)
		}
		for i, want := range tt.operands {
			if operandsRead[i] != want {
				t.Fatalf("invalid operand: got %q - want %q", operandsRead[i], want)
			}
		}
	}
}
