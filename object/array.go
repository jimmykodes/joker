package object

import (
	"fmt"
	"strings"
)

type Array struct {
	baseObject
	Elements []Object
}

func (a *Array) Type() Type { return ArrayType }
func (a *Array) Inspect() string {
	elements := make([]string, len(a.Elements))
	for i, element := range a.Elements {
		elements[i] = element.Inspect()
	}
	return fmt.Sprintf("[%s]", strings.Join(elements, ", "))
}

func (a *Array) Len() (*Integer, error) {
	return &Integer{Value: int64(len(a.Elements))}, nil
}

func (a *Array) Idx(obj Object) (Object, error) {
	o, ok := obj.(*Integer)
	if !ok {
		return nil, ErrUnsupportedType
	}
	if o.Value >= int64(len(a.Elements)) {
		return nil, fmt.Errorf("index out of range [%d] with length %d", o.Value, len(a.Elements))
	}
	return a.Elements[o.Value], nil
}
