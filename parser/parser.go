package parser

import (
	"github.com/jimmykodes/joker/ast"
	"github.com/jimmykodes/joker/lexer"
	"github.com/jimmykodes/joker/token"
)

type Parser struct {
	l      *lexer.Lexer
	errors []error

	curToken token.Token
	curLine  int
	curLit   string

	peekToken token.Token
	peekLine  int
	peekLit   string

	prefixParseFuncs map[token.Token]prefixParseFunc
	infixParseFuncs  map[token.Token]infixParseFunc
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	p.nextToken()
	p.nextToken()

	p.prefixParseFuncs = map[token.Token]prefixParseFunc{
		token.Ident:   p.parseIdentifier,
		token.NOT:     p.parsePrefixExpression,
		token.Minus:   p.parsePrefixExpression,
		token.Comment: p.parseCommentLiteral,
		token.Int:     p.parseIntegerLiteral,
		token.Float:   p.parseFloatLiteral,
		token.String:  p.parseStringLiteral,
		token.LBrack:  p.parseArrayLiteral,
		token.LBrace:  p.parseHashLiteral,
		token.True:    p.parseBoolean,
		token.False:   p.parseBoolean,
		token.LParen:  p.parseGroupedExpression,
		token.If:      p.parseIfExpression,
		token.While:   p.parseWhileExpression,
		token.Func:    p.parseFuncExpression,
	}

	p.infixParseFuncs = map[token.Token]infixParseFunc{
		token.Plus:   p.parseInfixExpression,
		token.Minus:  p.parseInfixExpression,
		token.Mult:   p.parseInfixExpression,
		token.Div:    p.parseInfixExpression,
		token.Mod:    p.parseInfixExpression,
		token.LT:     p.parseInfixExpression,
		token.GT:     p.parseInfixExpression,
		token.LTE:    p.parseInfixExpression,
		token.GTE:    p.parseInfixExpression,
		token.EQ:     p.parseInfixExpression,
		token.NEQ:    p.parseInfixExpression,
		token.LParen: p.parseCallExpression,
		token.LBrack: p.parseIndexExpression,
	}

	return p
}

func (p *Parser) ParseProgram() *ast.Program {
	prog := &ast.Program{
		Statements: make([]ast.Statement, 0),
	}
	for p.curToken != token.EOF {
		if stmt := p.parseStatement(); stmt != nil {
			prog.Statements = append(prog.Statements, stmt)
		}
		p.nextToken()
	}
	return prog
}

func (p *Parser) Errors() []error {
	return p.errors
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken {
	case token.Let:
		return p.parseLetStatement()
	case token.Return:
		return p.parseReturnStatement()
	case token.Continue:
		return p.parseContinueStatement()
	case token.Break:
		return p.parseBreakStatement()
	case token.Ident:
		if p.peekTokenIs(token.Assign) {
			return p.parseReassignStatement()
		}
		if p.peekTokenIs(token.Define) {
			return p.parseDefineStatement()
		}
	case token.Func:
		if p.peekTokenIs(token.Ident) {
			return p.parseFuncStatement()
		}
	}
	return p.parseExpressionStatement()
}

func (p *Parser) parseExpression(pre token.Precedence) ast.Expression {
	prefix := p.prefixParseFuncs[p.curToken]
	if prefix == nil {
		p.errors = append(p.errors, newParseError(p.curLine, "no prefix func found for token type: %s", p.curToken))
		return nil
	}

	leftExp := prefix()

	for !p.peekTokenIs(token.SemiCol) && pre < p.peekToken.Precedence() {
		infix := p.infixParseFuncs[p.peekToken]
		if infix == nil {
			return leftExp
		}

		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}
