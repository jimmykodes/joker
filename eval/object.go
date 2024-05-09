package eval

import "fmt"

type ObjectType int

//go:generate stringer -type=ObjectType
const (
	ErrType ObjectType = iota
	NilType
	IntType
	FloatType
	StringType
)

type Object interface {
	Type() ObjectType
}

type Addable interface {
	Add(right Object) Object
}

func Err(msg string) Object {
	return ErrObject{Value: msg}
}

func Errf(msg string, a ...any) Object {
	return ErrObject{Value: fmt.Sprintf(msg, a...)}
}

type ErrObject struct{ Value string }

func (o ErrObject) Type() ObjectType { return ErrType }

type IntObject struct{ Value int64 }

func (o IntObject) Type() ObjectType { return IntType }

func (o IntObject) Add(right Object) Object {
	r, ok := right.(IntObject)
	if !ok {
		return Errf("Add not supported between %T and %T", o, right)
	}
	return IntObject{Value: o.Value + r.Value}
}

type FloatObject struct{ Value float64 }

func (o FloatObject) Type() ObjectType { return FloatType }

type StringObject struct{ Value string }

func (o StringObject) Type() ObjectType { return StringType }
