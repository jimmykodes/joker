package object

import (
	"strings"
)

type String struct {
	baseObject
	Value string
}

func (s *String) Type() Type      { return StringType }
func (s *String) Inspect() string { return `"` + s.Value + `"` }

func (s *String) Add(obj Object) (Object, error) {
	if o, ok := obj.(*String); ok {
		return &String{Value: s.Value + o.Value}, nil
	}
	return nil, ErrUnsupportedType
}

func (s *String) Mult(obj Object) (Object, error) {
	switch o := obj.(type) {
	case *Integer:
		return &String{Value: strings.Repeat(s.Value, int(o.Value))}, nil
	default:
		return nil, ErrUnsupportedType
	}
}

func (s *String) LT(obj Object) (Object, error) {
	if o, ok := obj.(*String); ok {
		if s.Value < o.Value {
			return True, nil
		}
		return False, nil
	}
	return nil, ErrUnsupportedType
}

func (s *String) GT(obj Object) (Object, error) {
	if o, ok := obj.(*String); ok {
		if s.Value > o.Value {
			return True, nil
		}
		return False, nil
	}
	return nil, ErrUnsupportedType
}

func (s *String) LTE(obj Object) (Object, error) {
	if o, ok := obj.(*String); ok {
		if s.Value <= o.Value {
			return True, nil
		}
		return False, nil
	}
	return nil, ErrUnsupportedType
}

func (s *String) GTE(obj Object) (Object, error) {
	if o, ok := obj.(*String); ok {
		if s.Value >= o.Value {
			return True, nil
		}
		return False, nil
	}
	return nil, ErrUnsupportedType
}

func (s *String) EQ(obj Object) (Object, error) {
	if o, ok := obj.(*String); ok {
		if s.Value == o.Value {
			return True, nil
		}
		return False, nil
	}
	return nil, ErrUnsupportedType
}

func (s *String) NEQ(obj Object) (Object, error) {
	if o, ok := obj.(*String); ok {
		if s.Value != o.Value {
			return True, nil
		}
		return False, nil
	}
	return nil, ErrUnsupportedType
}

func (s *String) Idx(obj Object) (Object, error) {
	if o, ok := obj.(*Integer); ok {
		return &String{Value: string(s.Value[int(o.Value)])}, nil
	}
	return nil, ErrUnsupportedType
}
