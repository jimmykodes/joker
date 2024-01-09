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
		{
			op:       OpConstant,
			operands: []int{math.MaxUint16 - 1},
			want:     []byte{byte(OpConstant), 0xFF, 0xFE},
		},
		{
			op:       OpAdd,
			operands: []int{},
			want:     []byte{byte(OpAdd)},
		},
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
					t.Errorf("incorrect byte at pos %d: got %d - want %d", i, instruction[i], b)
				}
			}
		})
	}
}
