// Code generated by "stringer -type=Type"; DO NOT EDIT.

package object

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[NullType-0]
	_ = x[IntegerType-1]
	_ = x[FloatType-2]
	_ = x[BoolType-3]
	_ = x[StringType-4]
	_ = x[FunctionType-5]
	_ = x[CompiledFunctionType-6]
	_ = x[ClosureType-7]
	_ = x[BuiltinType-8]
	_ = x[ArrayType-9]
	_ = x[MapType-10]
	_ = x[ReturnType-11]
	_ = x[ContinueType-12]
	_ = x[BreakType-13]
	_ = x[ErrorType-14]
	_ = x[FileType-15]
}

const _Type_name = "NullTypeIntegerTypeFloatTypeBoolTypeStringTypeFunctionTypeCompiledFunctionTypeClosureTypeBuiltinTypeArrayTypeMapTypeReturnTypeContinueTypeBreakTypeErrorTypeFileType"

var _Type_index = [...]uint8{0, 8, 19, 28, 36, 46, 58, 78, 89, 100, 109, 116, 126, 138, 147, 156, 164}

func (i Type) String() string {
	if i < 0 || i >= Type(len(_Type_index)-1) {
		return "Type(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Type_name[_Type_index[i]:_Type_index[i+1]]
}
