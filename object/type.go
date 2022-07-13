package object

type Type int

//go:generate stringer -type=Type
const (
	NullType Type = iota
	IntegerType
	FloatType
	BoolType
	StringType
	FunctionType
	BuiltinType
	ArrayType
	MapType
	ReturnType
	ContinueType
	BreakType
	ErrorType
	ImportType
)
