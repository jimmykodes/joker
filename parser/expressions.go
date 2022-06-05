package parser

import (
	"fmt"
	"strconv"

	"github.com/jimmykodes/jk/ast"
	"github.com/jimmykodes/jk/token"
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

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()
	exp := p.parseExpression(Lowest)
	if !p.assertAndAdvance(p.peekTokenIs(token.RParen)) {
		p.errors = append(p.errors, fmt.Errorf("missing expected closing paren"))
		return nil
	}
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

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.BooleanLiteral{
		Token: p.curToken,
		Value: p.curTokenIs(token.True),
	}
}

func (p *Parser) parseIfExpression() ast.Expression {
	exp := &ast.IfExpression{Token: p.curToken}
	p.nextToken()
	exp.Condition = p.parseExpression(Lowest)
	if !p.assertAndAdvance(p.peekTokenIs(token.LBrace)) {
		p.errors = append(p.errors, fmt.Errorf("missing expected bracket"))
		return nil
	}
	exp.Consequence = p.parseBlockStatement()

	if !p.assertAndAdvance(p.peekTokenIs(token.Else)) {
		// no else, just continue
		return exp
	}

	if !p.assertAndAdvance(p.peekTokenIs(token.LBrace)) {
		p.errors = append(p.errors, fmt.Errorf("missing expected bracket"))
		return nil
	}

	exp.Alternative = p.parseBlockStatement()

	return exp
}

func (p *Parser) parseFuncExpression() ast.Expression {
	exp := &ast.FunctionLiteral{Token: p.curToken}
	if !p.assertAndAdvance(p.peekTokenIs(token.LParen)) {
		p.errors = append(p.errors, fmt.Errorf("missing required left paren"))
		return nil
	}
	p.nextToken()
	for !p.curTokenIs(token.RParen, token.EOF) {
		ident := p.parseIdentifier()
		if ident == nil {
			p.errors = append(p.errors, fmt.Errorf("invalid identifier"))
			return nil
		}
		exp.Parameters = append(exp.Parameters, ident.(*ast.Identifier))
		if p.peekTokenIs(token.Ident) {
			p.errors = append(p.errors, fmt.Errorf("missing comma between paramenters"))
			return nil
		}
		p.nextToken()
		if p.curTokenIs(token.Comma) {
			p.nextToken()
		}
	}

	if !p.assertAndAdvance(p.peekTokenIs(token.LBrace)) {
		p.errors = append(p.errors, fmt.Errorf("missing required left brace"))
		return nil
	}

	exp.Body = p.parseBlockStatement()

	return exp
}

func (p *Parser) parseCallExpression(f ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: f}
	if p.peekTokenIs(token.RParen) {
		p.nextToken()
		return exp
	}
	p.nextToken()
	exp.Arguments = append(exp.Arguments, p.parseExpression(Lowest))
	for p.peekTokenIs(token.Comma) {
		p.nextToken()
		p.nextToken()
		exp.Arguments = append(exp.Arguments, p.parseExpression(Lowest))
	}
	if !p.assertAndAdvance(p.peekTokenIs(token.RParen)) {
		p.errors = append(p.errors, fmt.Errorf("missing required closing paren"))
		return nil
	}
	return exp
}
