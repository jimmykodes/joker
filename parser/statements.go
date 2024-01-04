package parser

import (
	"fmt"

	"github.com/jimmykodes/joker/ast"
	"github.com/jimmykodes/joker/token"
)

func (p *Parser) parseLetStatement() ast.Statement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expect(p.peekTokenIs(token.Ident)) {
		p.errors = append(p.errors, invalidTokenError(p.curLine, token.Ident, p.peekToken))
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curLit}

	if !p.expect(p.peekTokenIs(token.Assign)) {
		p.errors = append(p.errors, invalidTokenError(p.curLine, token.Assign, p.peekToken))
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(token.LowestPrecedence)

	if !p.expect(p.peekTokenIs(token.SemiCol)) {
		p.errors = append(p.errors, invalidTokenError(p.curLine, token.SemiCol, p.peekToken))
		return nil
	}

	return stmt
}

func (p *Parser) parseFuncStatement() ast.Statement {
	stmt := &ast.FuncStatement{Token: p.curToken}

	if !p.expect(p.peekTokenIs(token.Ident)) {
		p.errors = append(p.errors, invalidTokenError(p.curLine, token.Ident, p.peekToken))
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curLit}
	stmt.Fn = p.parseFuncExpression().(*ast.FunctionLiteral)

	return stmt
}

func (p *Parser) parseReassignStatement() ast.Statement {
	stmt := &ast.ReassignStatement{
		Name: &ast.Identifier{Token: p.curToken, Value: p.curLit},
	}
	if !p.expect(p.peekTokenIs(token.Assign)) {
		p.errors = append(p.errors, fmt.Errorf("identifier not followed by assignment"))
		return nil
	}
	p.nextToken()

	stmt.Value = p.parseExpression(token.LowestPrecedence)

	if !p.expect(p.peekTokenIs(token.SemiCol)) {
		p.errors = append(p.errors, invalidTokenError(p.curLine, token.SemiCol, p.peekToken))
		return nil
	}

	return stmt
}

func (p *Parser) parseContinueStatement() ast.Statement {
	stmt := &ast.ContinueStatement{Token: p.curToken}
	if !p.expect(p.peekTokenIs(token.SemiCol)) {
		p.errors = append(p.errors, invalidTokenError(p.curLine, token.SemiCol, p.peekToken))
		return nil
	}
	return stmt
}

func (p *Parser) parseBreakStatement() ast.Statement {
	stmt := &ast.BreakStatement{Token: p.curToken}
	if !p.expect(p.peekTokenIs(token.SemiCol)) {
		p.errors = append(p.errors, invalidTokenError(p.curLine, token.SemiCol, p.peekToken))
		return nil
	}
	return stmt
}

func (p *Parser) parseReturnStatement() ast.Statement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()
	stmt.Value = p.parseExpression(token.LowestPrecedence)
	if p.peekTokenIs(token.SemiCol) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() ast.Statement {
	stmt := &ast.ExpressionStatement{
		Token:      p.curToken,
		Expression: p.parseExpression(token.LowestPrecedence),
	}
	if p.peekTokenIs(token.SemiCol) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	p.nextToken()
	for !p.curTokenIs(token.RBrace, token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}
	return block
}
