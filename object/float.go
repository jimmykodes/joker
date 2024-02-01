package object

import (
	"encoding/binary"
	"fmt"
	"math"
)

type Float struct {
	Value float64
}

func (f *Float) Type() Type      { return FloatType }
func (f *Float) Inspect() string { return fmt.Sprintf("%f", f.Value) }
func (f *Float) UnmarshalBytes(data []byte) (int, error) {
	if t := Type(data[0]); t != f.Type() {
		return 0, fmt.Errorf("invalid type: got %s - want %s", t, f.Type())
	}
	v := binary.BigEndian.Uint64(data[1:])
	f.Value = math.Float64frombits(v)
	return 9, nil
}

func (f *Float) MarshalBytes() ([]byte, error) {
	out := make([]byte, 9)
	out[0] = byte(f.Type())
	binary.BigEndian.PutUint64(out[1:], math.Float64bits(f.Value))
	return out, nil
}

func (f *Float) Bool() *Boolean {
	if f.Value != 0 {
		return True
	}
	return False
}

func (f *Float) HashKey() HashKey {
	return HashKey{Type: FloatType, Value: math.Float64bits(f.Value)}
}

func (f *Float) Negative() Object {
	return &Float{Value: -f.Value}
}

func (f *Float) Add(obj Object) Object {
	switch o := obj.(type) {
	case *Integer:
		return &Float{Value: f.Value + float64(o.Value)}
	case *Float:
		return &Float{Value: f.Value + o.Value}
	default:
		return ErrUnsupportedType
	}
}

func (f *Float) Sub(obj Object) Object {
	switch o := obj.(type) {
	case *Integer:
		return &Float{Value: f.Value - float64(o.Value)}
	case *Float:
		return &Float{Value: f.Value - o.Value}
	default:
		return ErrUnsupportedType
	}
}

func (f *Float) Mult(obj Object) Object {
	switch o := obj.(type) {
	case *Integer:
		return &Float{Value: f.Value * float64(o.Value)}
	case *Float:
		return &Float{Value: f.Value * o.Value}
	default:
		return ErrUnsupportedType
	}
}

func (f *Float) Div(obj Object) Object {
	switch o := obj.(type) {
	case *Integer:
		return &Float{Value: f.Value / float64(o.Value)}
	case *Float:
		return &Float{Value: f.Value / o.Value}
	default:
		return ErrUnsupportedType
	}
}

func (f *Float) LT(obj Object) Object {
	switch o := obj.(type) {
	case *Integer:
		if f.Value < float64(o.Value) {
			return True
		}
	case *Float:
		if f.Value < o.Value {
			return True
		}
	default:
		return ErrUnsupportedType
	}
	return False
}

func (f *Float) LTE(obj Object) Object {
	switch o := obj.(type) {
	case *Integer:
		if f.Value <= float64(o.Value) {
			return True
		}
	case *Float:
		if f.Value <= o.Value {
			return True
		}
	default:
		return ErrUnsupportedType
	}
	return False
}

func (f *Float) GT(obj Object) Object {
	switch o := obj.(type) {
	case *Integer:
		if f.Value > float64(o.Value) {
			return True
		}
	case *Float:
		if f.Value > o.Value {
			return True
		}
	default:
		return ErrUnsupportedType
	}
	return False
}

func (f *Float) GTE(obj Object) Object {
	switch o := obj.(type) {
	case *Integer:
		if f.Value >= float64(o.Value) {
			return True
		}
	case *Float:
		if f.Value >= o.Value {
			return True
		}
	default:
		return ErrUnsupportedType
	}
	return False
}

func (f *Float) EQ(obj Object) Object {
	switch o := obj.(type) {
	case *Integer:
		if f.Value == float64(o.Value) {
			return True
		}
	case *Float:
		if f.Value == o.Value {
			return True
		}
	default:
		return ErrUnsupportedType
	}
	return False
}

func (f *Float) NEQ(obj Object) Object {
	switch o := obj.(type) {
	case *Integer:
		if f.Value != float64(o.Value) {
			return True
		}
	case *Float:
		if f.Value != o.Value {
			return True
		}
	default:
		return ErrUnsupportedType
	}
	return False
}
