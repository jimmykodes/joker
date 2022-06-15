package object

import (
	"strconv"
)

type Integer struct {
	Value int64
}

func (i *Integer) Type() Type      { return IntegerType }
func (i *Integer) Inspect() string { return strconv.FormatInt(i.Value, 10) }

func (i *Integer) Bool() *Boolean {
	if i.Value != 0 {
		return True
	}
	return False
}

func (i *Integer) Negative() (Object, error) {
	return &Integer{Value: -i.Value}, nil
}

func (i *Integer) HashKey() HashKey {
	return HashKey{Type: IntegerType, Value: uint64(i.Value)}
}

func (i *Integer) Add(obj Object) Object {
	switch o := obj.(type) {
	case *Integer:
		return &Integer{Value: i.Value + o.Value}
	case *Float:
		return &Float{Value: float64(i.Value) + o.Value}
	default:
		return ErrUnsupportedType
	}
}

func (i *Integer) Sub(obj Object) Object {
	switch o := obj.(type) {
	case *Integer:
		return &Integer{Value: i.Value - o.Value}
	case *Float:
		return &Float{Value: float64(i.Value) - o.Value}
	default:
		return ErrUnsupportedType
	}
}

func (i *Integer) Mult(obj Object) Object {
	switch o := obj.(type) {
	case *Integer:
		return &Integer{Value: i.Value * o.Value}
	case *Float:
		return &Float{Value: float64(i.Value) * o.Value}
	default:
		return ErrUnsupportedType
	}
}

func (i *Integer) Div(obj Object) Object {
	switch o := obj.(type) {
	case *Integer:
		return &Integer{Value: i.Value / o.Value}
	case *Float:
		return &Float{Value: float64(i.Value) / o.Value}
	default:
		return ErrUnsupportedType
	}
}

func (i *Integer) Mod(obj Object) Object {
	switch o := obj.(type) {
	case *Integer:
		return &Integer{Value: i.Value % o.Value}
	default:
		return ErrUnsupportedType
	}
}

func (i *Integer) LT(obj Object) Object {
	switch o := obj.(type) {
	case *Integer:
		if i.Value < o.Value {
			return True
		}
	case *Float:
		if float64(i.Value) < o.Value {
			return True
		}
	default:
		return ErrUnsupportedType
	}
	return False
}

func (i *Integer) LTE(obj Object) Object {
	switch o := obj.(type) {
	case *Integer:
		if i.Value <= o.Value {
			return True
		}
	case *Float:
		if float64(i.Value) <= o.Value {
			return True
		}
	default:
		return ErrUnsupportedType
	}
	return False
}

func (i *Integer) GT(obj Object) Object {
	switch o := obj.(type) {
	case *Integer:
		if i.Value > o.Value {
			return True
		}
	case *Float:
		if float64(i.Value) > o.Value {
			return True
		}
	default:
		return ErrUnsupportedType
	}
	return False
}

func (i *Integer) GTE(obj Object) Object {
	switch o := obj.(type) {
	case *Integer:
		if i.Value >= o.Value {
			return True
		}
	case *Float:
		if float64(i.Value) >= o.Value {
			return True
		}
	default:
		return ErrUnsupportedType
	}
	return False
}

func (i *Integer) EQ(obj Object) Object {
	switch o := obj.(type) {
	case *Integer:
		if i.Value == o.Value {
			return True
		}
	case *Float:
		if float64(i.Value) == o.Value {
			return True
		}
	default:
		return ErrUnsupportedType
	}
	return False
}

func (i *Integer) NEQ(obj Object) Object {
	switch o := obj.(type) {
	case *Integer:
		if i.Value != o.Value {
			return True
		}
	case *Float:
		if float64(i.Value) != o.Value {
			return True
		}
	default:
		return ErrUnsupportedType
	}
	return False
}
