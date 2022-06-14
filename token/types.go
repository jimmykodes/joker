package token

import (
	"strconv"
)

type Type uint

const (
	Illegal Type = iota
	EOF

	String
	Float
	Int

	Ident
	Func
	Let
	If
	Else
	For
	While
	Continue
	Break
	Return
	True
	False

	LT  // <
	GT  // >
	LTE // <=
	GTE // >=
	EQ  // ==
	NEQ // !=
	NOT // !

	LParen // (
	RParen // )
	LBrace // {
	RBrace // }
	LBrack // [
	RBrack // ]

	Assign // =
	Plus   // +
	Minus  // -
	Mult   // *
	Div    // /
	Mod    // %

	QUOTE   // "
	Comma   // ,
	SemiCol // ;
)

var tokens = [...]string{
	Illegal:  "ILLEGAL",
	EOF:      "EOF",
	String:   "STRING",
	Float:    "FLOAT",
	Int:      "INT",
	Ident:    "IDENT",
	Func:     "FUNC",
	Let:      "LET",
	If:       "IF",
	Else:     "ELSE",
	For:      "FOR",
	While:    "WHILE",
	Continue: "CONTINUE",
	Break:    "BREAK",
	Return:   "RETURN",
	True:     "TRUE",
	False:    "FALSE",
	LT:       "<",
	GT:       ">",
	LTE:      "<=",
	GTE:      ">=",
	EQ:       "==",
	NEQ:      "!=",
	NOT:      "!",
	LParen:   "(",
	RParen:   ")",
	LBrace:   "{",
	RBrace:   "}",
	LBrack:   "[",
	RBrack:   "]",
	Assign:   "=",
	Plus:     "+",
	Minus:    "-",
	Mult:     "*",
	Div:      "/",
	Mod:      "%",
	QUOTE:    `"`,
	Comma:    ",",
	SemiCol:  ";",
}

func (t Type) String() string {
	s := ""
	if 0 <= t && int(t) <= len(tokens) {
		s = tokens[t]
	}
	if s == "" {
		s = "token(" + strconv.Itoa(int(t)) + ")"
	}
	return s
}
