package code

import (
	"math"
	"testing"
)

func TestInstructionsString(t *testing.T) {
	inst := []Instructions{
		Instruction(OpAdd),
		Instruction(OpConstant, 2),
		Instruction(OpConstant, 65535),
		Instruction(OpPop),
		Instruction(OpTrue),
		Instruction(OpFalse),
		Instruction(OpBang),
		Instruction(OpMinus),
		Instruction(OpNull),
		Instruction(OpSetGlobal, 2),
		Instruction(OpGetGlobal, 65535),
	}
	expect := `0000 OpAdd
0001 OpConstant 2
0004 OpConstant 65535
0007 OpPop
0008 OpTrue
0009 OpFalse
0010 OpBang
0011 OpMinus
0012 OpNull
0013 OpSetGlobal 2
0016 OpGetGlobal 65535
`
	var joined Instructions
	for _, ins := range inst {
		joined = append(joined, ins...)
	}
	if joined.String() != expect {
		t.Errorf("invalid string\ngot %q\nwant %q", joined.String(), expect)
	}
}

func TestInstruction(t *testing.T) {
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
			instruction := Instruction(tt.op, tt.operands...)
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
