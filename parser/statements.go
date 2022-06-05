package parser

import (
	"fmt"

	"github.com/jimmykodes/jk/ast"
	"github.com/jimmykodes/jk/token"
)

func (p *Parser) parseLetStatement() ast.Statement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.assertAndAdvance(p.peekTokenIs(token.Ident)) {
		fmt.Println("parse error at ident")
		p.errors = append(p.errors, invalidToken(token.Ident, p.peekToken.Type))
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.assertAndAdvance(p.peekTokenIs(token.Assign)) {
		fmt.Println("parse error at assing")
		p.errors = append(p.errors, invalidToken(token.Assign, p.peekToken.Type))
		return nil
	}

	// todo: eval expressions

	for !p.curTokenIs(token.SemiCol) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() ast.Statement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	// todo: parse expressions
	for !p.curTokenIs(token.SemiCol) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseExpressionStatement() ast.Statement {
	stmt := &ast.ExpressionStatement{
		Token:      p.curToken,
		Expression: p.parseExpression(Lowest),
	}
	if p.peekTokenIs(token.SemiCol) {
		p.nextToken()
	}
	return stmt
}
