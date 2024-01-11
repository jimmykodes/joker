package code

import "testing"

func TestInstructionsString(t *testing.T) {
	inst := []Instructions{
		Make(OpAdd),
		Make(OpConstant, 2),
		Make(OpConstant, 65535),
		Make(OpPop),
		Make(OpTrue),
		Make(OpFalse),
		Make(OpBang),
		Make(OpMinus),
		Make(OpNull),
		Make(OpSetGlobal, 2),
		Make(OpGetGlobal, 65535),
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
