package eval

type Env struct {
	outter *Env

	symbols map[string]Object
}

func NewEnv(outter *Env) *Env {
	return &Env{outter: outter, symbols: make(map[string]Object)}
}

func (e *Env) Set(name string, value Object) {
	e.symbols[name] = value
}

func (e *Env) Get(name string) Object {
	v, ok := e.symbols[name]
	if ok {
		return v
	}
	if e.outter == nil {
		return Errf("name '%s' not found", name)
	}
	return e.outter.Get(name)
}
