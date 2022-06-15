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
	baseObject
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
func (m *Map) Idx(obj Object) (Object, error) {
	hk, err := obj.HashKey()
	if err != nil {
		return nil, err
	}
	p, ok := m.Pairs[*hk]
	if !ok {
		return nil, fmt.Errorf("key not present")
	}

	return p.Value, nil
}
