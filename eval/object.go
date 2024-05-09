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

var Nil = NilObject{}

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

type NilObject struct{}

func (o NilObject) Type() ObjectType { return NilType }

type IntObject struct{ Value int64 }

func (o IntObject) Type() ObjectType { return IntType }

func (o IntObject) Add(right Object) Object {
	r, ok := right.(IntObject)
	if !ok {
		return Errf("Add not supported between %s and %s", o.Type(), right.Type())
	}
	return IntObject{Value: o.Value + r.Value}
}

type FloatObject struct{ Value float64 }

func (o FloatObject) Type() ObjectType { return FloatType }
func (o FloatObject) Add(right Object) Object {
	r, ok := right.(FloatObject)
	if !ok {
		return Errf("Add not supported between %s and %s", o.Type(), right.Type())
	}
	return FloatObject{Value: o.Value + r.Value}
}

type StringObject struct{ Value string }

func (o StringObject) Type() ObjectType { return StringType }
func (o StringObject) Add(right Object) Object {
	r, ok := right.(StringObject)
	if !ok {
		return Errf("Add not supported between %s and %s", o.Type(), right.Type())
	}
	return StringObject{Value: o.Value + r.Value}
}
