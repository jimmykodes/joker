package object

import (
	"fmt"
)

type Float struct {
	baseObject
	Value float64
}

func (f *Float) Type() Type      { return FloatType }
func (f *Float) Inspect() string { return fmt.Sprintf("%f", f.Value) }

func (f *Float) Negative() (Object, error) {
	return &Float{Value: -f.Value}, nil
}

func (f *Float) Add(obj Object) (Object, error) {
	switch o := obj.(type) {
	case *Integer:
		return &Float{Value: f.Value + float64(o.Value)}, nil
	case *Float:
		return &Float{Value: f.Value + o.Value}, nil
	default:
		return nil, ErrUnsupportedType
	}
}

func (f *Float) Minus(obj Object) (Object, error) {
	switch o := obj.(type) {
	case *Integer:
		return &Float{Value: f.Value - float64(o.Value)}, nil
	case *Float:
		return &Float{Value: f.Value - o.Value}, nil
	default:
		return nil, ErrUnsupportedType
	}
}

func (f *Float) Mult(obj Object) (Object, error) {
	switch o := obj.(type) {
	case *Integer:
		return &Float{Value: f.Value * float64(o.Value)}, nil
	case *Float:
		return &Float{Value: f.Value * o.Value}, nil
	default:
		return nil, ErrUnsupportedType
	}
}

func (f *Float) Div(obj Object) (Object, error) {
	switch o := obj.(type) {
	case *Integer:
		return &Float{Value: f.Value / float64(o.Value)}, nil
	case *Float:
		return &Float{Value: f.Value / o.Value}, nil
	default:
		return nil, ErrUnsupportedType
	}
}

func (f *Float) LT(obj Object) (Object, error) {
	switch o := obj.(type) {
	case *Integer:
		if f.Value < float64(o.Value) {
			return True, nil
		}
	case *Float:
		if f.Value < o.Value {
			return True, nil
		}
	default:
		return nil, ErrUnsupportedType
	}
	return False, nil
}

func (f *Float) GT(obj Object) (Object, error) {
	switch o := obj.(type) {
	case *Integer:
		if f.Value > float64(o.Value) {
			return True, nil
		}
	case *Float:
		if f.Value > o.Value {
			return True, nil
		}
	default:
		return nil, ErrUnsupportedType
	}
	return False, nil
}
func (f *Float) LTE(obj Object) (Object, error) {
	switch o := obj.(type) {
	case *Integer:
		if f.Value <= float64(o.Value) {
			return True, nil
		}
	case *Float:
		if f.Value <= o.Value {
			return True, nil
		}
	default:
		return nil, ErrUnsupportedType
	}
	return False, nil
}
func (f *Float) GTE(obj Object) (Object, error) {
	switch o := obj.(type) {
	case *Integer:
		if f.Value >= float64(o.Value) {
			return True, nil
		}
	case *Float:
		if f.Value >= o.Value {
			return True, nil
		}
	default:
		return nil, ErrUnsupportedType
	}
	return False, nil
}

func (f *Float) EQ(obj Object) (Object, error) {
	switch o := obj.(type) {
	case *Integer:
		if f.Value == float64(o.Value) {
			return True, nil
		}
	case *Float:
		if f.Value == o.Value {
			return True, nil
		}
	default:
		return nil, ErrUnsupportedType
	}
	return False, nil
}

func (f *Float) NEQ(obj Object) (Object, error) {
	switch o := obj.(type) {
	case *Integer:
		if f.Value != float64(o.Value) {
			return True, nil
		}
	case *Float:
		if f.Value != o.Value {
			return True, nil
		}
	default:
		return nil, ErrUnsupportedType
	}
	return False, nil
}