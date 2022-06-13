package parser

import (
	"github.com/jimmykodes/jk/token"
)

type Precedence uint

const (
	_ Precedence = iota
	Lowest
	EQ
	LGT
	Sum
	Product
	Prefix
	Call
	Index
)

var precedences = map[token.Type]Precedence{
	token.EQ:     EQ,
	token.NEQ:    EQ,
	token.LT:     LGT,
	token.LTE:    LGT,
	token.GT:     LGT,
	token.GTE:    LGT,
	token.Minus:  Sum,
	token.Plus:   Sum,
	token.Div:    Product,
	token.Mult:   Product,
	token.Mod:    Product,
	token.LParen: Call,
	token.LBrack: Index,
}

func (p *Parser) peekPrecedence() Precedence {
	if pre, ok := precedences[p.peekToken.Type]; ok {
		return pre
	}
	return Lowest
}

func (p *Parser) curPrecedence() Precedence {
	if pre, ok := precedences[p.curToken.Type]; ok {
		return pre
	}
	return Lowest
}
