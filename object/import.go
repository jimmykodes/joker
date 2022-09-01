package object

import (
	"strings"
)

type Import struct {
	Env  *Environment
	File string
}

func (i *Import) Access() *Environment {
	return i.Env
}

func (i *Import) Type() Type { return ImportType }
func (i *Import) Inspect() string {
	out := make([]string, 0, len(i.Env.store))
	for _, k := range i.Env.Keys() {
		out = append(out, k+": "+i.Env.store[k].Type().String())
	}
	return strings.Join(out, "\n")
}
