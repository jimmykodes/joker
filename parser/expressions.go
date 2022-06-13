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

func (p *Parser) parseArrayLiteral() ast.Expression {
	a := &ast.ArrayLiteral{Token: p.curToken}
	a.Elements = p.parseExpressionList(token.RBrack)
	return a
}

func (p *Parser) parseExpressionList(end token.Type) []ast.Expression {
	var list []ast.Expression
	if p.peekTokenIs(end, token.EOF) {
		return list
	}
	p.nextToken()
	list = append(list, p.parseExpression(Lowest))
	for p.peekTokenIs(token.Comma) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(Lowest))
	}
	if !p.assertAndAdvance(p.peekTokenIs(end)) {
		p.errors = append(p.errors, fmt.Errorf("missing end token"))
		return nil
	}
	return list
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

func (p *Parser) parseWhileExpression() ast.Expression {
	exp := &ast.WhileExpression{Token: p.curToken}
	p.nextToken()
	exp.Condition = p.parseExpression(Lowest)
	if !p.assertAndAdvance(p.peekTokenIs(token.LBrace)) {
		p.errors = append(p.errors, fmt.Errorf("missing expected bracket"))
		return nil
	}
	exp.Body = p.parseBlockStatement()
	return exp
}

func (p *Parser) parseFuncExpression() ast.Expression {
	exp := &ast.FunctionLiteral{Token: p.curToken}
	if !p.assertAndAdvance(p.peekTokenIs(token.LParen)) {
		p.errors = append(p.errors, fmt.Errorf("missing required left paren"))
		return nil
	}

	params := p.parseExpressionList(token.RParen)
	exp.Parameters = make([]*ast.Identifier, len(params))
	for i, param := range params {
		cast, ok := param.(*ast.Identifier)
		if !ok {
			p.errors = append(p.errors, fmt.Errorf("invalid type for func param. got %T - want %T", param, &ast.Identifier{}))
			return nil
		}
		exp.Parameters[i] = cast
	}

	if !p.assertAndAdvance(p.peekTokenIs(token.LBrace)) {
		p.errors = append(p.errors, fmt.Errorf("missing required left brace"))
		return nil
	}

	exp.Body = p.parseBlockStatement()

	return exp
}

func (p *Parser) parseCallExpression(f ast.Expression) ast.Expression {
	return &ast.CallExpression{
		Token:     p.curToken,
		Function:  f,
		Arguments: p.parseExpressionList(token.RParen),
	}
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	e := &ast.IndexExpression{Token: p.curToken, Left: left}
	p.nextToken()
	e.Index = p.parseExpression(Lowest)
	if !p.assertAndAdvance(p.peekTokenIs(token.RBrack)) {
		p.errors = append(p.errors, fmt.Errorf("missing closing bracket"))
		return nil
	}
	return e
}
