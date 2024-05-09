package parser

import (
	"bytes"
	"strconv"

	"github.com/jimmykodes/joker/ast"
	"github.com/jimmykodes/joker/lexer"
	"github.com/jimmykodes/joker/token"
)

func New(lex *lexer.Lexer) (*Parser, error) {
	p := Parser{lex: lex}
	if err := p.advance(); err != nil {
		return nil, err
	}
	return &p, nil
}

type Parser struct {
	lex       *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
}

func (p *Parser) Parse() (ast.Expr, error) {
	return p.expression()
}

func (p *Parser) advance() error {
	p.curToken = p.peekToken
	var err error
	p.peekToken, err = p.lex.Next()
	return err
}

func (p *Parser) peekTokenIs(t ...token.Type) bool {
	for _, tt := range t {
		if p.peekToken.Type == tt {
			return true
		}
	}
	return false
}

func (p *Parser) curTokenIs(t ...token.Type) bool {
	for _, tt := range t {
		if p.curToken.Type == tt {
			return true
		}
	}
	return false
}

func (p *Parser) expression() (ast.Expr, error) {
	return p.equality()
}

func (p *Parser) equality() (ast.Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}
	for p.peekTokenIs(token.EQ, token.NEQ) {
		if err := p.advance(); err != nil {
			return nil, err
		}

		op := p.curToken
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}

		expr = &ast.BinaryExpr{Left: expr, Right: right, Operator: op}
	}

	return expr, nil
}

func (p *Parser) comparison() (ast.Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.peekTokenIs(token.LT, token.LTEQ, token.GT, token.GTE) {
		if err := p.advance(); err != nil {
			return nil, err
		}

		op := p.curToken
		right, err := p.term()
		if err != nil {
			return nil, err
		}

		expr = &ast.BinaryExpr{Left: expr, Right: right, Operator: op}
	}

	return expr, nil
}

func (p *Parser) term() (ast.Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.peekTokenIs(token.Plus, token.Minus) {
		if err := p.advance(); err != nil {
			return nil, err
		}

		op := p.curToken
		right, err := p.factor()
		if err != nil {
			return nil, err
		}

		expr = &ast.BinaryExpr{Left: expr, Right: right, Operator: op}
	}

	return expr, nil
}

func (p *Parser) factor() (ast.Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.peekTokenIs(token.Div, token.Mult, token.Mod) {
		if err := p.advance(); err != nil {
			return nil, err
		}

		op := p.curToken
		right, err := p.unary()
		if err != nil {
			return nil, err
		}

		expr = &ast.BinaryExpr{Left: expr, Right: right, Operator: op}
	}

	return expr, nil
}

func (p *Parser) unary() (ast.Expr, error) {
	if p.peekTokenIs(token.Bang, token.Minus) {
		if err := p.advance(); err != nil {
			return nil, err
		}

		op := p.curToken
		right, err := p.unary()
		if err != nil {
			return nil, err
		}

		return &ast.UnaryExpr{Right: right, Operator: op}, nil
	}

	return p.primary()
}

func (p *Parser) primary() (ast.Expr, error) {
	if err := p.advance(); err != nil {
		return nil, err
	}
	c := p.curToken
	if p.curTokenIs(token.False, token.True, token.Nil) {
		switch c.Type {
		case token.True:
			return &ast.BoolLitExpr{Token: c, Value: true}, nil
		case token.False:
			return &ast.BoolLitExpr{Token: c, Value: false}, nil
		case token.Nil:
			return &ast.NilLitExpr{Token: c}, nil
		}
	}
	if p.curTokenIs(token.Int, token.Hex, token.Oct, token.Bin) {
		var (
			v   int64
			err error
			s   string = string(bytes.ReplaceAll(c.Value, []byte{'_'}, nil))
		)
		switch c.Type {
		case token.Int:
			v, err = strconv.ParseInt(s, 10, 64)
		case token.Hex:
			v, err = strconv.ParseInt(s[2:], 16, 64)
		case token.Oct:
			v, err = strconv.ParseInt(s[2:], 8, 64)
		case token.Bin:
			v, err = strconv.ParseInt(s[2:], 2, 64)
		}
		if err != nil {
			return nil, ParserError{Token: p.curToken, Message: "invalid int literal", Err: err}
		}
		return &ast.IntLitExpr{Token: p.curToken, Value: v}, nil
	}
	if p.curTokenIs(token.Float) {
		v, err := strconv.ParseFloat(string(c.Value), 64)
		if err != nil {
			return nil, ParserError{Token: p.curToken, Message: "invalid float literal", Err: err}
		}
		return &ast.FloatLitExpr{Token: c, Value: v}, nil
	}

	if p.curTokenIs(token.String) {
		return &ast.StringLitExpr{Token: c, Value: string(c.Value)}, nil
	}

	if p.curTokenIs(token.LPar) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		if !p.peekTokenIs(token.RPar) {
			return nil, ParserError{Token: p.peekToken, Message: "missing expected ')'"}
		}
		if err := p.advance(); err != nil {
			return nil, err
		}
		return &ast.GroupingExpr{Expr: expr, Token: c}, nil
	}

	return nil, ParserError{Token: p.curToken, Message: "invalid token"}
}
