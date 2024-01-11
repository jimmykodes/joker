package code

import (
	"math"
	"testing"
)

func TestMake(t *testing.T) {
	tests := []struct {
		op       Opcode
		operands []int
		want     []byte
	}{
		// stack mod
		{OpConstant, []int{math.MaxUint16 - 1}, []byte{byte(OpConstant), 0xFF, 0xFE}},
		{OpPop, []int{}, []byte{byte(OpPop)}},

		// arithmetic
		{OpAdd, []int{}, []byte{byte(OpAdd)}},
		{OpSub, []int{}, []byte{byte(OpSub)}},
		{OpMult, []int{}, []byte{byte(OpMult)}},
		{OpDiv, []int{}, []byte{byte(OpDiv)}},
		{OpMod, []int{}, []byte{byte(OpMod)}},

		// bool
		{OpTrue, []int{}, []byte{byte(OpTrue)}},
		{OpFalse, []int{}, []byte{byte(OpFalse)}},
		{OpNull, []int{}, []byte{byte(OpNull)}},

		// comparison
		{OpEQ, []int{}, []byte{byte(OpEQ)}},
		{OpNEQ, []int{}, []byte{byte(OpNEQ)}},
		{OpGT, []int{}, []byte{byte(OpGT)}},
		{OpGTE, []int{}, []byte{byte(OpGTE)}},

		// prefix
		{OpMinus, []int{}, []byte{byte(OpMinus)}},
		{OpBang, []int{}, []byte{byte(OpBang)}},

		// jump
		{OpJump, []int{math.MaxUint16 - 1}, []byte{byte(OpJump), 0xFF, 0xFE}},
		{OpJumpNotTruthy, []int{math.MaxUint16 - 1}, []byte{byte(OpJumpNotTruthy), 0xFF, 0xFE}},

		// variables
		{OpSetGlobal, []int{math.MaxUint16 - 1}, []byte{byte(OpSetGlobal), 0xFF, 0xFE}},
		{OpGetGlobal, []int{math.MaxUint16 - 1}, []byte{byte(OpGetGlobal), 0xFF, 0xFE}},
	}
	for _, tt := range tests {
		t.Run(tt.op.String(), func(t *testing.T) {
			instruction := Make(tt.op, tt.operands...)
			if len(instruction) != len(tt.want) {
				t.Errorf("instruction has incorrect length: got %d - want %d", len(instruction), len(tt.want))
				return
			}
			for i, b := range tt.want {
				if instruction[i] != b {
					t.Errorf("incorrect instruction at pos %d: got %s - want %s", i, Opcode(instruction[i]), Opcode(b))
				}
			}
		})
	}
}
