package token

import "fmt"

type Token struct {
	Type  Type
	Line  int
	Value []byte
}

func (t Token) String() string {
	if len(t.Value) > 0 {
		return fmt.Sprintf("%d: %s - %s", t.Line, t.Type.String(), string(t.Value))
	}
	return fmt.Sprintf("%d: %s", t.Line, t.Type.String())
}
