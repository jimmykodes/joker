package object

type Object interface {
	Type() Type
	Inspect() string

	Negative() (Object, error)
	Bool() (*Boolean, error)
	HashKey() (*HashKey, error)
	Len() (*Integer, error)

	Add(Object) (Object, error)
	Minus(Object) (Object, error)
	Mult(Object) (Object, error)
	Div(Object) (Object, error)
	Mod(Object) (Object, error)
	LT(Object) (Object, error)
	GT(Object) (Object, error)
	LTE(Object) (Object, error)
	GTE(Object) (Object, error)
	EQ(Object) (Object, error)
	NEQ(Object) (Object, error)
	Idx(Object) (Object, error)
}

type baseObject struct{}

func (b *baseObject) Bool() (*Boolean, error) {
	return nil, ErrUnsupportedOperation
}
func (b *baseObject) Negative() (Object, error) {
	return nil, ErrUnsupportedOperation
}
func (b *baseObject) HashKey() (*HashKey, error) {
	return nil, ErrUnsupportedOperation
}
func (b *baseObject) Len() (*Integer, error) {
	return nil, ErrUnsupportedOperation
}
func (b *baseObject) Add(Object) (Object, error) {
	return nil, ErrUnsupportedOperation
}
func (b *baseObject) Minus(Object) (Object, error) {
	return nil, ErrUnsupportedOperation
}
func (b *baseObject) Mult(Object) (Object, error) {
	return nil, ErrUnsupportedOperation
}
func (b *baseObject) Div(Object) (Object, error) {
	return nil, ErrUnsupportedOperation
}
func (b *baseObject) Mod(Object) (Object, error) {
	return nil, ErrUnsupportedOperation
}
func (b *baseObject) LT(Object) (Object, error) {
	return nil, ErrUnsupportedOperation
}
func (b *baseObject) GT(Object) (Object, error) {
	return nil, ErrUnsupportedOperation
}
func (b *baseObject) LTE(Object) (Object, error) {
	return nil, ErrUnsupportedOperation
}
func (b *baseObject) GTE(Object) (Object, error) {
	return nil, ErrUnsupportedOperation
}
func (b *baseObject) EQ(Object) (Object, error) {
	return nil, ErrUnsupportedOperation
}
func (b *baseObject) NEQ(Object) (Object, error) {
	return nil, ErrUnsupportedOperation
}
func (b *baseObject) Idx(Object) (Object, error) {
	return nil, ErrUnsupportedOperation
}
