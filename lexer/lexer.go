package lexer

import (
	"github.com/jimmykodes/jk/token"
)

func New(input string) *Lexer {
	l := &Lexer{input: input, lineNum: 1}
	l.next()
	return l
}

type Lexer struct {
	input        string
	position     int
	readPosition int
	lineNum      int
	ch           byte
}

func (l *Lexer) NextToken() (token.Token, int, string) {

	l.stripWhitespace()
	var (
		tok token.Token
		lit string
	)
	switch {
	case isLetter(l.ch):
		tok, lit = l.readIdent()
		return tok, l.lineNum, lit
	case isDigit(l.ch):
		tok, lit = l.readNumber()
		return tok, l.lineNum, lit
	default:
		switch l.ch {
		case 0:
			tok = token.EOF
		case '<':
			tok = l.switchEQ(token.LT, token.LTE)
		case '>':
			tok = l.switchEQ(token.GT, token.GTE)
		case '!':
			tok = l.switchEQ(token.NOT, token.NEQ)
		case '=':
			tok = l.switchEQ(token.Assign, token.EQ)
		case '(':
			tok = token.LParen
		case ')':
			tok = token.RParen
		case '{':
			tok = token.LBrace
		case '}':
			tok = token.RBrace
		case '[':
			tok = token.LBrack
		case ']':
			tok = token.RBrack
		case '+':
			tok = token.Plus
		case '-':
			tok = token.Minus
		case '*':
			tok = token.Mult
		case '/':
			tok = token.Div
		case '%':
			tok = token.Mod
		case ',':
			tok = token.Comma
		case ';':
			tok = token.SemiCol
		case ':':
			tok = token.Colon
		case '"':
			l.next()
			lit = l.readMultiple(func(b byte) bool { return b != '"' })
			tok = token.String
		}
	}
	if lit == "" {
		lit = tok.String()
	}
	l.next()
	return tok, l.lineNum, lit
}

func (l *Lexer) next() {
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
		l.next()
	}
	return l.input[startPos:l.position]
}

func (l *Lexer) stripWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		if l.ch == '\n' || l.ch == '\r' {
			l.lineNum++
		}
		l.next()
	}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) readNumber() (token.Token, string) {
	startPos := l.position
	tok := token.Illegal

	if l.ch != '.' {
		tok = token.Int
		l.readDigits()
	}

	if l.ch == '.' {
		tok = token.Float
		l.next()
		l.readDigits()
	}
	return tok, l.input[startPos:l.position]
}

func (l *Lexer) readIdent() (token.Token, string) {
	startPos := l.position
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' {
		l.next()
	}
	ident := l.input[startPos:l.position]
	return token.Lookup(ident), ident
}

func (l *Lexer) readDigits() {
	for isDigit(l.ch) {
		l.next()
	}
}

func (l *Lexer) switchEQ(tok0, tok1 token.Token) token.Token {
	if l.peekChar() == '=' {
		l.advancePos()
		return tok1
	}
	return tok0
}
