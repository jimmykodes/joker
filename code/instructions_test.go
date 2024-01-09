package code

import "testing"

func TestInstructionsString(t *testing.T) {
	inst := []Instructions{
		Make(OpConstant, 1),
		Make(OpConstant, 2),
		Make(OpConstant, 65535),
	}
	expect := `0000 OpConstant 1
0003 OpConstant 2
0006 OpConstant 65535
`
	var joined Instructions
	for _, ins := range inst {
		joined = append(joined, ins...)
	}
	if joined.String() != expect {
		t.Errorf("invalid string\ngot %q\nwant %q", joined.String(), expect)
	}
}
