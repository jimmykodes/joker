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
		{OpConstant, []int{65535}, 2},
		{OpAdd, []int{}, 0},
		{OpPop, []int{}, 0},
	}
	for _, tt := range tests {
		instruction := Make(tt.op, tt.operands...)
		def, err := Lookup(byte(tt.op))
		if err != nil {
			t.Fatalf("definition not found: %s", err)
		}
		operandsRead, n := ReadOperands(def, instruction[1:])
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
