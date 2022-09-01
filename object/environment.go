package object

import (
	"io"
	"os"
	"sort"
)

func NewEnvironment(opts ...EnvOption) *Environment {
	e := &Environment{store: make(map[string]Object)}
	for _, opt := range opts {
		e = opt(e)
	}
	return e
}

type EnvOption func(*Environment) *Environment

func EncloseOuterOption(outer *Environment) EnvOption {
	return func(e *Environment) *Environment {
		e.outer = outer
		return e
	}
}

func WithOut(out io.Writer) EnvOption {
	return func(e *Environment) *Environment {
		e.out = out
		return e
	}
}

type Environment struct {
	store map[string]Object
	keys  []string
	outer *Environment
	out   io.Writer
}

func (e *Environment) Get(name string) (Object, bool) {
	o, ok := e.store[name]
	if !ok && e.outer != nil {
		o, ok = e.outer.Get(name)
	}
	return o, ok
}

func (e *Environment) Set(name string, val Object) {
	e.store[name] = val
	e.keys = append(e.keys, name)
}

func (e *Environment) Keys() []string {
	if !sort.StringsAreSorted(e.keys) {
		sort.Strings(e.keys)
	}
	return e.keys
}

func (e *Environment) Out() io.Writer {
	if e.out != nil {
		return e.out
	}
	if e.outer != nil && e.outer.Out() != nil {
		return e.outer.Out()
	}
	return os.Stdout
}
