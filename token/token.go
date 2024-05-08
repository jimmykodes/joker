package token

import "fmt"

type Token struct {
	Type  Type
	Line  int
	Value []byte
}

func (t Token) String() string {
	if len(t.Value) > 0 {
		return fmt.Sprintf("%s - %s <lineNum: %d>", t.Type.String(), string(t.Value), t.Line)
	}
	return fmt.Sprintf("%s <lineNum: %d>", t.Type.String(), t.Line)
}
