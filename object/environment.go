package object

func NewEnvironment() *Environment {
	return &Environment{store: make(map[string]Object)}
}

type Environment struct {
	store map[string]Object
}

func (e *Environment) Get(name string) (Object, bool) {
	o, ok := e.store[name]
	return o, ok
}

func (e *Environment) Set(name string, val Object) {
	e.store[name] = val
}
