package object

import (
	"encoding/binary"
	"fmt"
	"strconv"
)

type Integer struct {
	Value int64
}

func (i *Integer) Type() Type      { return IntegerType }
func (i *Integer) Inspect() string { return strconv.FormatInt(i.Value, 10) }

func (i *Integer) UnmarshalBytes(data []byte) (int, error) {
	if t := Type(data[0]); t != i.Type() {
		return 0, fmt.Errorf("invalid type: got %s - want %s", t, i.Type())
	}
	i.Value = int64(binary.BigEndian.Uint64(data[1:]))
	return 9, nil
}

func (i *Integer) MarshalBytes() ([]byte, error) {
	out := make([]byte, 9)
	out[0] = byte(IntegerType)
	binary.BigEndian.PutUint64(out[1:], uint64(i.Value))
	return out, nil
}

func (i *Integer) Bool() *Boolean {
	if i.Value != 0 {
		return True
	}
	return False
}

func (i *Integer) Negative() Object {
	return &Integer{Value: -i.Value}
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
