package parser

import (
	"fmt"

	"github.com/jimmykodes/jk/ast"
	"github.com/jimmykodes/jk/lexer"
	"github.com/jimmykodes/jk/token"
)

type Parser struct {
	l         *lexer.Lexer
	errors    []error
	curToken  token.Token
	peekToken token.Token

	prefixParseFuncs map[token.Type]prefixParseFunc
	infixParseFuncs  map[token.Type]infixParseFunc
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	p.nextToken()
	p.nextToken()

	p.prefixParseFuncs = map[token.Type]prefixParseFunc{
		token.Ident:  p.parseIdentifier,
		token.NOT:    p.parsePrefixExpression,
		token.Minus:  p.parsePrefixExpression,
		token.Int:    p.parseIntegerLiteral,
		token.Float:  p.parseFloatLiteral,
		token.String: p.parseStringLiteral,
		token.True:   p.parseBoolean,
		token.False:  p.parseBoolean,
		token.LParen: p.parseGroupedExpression,
		token.If:     p.parseIfExpression,
		token.Func:   p.parseFuncExpression,
	}

	p.infixParseFuncs = map[token.Type]infixParseFunc{
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
	}

	return p
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

func (p *Parser) Errors() []error {
	return p.errors
}

func (p *Parser) parseStatement() ast.Statement {
	fmt.Println("parsing statement")
	switch p.curToken.Type {
	case token.Let:
		return p.parseLetStatement()
	case token.Return:
		return p.parseReturnStatement()
	case token.Ident:
		if p.peekTokenIs(token.Assign) {
			return p.parseReassignStatement()
		}
		fallthrough
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseExpression(pre Precedence) ast.Expression {
	prefix := p.prefixParseFuncs[p.curToken.Type]
	if prefix == nil {
		p.errors = append(p.errors, fmt.Errorf("no prefix func found for token type: %s", p.curToken.Type))
		return nil
	}

	leftExp := prefix()

	for !p.peekTokenIs(token.SemiCol) && pre < p.peekPrecedence() {
		infix := p.infixParseFuncs[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}
