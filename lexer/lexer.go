package lexer

import (
	"github.com/jimmykodes/jk/token"
)

func New(input string) *Lexer {
	l := &Lexer{input: input, lineNum: 1}
	l.readChar()
	return l
}

type Lexer struct {
	input        string
	position     int
	readPosition int
	lineNum      int
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
			tok = newFixedToken(token.LTE, l.lineNum)
		} else {
			tok = newFixedToken(token.LT, l.lineNum)
		}
	case '>':
		if next := l.peekChar(); next == '=' {
			l.advancePos()
			tok = newFixedToken(token.GTE, l.lineNum)
		} else {
			tok = newFixedToken(token.GT, l.lineNum)
		}
	case '!':
		if next := l.peekChar(); next == '=' {
			l.advancePos()
			tok = newFixedToken(token.NEQ, l.lineNum)
		} else {
			tok = newFixedToken(token.NOT, l.lineNum)
		}
	case '=':
		if next := l.peekChar(); next == '=' {
			l.advancePos()
			tok = newFixedToken(token.EQ, l.lineNum)
		} else {
			tok = newFixedToken(token.Assign, l.lineNum)
		}
	case '(':
		tok = newFixedToken(token.LParen, l.lineNum)
	case ')':
		tok = newFixedToken(token.RParen, l.lineNum)
	case '{':
		tok = newFixedToken(token.LBrace, l.lineNum)
	case '}':
		tok = newFixedToken(token.RBrace, l.lineNum)
	case '[':
		tok = newFixedToken(token.LBrack, l.lineNum)
	case ']':
		tok = newFixedToken(token.RBrack, l.lineNum)
	case '+':
		tok = newFixedToken(token.Plus, l.lineNum)
	case '-':
		tok = newFixedToken(token.Minus, l.lineNum)
	case '*':
		tok = newFixedToken(token.Mult, l.lineNum)
	case '/':
		tok = newFixedToken(token.Div, l.lineNum)
	case '%':
		tok = newFixedToken(token.Mod, l.lineNum)
	case ',':
		tok = newFixedToken(token.Comma, l.lineNum)
	case ';':
		tok = newFixedToken(token.SemiCol, l.lineNum)
	case ':':
		tok = newFixedToken(token.Colon, l.lineNum)
	case '"':
		l.readChar()
		tok.Literal = l.readMultiple(func(b byte) bool { return b != '"' })
		tok.Type = token.String
		tok.Line = l.lineNum
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readMultiple(isLetter)
			tok.Type = token.IdentType(tok.Literal)
			tok.Line = l.lineNum
			return tok
		}
		if isDigit(l.ch) {
			tok.Literal = l.readMultiple(isDigit)
			tok.Type = token.NumericType(tok.Literal)
			tok.Line = l.lineNum
			return tok
		}
		tok = newToken(token.Illegal, l.ch, l.lineNum)
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
	for tester(l.ch) && l.ch != 0 {
		l.readChar()
	}
	return l.input[startPos:l.position]
}

func (l *Lexer) stripWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		if l.ch == '\n' || l.ch == '\r' {
			l.lineNum++
		}
		l.readChar()
	}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9' || ch == '.'
}

func newToken(t token.Type, ch byte, line int) token.Token {
	return token.Token{Type: t, Literal: string(ch), Line: line}
}

func newFixedToken(t token.Type, line int) token.Token {
	return token.Token{Type: t, Literal: t.String(), Line: line}
}
