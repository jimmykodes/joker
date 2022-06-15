package object

import (
	"strconv"
)

type Integer struct {
	baseObject
	Value int64
}

func (i *Integer) Type() Type      { return IntegerType }
func (i *Integer) Inspect() string { return strconv.FormatInt(i.Value, 10) }

func (i *Integer) Bool() (*Boolean, error) {
	if i.Value != 0 {
		return True, nil
	}
	return False, nil
}

func (i *Integer) Negative() (Object, error) {
	return &Integer{Value: -i.Value}, nil
}

func (i *Integer) HashKey() (*HashKey, error) {
	return &HashKey{Type: IntegerType, Value: uint64(i.Value)}, nil
}

func (i *Integer) Add(obj Object) (Object, error) {
	switch o := obj.(type) {
	case *Integer:
		return &Integer{Value: i.Value + o.Value}, nil
	case *Float:
		return &Float{Value: float64(i.Value) + o.Value}, nil
	default:
		return nil, ErrUnsupportedType
	}
}

func (i *Integer) Minus(obj Object) (Object, error) {
	switch o := obj.(type) {
	case *Integer:
		return &Integer{Value: i.Value - o.Value}, nil
	case *Float:
		return &Float{Value: float64(i.Value) - o.Value}, nil
	default:
		return nil, ErrUnsupportedType
	}
}

func (i *Integer) Mult(obj Object) (Object, error) {
	switch o := obj.(type) {
	case *Integer:
		return &Integer{Value: i.Value * o.Value}, nil
	case *Float:
		return &Float{Value: float64(i.Value) * o.Value}, nil
	default:
		return nil, ErrUnsupportedType
	}
}

func (i *Integer) Div(obj Object) (Object, error) {
	switch o := obj.(type) {
	case *Integer:
		return &Integer{Value: i.Value / o.Value}, nil
	case *Float:
		return &Float{Value: float64(i.Value) / o.Value}, nil
	default:
		return nil, ErrUnsupportedType
	}
}

func (i *Integer) Mod(obj Object) (Object, error) {
	switch o := obj.(type) {
	case *Integer:
		return &Integer{Value: i.Value % o.Value}, nil
	default:
		return nil, ErrUnsupportedType
	}
}

func (i *Integer) LT(obj Object) (Object, error) {
	switch o := obj.(type) {
	case *Integer:
		if i.Value < o.Value {
			return True, nil
		}
	case *Float:
		if float64(i.Value) < o.Value {
			return True, nil
		}
	default:
		return nil, ErrUnsupportedType
	}
	return False, nil
}

func (i *Integer) GT(obj Object) (Object, error) {
	switch o := obj.(type) {
	case *Integer:
		if i.Value > o.Value {
			return True, nil
		}
	case *Float:
		if float64(i.Value) > o.Value {
			return True, nil
		}
	default:
		return nil, ErrUnsupportedType
	}
	return False, nil
}
func (i *Integer) LTE(obj Object) (Object, error) {
	switch o := obj.(type) {
	case *Integer:
		if i.Value <= o.Value {
			return True, nil
		}
	case *Float:
		if float64(i.Value) <= o.Value {
			return True, nil
		}
	default:
		return nil, ErrUnsupportedType
	}
	return False, nil
}
func (i *Integer) GTE(obj Object) (Object, error) {
	switch o := obj.(type) {
	case *Integer:
		if i.Value >= o.Value {
			return True, nil
		}
	case *Float:
		if float64(i.Value) >= o.Value {
			return True, nil
		}
	default:
		return nil, ErrUnsupportedType
	}
	return False, nil
}

func (i *Integer) EQ(obj Object) (Object, error) {
	switch o := obj.(type) {
	case *Integer:
		if i.Value == o.Value {
			return True, nil
		}
	case *Float:
		if float64(i.Value) == o.Value {
			return True, nil
		}
	default:
		return nil, ErrUnsupportedType
	}
	return False, nil
}

func (i *Integer) NEQ(obj Object) (Object, error) {
	switch o := obj.(type) {
	case *Integer:
		if i.Value != o.Value {
			return True, nil
		}
	case *Float:
		if float64(i.Value) != o.Value {
			return True, nil
		}
	default:
		return nil, ErrUnsupportedType
	}
	return False, nil
}
