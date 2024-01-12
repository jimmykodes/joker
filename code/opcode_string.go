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
	_ = x[OpNull-9]
	_ = x[OpEQ-10]
	_ = x[OpNEQ-11]
	_ = x[OpGT-12]
	_ = x[OpGTE-13]
	_ = x[OpMinus-14]
	_ = x[OpBang-15]
	_ = x[OpJump-16]
	_ = x[OpJumpNotTruthy-17]
	_ = x[OpSetGlobal-18]
	_ = x[OpGetGlobal-19]
	_ = x[OpArray-20]
	_ = x[OpMap-21]
	_ = x[OpIndex-22]
	_ = x[lastOpcode-23]
}

const _Opcode_name = "OpConstantOpPopOpAddOpSubOpMultOpDivOpModOpTrueOpFalseOpNullOpEQOpNEQOpGTOpGTEOpMinusOpBangOpJumpOpJumpNotTruthyOpSetGlobalOpGetGlobalOpArrayOpMapOpIndexlastOpcode"

var _Opcode_index = [...]uint8{0, 10, 15, 20, 25, 31, 36, 41, 47, 54, 60, 64, 69, 73, 78, 85, 91, 97, 112, 123, 134, 141, 146, 153, 163}

func (i Opcode) String() string {
	if i >= Opcode(len(_Opcode_index)-1) {
		return "Opcode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Opcode_name[_Opcode_index[i]:_Opcode_index[i+1]]
}
