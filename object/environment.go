package object

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

type Environment struct {
	store map[string]Object
	outer *Environment
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
}
