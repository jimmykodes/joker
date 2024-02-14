package object

var ErrUnsupportedType = &Error{Message: "unsupported type for operation"}

type Encodable interface {
	MarshalBytes() ([]byte, error)
	UnmarshalBytes([]byte) (int, error)
}

type Object interface {
	Type() Type
	Inspect() string
}

type Stringer interface {
	String() string
}

type Hashable interface {
	HashKey() HashKey
}

type Adder interface {
	Add(Object) Object
}

type Subber interface {
	Sub(Object) Object
}

type MultDiver interface {
	Mult(Object) Object
	Div(Object) Object
}

type Modder interface {
	Mod(Object) Object
}

type Lenner interface {
	Len() *Integer
}

type Inequality interface {
	LT(Object) Object
	LTE(Object) Object
	GT(Object) Object
	GTE(Object) Object
}

type Equal interface {
	EQ(Object) Object
	NEQ(Object) Object
}

type Indexer interface {
	Idx(Object) Object
}

type Booler interface {
	Bool() *Boolean
}

type Negater interface {
	Negative() Object
}

type Settable interface {
	Set(key, value Object) Object
}

type Closer interface {
	Close() Object
}

type Reader interface {
	Read() Object
}

type Readliner interface {
	Readline() Object
}

type Writer interface {
	Write(Object) Object
}

type Continue struct{}

func (c *Continue) Type() Type      { return ContinueType }
func (c *Continue) Inspect() string { return "continue" }

type Break struct{}

func (b *Break) Type() Type      { return BreakType }
func (b *Break) Inspect() string { return "break" }

type Return struct {
	Value Object
}

func (r *Return) Type() Type      { return ReturnType }
func (r *Return) Inspect() string { return r.Value.Inspect() }

func ErrorFromGo(err error) *Error {
	return &Error{Message: err.Error()}
}

type Error struct {
	Message string
}

func (e *Error) Type() Type      { return ErrorType }
func (e *Error) Inspect() string { return e.Message }
func (e *Error) Error() string   { return e.Message }
