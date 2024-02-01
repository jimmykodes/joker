package object

type Null struct{}

func (n *Null) Type() Type      { return NullType }
func (n *Null) Inspect() string { return "null" }

func (n *Null) Bool() *Boolean {
	return False
}
