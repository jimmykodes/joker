// Code generated by "stringer -type Opcode"; DO NOT EDIT.

package code

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[OpConstant-0]
	_ = x[OpPop-1]
	_ = x[OpAdd-2]
	_ = x[OpSub-3]
	_ = x[OpMult-4]
	_ = x[OpDiv-5]
	_ = x[OpMod-6]
	_ = x[OpTrue-7]
	_ = x[OpFalse-8]
	_ = x[OpEQ-9]
	_ = x[OpNEQ-10]
	_ = x[OpGT-11]
	_ = x[OpGTE-12]
	_ = x[OpMinus-13]
	_ = x[OpBang-14]
	_ = x[OpJump-15]
	_ = x[OpJumpNotTruthy-16]
	_ = x[lastOpcode-17]
}

const _Opcode_name = "OpConstantOpPopOpAddOpSubOpMultOpDivOpModOpTrueOpFalseOpEQOpNEQOpGTOpGTEOpMinusOpBangOpJumpOpJumpNotTruthylastOpcode"

var _Opcode_index = [...]uint8{0, 10, 15, 20, 25, 31, 36, 41, 47, 54, 58, 63, 67, 72, 79, 85, 91, 106, 116}

func (i Opcode) String() string {
	if i >= Opcode(len(_Opcode_index)-1) {
		return "Opcode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Opcode_name[_Opcode_index[i]:_Opcode_index[i+1]]
}
