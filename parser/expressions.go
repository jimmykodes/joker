package parser

import (
	"strconv"

	"github.com/jimmykodes/jk/ast"
)

type (
	prefixParseFunc func() ast.Expression
	infixParseFunc  func(ast.Expression) ast.Expression
)

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	exp := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}
	p.nextToken()
	exp.Right = p.parseExpression(Prefix)
	return exp
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	exp := &ast.InfixExpression{
		Token:    p.curToken,
		Left:     left,
		Operator: p.curToken.Literal,
	}
	pre := p.curPrecedence()
	p.nextToken()
	exp.Right = p.parseExpression(pre)
	return exp
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{
		Token: p.curToken,
		Value: p.curToken.Literal,
	}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	i, err := strconv.ParseInt(p.curToken.Literal, 10, 64)
	if err != nil {
		p.errors = append(p.errors, err)
		return nil
	}
	return &ast.IntegerLiteral{Token: p.curToken, Value: i}
}

func (p *Parser) parseFloatLiteral() ast.Expression {
	i, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		p.errors = append(p.errors, err)
		return nil
	}
	return &ast.FloatLiteral{Token: p.curToken, Value: i}
}
