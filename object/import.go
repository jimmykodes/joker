package object

type Import struct {
	Env  *Environment
	File string
}

func (i *Import) Access() *Environment {
	return i.Env
}

func (i *Import) Type() Type      { return ImportType }
func (i *Import) Inspect() string { return "import " + i.File }
