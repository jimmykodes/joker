package token

import (
	"strings"
)

type Token struct {
	Type    Type
	Literal string
	Line    int
}

var keywords = map[string]Type{
	"fn":       Func,
	"let":      Let,
	"if":       If,
	"else":     Else,
	"for":      For,
	"while":    While,
	"continue": Continue,
	"break":    Break,
	"return":   Return,
	"true":     True,
	"false":    False,
}

func IdentType(ident string) Type {
	if t, ok := keywords[ident]; ok {
		return t
	}
	return Ident
}

func NumericType(ident string) Type {
	if strings.Contains(ident, ".") {
		return Float
	}
	return Int
}
