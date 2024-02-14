package object

import (
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"strings"
)

type String struct {
	Value string
}

func (s *String) Type() Type      { return StringType }
func (s *String) Inspect() string { return `"` + s.Value + `"` }
func (s *String) String() string  { return s.Value }

func (s *String) UnmarshalBytes(data []byte) (int, error) {
	if t := Type(data[0]); t != s.Type() {
		return 0, fmt.Errorf("invalid type: got %s - want %s", t, s.Type())
	}
	strLen := binary.BigEndian.Uint64(data[1:])

	s.Value = string(data[9 : strLen+9])

	return int(strLen) + 9, nil
}

func (s *String) MarshalBytes() ([]byte, error) {
	out := make([]byte, 9, len(s.Value)+9)
	out[0] = byte(s.Type())
	binary.BigEndian.PutUint64(out[1:], uint64(len(s.Value)))
	return append(out, []byte(s.Value)...), nil
}

func (s *String) Bool() *Boolean {
	if s.Value != "" {
		return True
	}
	return False
}

func (s *String) Len() *Integer {
	return &Integer{Value: int64(len(s.Value))}
}

func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return HashKey{
		Type:  StringType,
		Value: h.Sum64(),
	}
}

func (s *String) Add(obj Object) Object {
	if o, ok := obj.(*String); ok {
		return &String{Value: s.Value + o.Value}
	}
	return ErrUnsupportedType
}

func (s *String) Mult(obj Object) Object {
	switch o := obj.(type) {
	case *Integer:
		return &String{Value: strings.Repeat(s.Value, int(o.Value))}
	default:
		return ErrUnsupportedType
	}
}

func (s *String) LT(obj Object) Object {
	if o, ok := obj.(*String); ok {
		if s.Value < o.Value {
			return True
		}
		return False
	}
	return ErrUnsupportedType
}

func (s *String) GT(obj Object) Object {
	if o, ok := obj.(*String); ok {
		if s.Value > o.Value {
			return True
		}
		return False
	}
	return ErrUnsupportedType
}

func (s *String) LTE(obj Object) Object {
	if o, ok := obj.(*String); ok {
		if s.Value <= o.Value {
			return True
		}
		return False
	}
	return ErrUnsupportedType
}

func (s *String) GTE(obj Object) Object {
	if o, ok := obj.(*String); ok {
		if s.Value >= o.Value {
			return True
		}
		return False
	}
	return ErrUnsupportedType
}

func (s *String) EQ(obj Object) Object {
	if o, ok := obj.(*String); ok {
		if s.Value == o.Value {
			return True
		}
		return False
	}
	return ErrUnsupportedType
}

func (s *String) NEQ(obj Object) Object {
	if o, ok := obj.(*String); ok {
		if s.Value != o.Value {
			return True
		}
		return False
	}
	return ErrUnsupportedType
}

func (s *String) Idx(obj Object) Object {
	o, ok := obj.(*Integer)
	if !ok {
		return ErrUnsupportedType
	}
	if o.Value >= int64(len(s.Value)) {
		return &Error{Message: fmt.Sprintf("index out of range [%d] with length %d", o.Value, len(s.Value))}
	}
	return &String{Value: string(s.Value[int(o.Value)])}
}
