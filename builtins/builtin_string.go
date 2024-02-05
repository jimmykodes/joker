// Code generated by "stringer -type builtin -linecomment"; DO NOT EDIT.

package builtins

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[start-0]
	_ = x[Int-1]
	_ = x[Float-2]
	_ = x[String-3]
	_ = x[Len-4]
	_ = x[Pop-5]
	_ = x[Print-6]
	_ = x[Append-7]
	_ = x[Set-8]
	_ = x[Slice-9]
	_ = x[Argv-10]
	_ = x[Open-11]
	_ = x[Read-12]
	_ = x[Readline-13]
	_ = x[Write-14]
	_ = x[Close-15]
	_ = x[end-16]
}

const _builtin_name = "startintfloatstringlenpopprintappendsetsliceargvopenreadreadlinewritecloseend"

var _builtin_index = [...]uint8{0, 5, 8, 13, 19, 22, 25, 30, 36, 39, 44, 48, 52, 56, 64, 69, 74, 77}

func (i builtin) String() string {
	if i < 0 || i >= builtin(len(_builtin_index)-1) {
		return "builtin(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _builtin_name[_builtin_index[i]:_builtin_index[i+1]]
}
