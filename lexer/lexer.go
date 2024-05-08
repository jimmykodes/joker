package lexer

import "github.com/jimmykodes/joker/token"

func New(data []byte) *Lexer {
	return &Lexer{
		lineNum: 1,
		data:    data,
	}
}

type Lexer struct {
	startPos   int
	currentPos int
	lineNum    int
	data       []byte
}

func (l *Lexer) Next() (token.Token, error) {
	l.strip()
	t := token.Token{Line: l.lineNum}

	l.startPos = l.currentPos
	switch c := l.current(); c {
	case '(':
		t.Type = token.LPar
	case ')':
		t.Type = token.RPar
	case '{':
		t.Type = token.LBrace
	case '}':
		t.Type = token.RBrace
	case '[':
		t.Type = token.LBrack
	case ']':
		t.Type = token.RBrack
	case ',':
		t.Type = token.Comma
	case '.':
		t.Type = token.Dot
	case '|':
		t.Type = token.Pipe
	case ';':
		t.Type = token.SemiColon
	case '+':
		t.Type = token.Plus
	case '-':
		t.Type = token.Minus
	case '*':
		t.Type = token.Mult
	case '%':
		t.Type = token.Mod

	case '/':
		t.Type = l.match('/', token.Comment, token.Div)
		if t.Type == token.Comment {
			for l.peek() != '\n' && l.peek() != 0 {
				l.advance()
			}
			t.Value = l.data[l.startPos : l.currentPos+1]
		}

	case '!':
		t.Type = l.matchEq(token.Bang, token.BangEq)
	case '=':
		t.Type = l.matchEq(token.Eq, token.EqEq)
	case '>':
		t.Type = l.matchEq(token.Gr, token.GrEq)
	case '<':
		t.Type = l.matchEq(token.Ls, token.LsEq)

	case '"':
		t.Type = token.String
		for !l.peekTokenIs('"') {
			l.advance()
		}
		l.advance() // final "
		t.Value = l.data[l.startPos+1 : l.currentPos]

	case 0:
		t.Type = token.EOF

	case '0':
		switch l.peek() {
		case 'x', 'X':
			t.Type = token.Hex
			l.advance()
			for isDigit(l.peek()) || (l.peek() >= 'a' && l.peek() <= 'f') || (l.peek() >= 'A' && l.peek() <= 'F') || l.peek() == '_' {
				l.advance()
			}
		case 'b', 'B':
			t.Type = token.Bin
			l.advance()
			for l.peek() == '0' || l.peek() == '1' || l.peek() == '_' {
				l.advance()
			}
		case 'o', 'O':
			t.Type = token.Oct
			l.advance()
			for (l.peek() >= '0' && l.peek() <= '7') || l.peek() == '_' {
				l.advance()
			}
		case '.':
			t.Type = token.Float
			l.advance()
			for isDigit(l.peek()) || l.peek() == '_' {
				l.advance()
			}
		default:
			t.Type = token.Int
			for isDigit(l.peek()) || l.peek() == '_' {
				l.advance()
			}
		}
		t.Value = l.data[l.startPos : l.currentPos+1]

	default:
		if isDigit(c) {
			t.Type = token.Int
			for isDigit(l.peek()) || l.peek() == '_' {
				l.advance()
			}
			if l.peekTokenIs('.') {
				t.Type = token.Float
				l.advance()
				for isDigit(l.peek()) || l.peek() == '_' {
					l.advance()
				}
			}
			t.Value = l.data[l.startPos : l.currentPos+1]
		} else if isAlpha(c) {
			for p := l.peek(); isAlpha(p) || isDigit(p) || p == '_'; p = l.peek() {
				l.advance()
			}
			v := l.data[l.startPos : l.currentPos+1]
			t.Type = token.LookupIdent(string(v))
			if t.Type == token.Ident {
				t.Value = v
			}
		}
	}

	l.advance()

	return t, nil
}

func (l *Lexer) strip() {
	for {
		switch l.current() {
		case '\n':
			l.lineNum++
			fallthrough
		case '\t', ' ':
			l.advance()
		default:
			return
		}
	}
}

func (l *Lexer) current() byte {
	if l.currentPos >= len(l.data) {
		return 0
	}
	return l.data[l.currentPos]
}

func (l *Lexer) peek() byte {
	if l.currentPos+1 >= len(l.data) {
		return 0
	}
	return l.data[l.currentPos+1]
}

func (l *Lexer) advance() {
	l.currentPos++
}

func (l *Lexer) peekTokenIs(b byte) bool {
	return l.peek() == b
}

func (l *Lexer) match(b byte, eq, neq token.Type) token.Type {
	if l.peekTokenIs(b) {
		l.advance()
		return eq
	}
	return neq
}

func (l *Lexer) matchEq(neq, eq token.Type) token.Type {
	return l.match('=', eq, neq)
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isAlpha(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}
