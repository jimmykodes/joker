package parser

import (
	"github.com/jimmykodes/jk/ast"
	"github.com/jimmykodes/jk/lexer"
	"github.com/jimmykodes/jk/token"
)

type Parser struct {
	l *lexer.Lexer

	curToken, peekToken token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	prog := &ast.Program{
		Statements: make([]ast.Statement, 0),
	}
	for p.curToken.Type != token.EOF {
		if stmt := p.parseStatement(); stmt != nil {
			prog.Statements = append(prog.Statements, stmt)
		}
		p.nextToken()
	}
	return prog
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.Let:
		return p.parseLetStatement()
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() ast.Statement {
	stmt := &ast.LetStatement{
		Token: p.curToken,
	}
	if !p.assertAndAdvance(p.peekTokenIs(token.Ident)) {
		return nil // maybe raise error?
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.assertAndAdvance(p.peekTokenIs(token.Assign)) {
		return nil // maybe raise error?
	}

	for !p.curTokenIs(token.SemiCol) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) assertAndAdvance(b bool) bool {
	if b {
		p.nextToken()
	}
	return b
}

func (p *Parser) peekTokenIs(t token.Type) bool {
	return p.peekToken.Type == t
}

func (p *Parser) curTokenIs(t token.Type) bool {
	return p.curToken.Type == t
}
