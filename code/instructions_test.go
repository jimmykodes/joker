package code

import "testing"

func TestInstructionsString(t *testing.T) {
	inst := []Instructions{
		Make(OpAdd),
		Make(OpConstant, 2),
		Make(OpConstant, 65535),
		Make(OpPop),
	}
	expect := `0000 OpAdd
0001 OpConstant 2
0004 OpConstant 65535
0007 OpPop
`
	var joined Instructions
	for _, ins := range inst {
		joined = append(joined, ins...)
	}
	if joined.String() != expect {
		t.Errorf("invalid string\ngot %q\nwant %q", joined.String(), expect)
	}
}
