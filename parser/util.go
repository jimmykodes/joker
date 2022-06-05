package parser

import (
	"fmt"

	"github.com/jimmykodes/jk/token"
)

func invalidToken(expected, got token.Type) error {
	return fmt.Errorf("parser: invalid token. expected: %s - got: %s", expected, got)
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) assertAndAdvance(b bool) bool {
	if b {
		p.nextToken()
	}
	return b
}

func (p *Parser) peekTokenIs(t token.Type) bool {
	return p.peekToken.Type == t
}

func (p *Parser) curTokenIs(t token.Type, ts ...token.Type) bool {
	if p.curToken.Type == t {
		return true
	}
	for _, tt := range ts {
		if tt == t {
			return true
		}
	}
	return false
}
