package parser

import (
	"strconv"

	"github.com/jimmykodes/joker/ast"
	"github.com/jimmykodes/joker/token"
)

func (p *Parser) parseCommentLiteral() ast.Expression {
	return &ast.CommentLiteral{
		Token: p.curToken,
		Value: p.curLit,
	}
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{
		Token: p.curToken,
		Value: p.curLit,
	}
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	a := &ast.ArrayLiteral{Token: p.curToken}
	a.Elements = p.parseExpressionList(token.RBrack)
	return a
}

func (p *Parser) parseHashLiteral() ast.Expression {
	h := &ast.MapLiteral{Token: p.curToken, Pairs: make(map[ast.Expression]ast.Expression)}
	for !p.peekTokenIs(token.RBrace) {
		p.nextToken()
		key := p.parseExpression(token.LowestPrecedence)
		if !p.expect(p.peekTokenIs(token.Colon)) {
			p.errors = append(p.errors, invalidTokenError(p.curLine, token.Colon, p.peekToken))
			return nil
		}
		p.nextToken()
		val := p.parseExpression(token.LowestPrecedence)
		h.Pairs[key] = val
		if !p.peekTokenIs(token.RBrace, token.Comma) {
			p.errors = append(p.errors, newParseError(
				p.curLine,
				"invalid token. expected: %s or %s - got: %s",
				token.Comma,
				token.RBrace,
				p.peekToken,
			))
			return nil
		}
		if p.peekTokenIs(token.Comma) {
			p.nextToken()
		}
	}
	p.nextToken()
	return h
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	i, err := strconv.ParseInt(p.curLit, 10, 64)
	if err != nil {
		p.errors = append(p.errors, err)
		return nil
	}
	return &ast.IntegerLiteral{Token: p.curToken, Value: i}
}

func (p *Parser) parseFloatLiteral() ast.Expression {
	i, err := strconv.ParseFloat(p.curLit, 64)
	if err != nil {
		p.errors = append(p.errors, err)
		return nil
	}
	return &ast.FloatLiteral{Token: p.curToken, Value: i}
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.BooleanLiteral{
		Token: p.curToken,
		Value: p.curTokenIs(token.True),
	}
}
