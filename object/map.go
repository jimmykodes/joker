package object

import (
	"fmt"
	"strings"
)

type HashKey struct {
	Type  Type
	Value uint64
}

type HashPair struct {
	Key   Object
	Value Object
}

type Map struct {
	Pairs map[HashKey]HashPair
}

func (m *Map) Type() Type { return MapType }
func (m *Map) Inspect() string {
	var pairs []string
	for _, pair := range m.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}
	return fmt.Sprintf("{%s}", strings.Join(pairs, ", "))
}

func (m *Map) Idx(obj Object) Object {
	hashable, ok := obj.(Hashable)
	if !ok {
		return ErrUnsupportedType
	}
	hk := hashable.HashKey()
	p, ok := m.Pairs[hk]
	if !ok {
		return &Error{Message: "key not present"}
	}
	return p.Value
}

func (m *Map) Set(key, value Object) Object {
	hashable, ok := key.(Hashable)
	if !ok {
		return ErrUnsupportedType
	}
	hk := hashable.HashKey()
	m.Pairs[hk] = HashPair{key, value}
	return nil
}
