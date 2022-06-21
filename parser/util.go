package parser

import (
	"github.com/jimmykodes/jk/token"
)

func (p *Parser) nextToken() {
	p.curToken, p.curLine, p.curLit = p.peekToken, p.peekLine, p.peekLit
	p.peekToken, p.peekLine, p.peekLit = p.l.NextToken()
}

func (p *Parser) expect(b bool) bool {
	if b {
		p.nextToken()
	}
	return b
}

func (p *Parser) peekTokenIs(t token.Token, ts ...token.Token) bool {
	return tokenIs(p.peekToken, t, ts...)
}

func (p *Parser) curTokenIs(t token.Token, ts ...token.Token) bool {
	return tokenIs(p.curToken, t, ts...)
}

func tokenIs(tok token.Token, t1 token.Token, ts ...token.Token) bool {
	if tok == t1 {
		return true
	}
	for _, t := range ts {
		if tok == t {
			return true
		}
	}
	return false
}
