package parser

import (
	"strconv"

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
	}
	p.nextToken()
	return h
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
	return &ast.CallExpression{
		Token:     p.curToken,
		Function:  f,
		Arguments: p.parseExpressionList(token.RParen),
	}
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
