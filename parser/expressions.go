package parser

import (
	"github.com/jimmykodes/joker/ast"
	"github.com/jimmykodes/joker/token"
)

type (
	prefixParseFunc func() ast.Expression
	infixParseFunc  func(ast.Expression) ast.Expression
)

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curLit}
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	exp := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curLit,
	}
	p.nextToken()
	exp.Right = p.parseExpression(token.PrefixPrecedence)
	return exp
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()
	exp := p.parseExpression(token.LowestPrecedence)
	if !p.expect(p.peekTokenIs(token.RParen)) {
		p.errors = append(p.errors, invalidTokenError(p.curLine, token.RParen, p.peekToken))
		return nil
	}
	return exp
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	exp := &ast.InfixExpression{
		Token:    p.curToken,
		Left:     left,
		Operator: p.curLit,
	}
	pre := p.curToken.Precedence()
	p.nextToken()
	exp.Right = p.parseExpression(pre)
	return exp
}

func (p *Parser) parseExpressionList(end token.Token) []ast.Expression {
	var list []ast.Expression
	p.nextToken()
	if p.curTokenIs(end, token.EOF) {
		return list
	}
	list = append(list, p.parseExpression(token.LowestPrecedence))
	for p.peekTokenIs(token.Comma) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(token.LowestPrecedence))
	}
	if !p.expect(p.peekTokenIs(end)) {
		p.errors = append(p.errors, invalidTokenError(p.curLine, end, p.peekToken))
		return nil
	}
	return list
}

func (p *Parser) parseIfExpression() ast.Expression {
	exp := &ast.IfExpression{Token: p.curToken}
	p.nextToken()
	exp.Condition = p.parseExpression(token.LowestPrecedence)
	if !p.expect(p.peekTokenIs(token.LBrace)) {
		p.errors = append(p.errors, invalidTokenError(p.curLine, token.LBrace, p.peekToken))
		return nil
	}
	exp.Consequence = p.parseBlockStatement()

	if !p.expect(p.peekTokenIs(token.Else)) {
		// no else, just continue
		return exp
	}

	if !p.expect(p.peekTokenIs(token.LBrace)) {
		p.errors = append(p.errors, invalidTokenError(p.curLine, token.LBrace, p.peekToken))
		return nil
	}

	exp.Alternative = p.parseBlockStatement()

	return exp
}

func (p *Parser) parseForExpression() ast.Expression {
	exp := &ast.ForExpression{Token: p.curToken}
	p.nextToken()

	exp.Init = p.parseStatement()
	p.nextToken()

	exp.Condition = p.parseStatement()
	p.nextToken()

	exp.Increment = p.parseStatement()
	p.nextToken()

	exp.Body = p.parseBlockStatement()
	return exp
}

func (p *Parser) parseWhileExpression() ast.Expression {
	exp := &ast.WhileExpression{Token: p.curToken}
	p.nextToken()
	exp.Condition = p.parseExpression(token.LowestPrecedence)
	if !p.expect(p.peekTokenIs(token.LBrace)) {
		p.errors = append(p.errors, invalidTokenError(p.curLine, token.LBrace, p.peekToken))
		return nil
	}
	exp.Body = p.parseBlockStatement()
	return exp
}

func (p *Parser) parseFuncExpression() ast.Expression {
	exp := &ast.FunctionLiteral{Token: p.curToken}
	if !p.expect(p.peekTokenIs(token.LParen)) {
		p.errors = append(p.errors, invalidTokenError(p.curLine, token.LParen, p.peekToken))
		return nil
	}

	params := p.parseExpressionList(token.RParen)
	exp.Parameters = make([]*ast.Identifier, len(params))
	for i, param := range params {
		cast, ok := param.(*ast.Identifier)
		if !ok {
			p.errors = append(p.errors, newParseError(p.curLine, "invalid type for func param. got %T - want %T", param, &ast.Identifier{}))
			return nil
		}
		exp.Parameters[i] = cast
	}

	if !p.expect(p.peekTokenIs(token.LBrace)) {
		p.errors = append(p.errors, invalidTokenError(p.curLine, token.LBrace, p.peekToken))
		return nil
	}

	exp.Body = p.parseBlockStatement()

	return exp
}

func (p *Parser) parseCallExpression(f ast.Expression) ast.Expression {
	exp := &ast.CallExpression{
		Token:    p.curToken,
		Function: f,
	}
	exp.Arguments = p.parseExpressionList(token.RParen)
	return exp
}

func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	e := &ast.IndexExpression{Token: p.curToken, Left: left}
	p.nextToken()
	e.Index = p.parseExpression(token.LowestPrecedence)
	if !p.expect(p.peekTokenIs(token.RBrack)) {
		p.errors = append(p.errors, invalidTokenError(p.curLine, token.RBrack, p.peekToken))
		return nil
	}
	return e
}
