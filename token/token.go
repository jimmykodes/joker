package token

import (
	"strconv"
)

type Token int

const (
	Illegal Token = iota
	EOF

	literalBeg
	String
	Float
	Int
	Ident
	literalEnd

	keywordBeg
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
	keywordEnd

	operatorBeg
	LT  // <
	GT  // >
	LTE // <=
	GTE // >=
	EQ  // ==
	NEQ // !=
	NOT // !

	Assign // =
	Define // :=
	Plus   // +
	Minus  // -
	Mult   // *
	Div    // /
	Mod    // %

	LParen // (
	RParen // )
	LBrace // {
	RBrace // }
	LBrack // [
	RBrack // ]

	Comma   // ,
	SemiCol // ;
	Colon   // :
	operatorEnd
)

var tokens = [...]string{
	Illegal:  "ILLEGAL",
	EOF:      "EOF",
	String:   "STRING",
	Float:    "FLOAT",
	Int:      "INT",
	Ident:    "IDENT",
	Func:     "fn",
	Let:      "let",
	If:       "if",
	Else:     "else",
	For:      "for",
	While:    "while",
	Continue: "continue",
	Break:    "break",
	Return:   "return",
	True:     "true",
	False:    "false",
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
	Define:   ":=",
	Plus:     "+",
	Minus:    "-",
	Mult:     "*",
	Div:      "/",
	Mod:      "%",
	Comma:    ",",
	SemiCol:  ";",
	Colon:    ":",
}

func (t Token) String() string {
	s := ""
	if 0 <= t && int(t) <= len(tokens) {
		s = tokens[t]
	}
	if s == "" {
		s = "token(" + strconv.Itoa(int(t)) + ")"
	}
	return s
}

func (t Token) IsKeyword() bool {
	return inRange(t, keywordBeg, keywordEnd)
}

func (t Token) IsLiteral() bool {
	return inRange(t, literalBeg, literalEnd)
}

func (t Token) IsOperator() bool {
	return inRange(t, operatorBeg, operatorEnd)
}

type Precedence uint

const (
	_ Precedence = iota
	LowestPrecedence
	EQPrecedence
	LGTPrecedence
	SumPrecedence
	ProductPrecedence
	PrefixPrecedence
	CallPrecedence
	IndexPrecedence
)

func (t Token) Precedence() Precedence {
	switch t {
	case EQ, NEQ:
		return EQPrecedence
	case LT, LTE, GT, GTE:
		return LGTPrecedence
	case Minus, Plus:
		return SumPrecedence
	case Div, Mult, Mod:
		return ProductPrecedence
	case LParen:
		return CallPrecedence
	case LBrack:
		return IndexPrecedence
	default:
		return LowestPrecedence
	}
}

var keywords map[string]Token

func init() {
	keywords = make(map[string]Token)
	for i := keywordBeg + 1; i < keywordEnd; i++ {
		keywords[tokens[i]] = i
	}
}

func Lookup(ident string) Token {
	if t, ok := keywords[ident]; ok {
		return t
	}
	return Ident
}

func inRange(item, beg, end Token) bool {
	return beg < item && item < end
}
