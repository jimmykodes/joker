package lexer

import (
	"github.com/jimmykodes/jk/token"
)

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func (l *Lexer) NextToken() token.Token {
	l.stripWhitespace()

	var tok token.Token
	switch l.ch {
	case 0:
		tok = token.Token{Type: token.EOF}
	case '<':
		if next := l.peekChar(); next == '=' {
			l.advancePos()
			tok = newFixedToken(token.LTE)
		} else {
			tok = newFixedToken(token.LT)
		}
	case '>':
		if next := l.peekChar(); next == '=' {
			l.advancePos()
			tok = newFixedToken(token.GTE)
		} else {
			tok = newFixedToken(token.GT)
		}
	case '!':
		if next := l.peekChar(); next == '=' {
			l.advancePos()
			tok = newFixedToken(token.NEQ)
		} else {
			tok = newFixedToken(token.NOT)
		}
	case '=':
		if next := l.peekChar(); next == '=' {
			l.advancePos()
			tok = newFixedToken(token.EQ)
		} else {
			tok = newFixedToken(token.Assign)
		}
	case '(':
		tok = newFixedToken(token.LParen)
	case ')':
		tok = newFixedToken(token.RParen)
	case '{':
		tok = newFixedToken(token.LBrace)
	case '}':
		tok = newFixedToken(token.RBrace)
	case '[':
		tok = newFixedToken(token.LBrack)
	case ']':
		tok = newFixedToken(token.RBrack)
	case '+':
		tok = newFixedToken(token.Plus)
	case '-':
		tok = newFixedToken(token.Minus)
	case '*':
		tok = newFixedToken(token.Mult)
	case '/':
		tok = newFixedToken(token.Div)
	case '%':
		tok = newFixedToken(token.Mod)
	case ',':
		tok = newFixedToken(token.Comma)
	case ';':
		tok = newFixedToken(token.SemiCol)
	case '"':
		tok = newFixedToken(token.QUOTE)
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readMultiple(isLetter)
			tok.Type = token.IdentType(tok.Literal)
			return tok
		}
		if isDigit(l.ch) {
			tok.Literal = l.readMultiple(isDigit)
			tok.Type = token.NumericType(tok.Literal)
			return tok
		}
		tok = newToken(token.Illegal, l.ch)
	}
	l.readChar()
	return tok
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.advancePos()
}
func (l *Lexer) advancePos() {
	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) readMultiple(tester func(byte) bool) string {
	startPos := l.position
	for tester(l.ch) {
		l.readChar()
	}
	return l.input[startPos:l.position]
}

func (l *Lexer) stripWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9' || ch == '.'
}

func newToken(t token.Type, ch byte) token.Token {
	return token.Token{Type: t, Literal: string(ch)}
}

func newFixedToken(t token.Type) token.Token {
	return token.Token{Type: t, Literal: t.String()}
}
