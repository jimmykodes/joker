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

func (p *Parser) peekTokenIs(t token.Type, ts ...token.Type) bool {
	return tokenIs(p.peekToken, t, ts...)
}

func (p *Parser) curTokenIs(t token.Type, ts ...token.Type) bool {
	return tokenIs(p.curToken, t, ts...)
}

func tokenIs(tok token.Token, t1 token.Type, ts ...token.Type) bool {
	if tok.Type == t1 {
		return true
	}
	for _, t := range ts {
		if tok.Type == t {
			return true
		}
	}
	return false
}
