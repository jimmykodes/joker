package parser

import (
	"fmt"

	"github.com/jimmykodes/jk/ast"
	"github.com/jimmykodes/jk/token"
)

func (p *Parser) parseLetStatement() ast.Statement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.assertAndAdvance(p.peekTokenIs(token.Ident)) {
		p.errors = append(p.errors, invalidTokenError(p.curToken.Line, token.Ident, p.peekToken.Type))
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.assertAndAdvance(p.peekTokenIs(token.Assign)) {
		p.errors = append(p.errors, invalidTokenError(p.curToken.Line, token.Assign, p.peekToken.Type))
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(Lowest)

	if !p.assertAndAdvance(p.peekTokenIs(token.SemiCol)) {
		p.errors = append(p.errors, invalidTokenError(p.curToken.Line, token.SemiCol, p.peekToken.Type))
		return nil
	}

	return stmt
}

func (p *Parser) parseFuncStatement() ast.Statement {
	stmt := &ast.FuncStatement{Token: p.curToken}

	if !p.assertAndAdvance(p.peekTokenIs(token.Ident)) {
		p.errors = append(p.errors, invalidTokenError(p.curToken.Line, token.Ident, p.peekToken.Type))
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	stmt.Fn = p.parseFuncExpression().(*ast.FunctionLiteral)

	if !p.assertAndAdvance(p.peekTokenIs(token.SemiCol)) {
		p.errors = append(p.errors, invalidTokenError(p.curToken.Line, token.SemiCol, p.peekToken.Type))
		return nil
	}

	return stmt
}

func (p *Parser) parseReassignStatement() ast.Statement {
	stmt := &ast.ReassignStatement{
		Name: &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal},
	}
	if !p.assertAndAdvance(p.peekTokenIs(token.Assign)) {
		p.errors = append(p.errors, fmt.Errorf("identifier not followed by assignment"))
		return nil
	}
	p.nextToken()

	stmt.Value = p.parseExpression(Lowest)

	if !p.assertAndAdvance(p.peekTokenIs(token.SemiCol)) {
		p.errors = append(p.errors, invalidTokenError(p.curToken.Line, token.SemiCol, p.peekToken.Type))
		return nil
	}

	return stmt
}

func (p *Parser) parseContinueStatement() ast.Statement {
	stmt := &ast.ContinueStatement{Token: p.curToken}
	if !p.assertAndAdvance(p.peekTokenIs(token.SemiCol)) {
		p.errors = append(p.errors, invalidTokenError(p.curToken.Line, token.SemiCol, p.peekToken.Type))
		return nil
	}
	return stmt
}

func (p *Parser) parseBreakStatement() ast.Statement {
	stmt := &ast.BreakStatement{Token: p.curToken}
	if !p.assertAndAdvance(p.peekTokenIs(token.SemiCol)) {
		p.errors = append(p.errors, invalidTokenError(p.curToken.Line, token.SemiCol, p.peekToken.Type))
		return nil
	}
	return stmt
}

func (p *Parser) parseReturnStatement() ast.Statement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()
	stmt.Value = p.parseExpression(Lowest)
	if p.peekTokenIs(token.SemiCol) {
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
