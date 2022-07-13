package lexer

import (
	"testing"

	"github.com/jimmykodes/joker/token"
)

// func TestLexer_NextToken(t *testing.T) {
// 	tests := []struct {
// 		name  string
// 		input string
// 		want  []token.Token
// 	}{
// 		{
// 			name:  "basic tokens",
// 			input: "[](){};",
// 			want: []token.Token{
// 				newFixedToken(token.LBrack),
// 				newFixedToken(token.RBrack),
// 				newFixedToken(token.LParen),
// 				newFixedToken(token.RParen),
// 				newFixedToken(token.LBrace),
// 				newFixedToken(token.RBrace),
// 				newFixedToken(token.SemiCol),
// 			},
// 		},
// 		{
// 			name:  "multichar tokens",
// 			input: "!===>=<=",
// 			want: []token.Token{
// 				newFixedToken(token.NEQ),
// 				newFixedToken(token.EQ),
// 				newFixedToken(token.GTE),
// 				newFixedToken(token.LTE),
// 			},
// 		},
// 		{
// 			name:  "assignment of string",
// 			input: `let my_val = "test";`,
// 			want: []token.Token{
// 				newFixedToken(token.Let),
// 				{Type: token.Ident, Literal: "my_val"},
// 				newFixedToken(token.Assign),
// 				{Type: token.String, Literal: "test"},
// 			},
// 		},
// 		{
// 			name:  "assignment of int",
// 			input: "let my_int = 5;",
// 			want: []token.Token{
// 				{Type: token.Let, Literal: "let"},
// 				{Type: token.Ident, Literal: "my_int"},
// 				{Type: token.Assign, Literal: "="},
// 				{Type: token.Int, Literal: "5"},
// 				{Type: token.SemiCol, Literal: ";"},
// 			},
// 		},
// 		{
// 			name:  "assignment of float",
// 			input: "let my_float = 5.0;",
// 			want: []token.Token{
// 				{Type: token.Let, Literal: "let"},
// 				{Type: token.Ident, Literal: "my_float"},
// 				{Type: token.Assign, Literal: "="},
// 				{Type: token.Float, Literal: "5.0"},
// 				{Type: token.SemiCol, Literal: ";"},
// 			},
// 		},
// 		{
// 			name:  "line w/out semicolon",
// 			input: "a * b",
// 			want: []token.Token{
// 				newToken(token.Ident, 'a'),
// 				newFixedToken(token.Mult),
// 				newToken(token.Ident, 'b'),
// 			},
// 		},
// 		{
// 			name:  "function definition",
// 			input: `let add = fn(a, b) { a + b };`,
// 			want: []token.Token{
// 				{Type: token.Let, Literal: "let"},
// 				{Type: token.Ident, Literal: "add"},
// 				{Type: token.Assign, Literal: "="},
// 				{Type: token.Func, Literal: "fn"},
// 				{Type: token.LParen, Literal: "("},
// 				{Type: token.Ident, Literal: "a"},
// 				{Type: token.Comma, Literal: ","},
// 				{Type: token.Ident, Literal: "b"},
// 				{Type: token.RParen, Literal: ")"},
// 				{Type: token.LBrace, Literal: "{"},
// 				{Type: token.Ident, Literal: "a"},
// 				{Type: token.Plus, Literal: "+"},
// 				{Type: token.Ident, Literal: "b"},
// 				{Type: token.RBrace, Literal: "}"},
// 				{Type: token.SemiCol, Literal: ";"},
// 			},
// 		},
// 		{
// 			name:  "function call",
// 			input: "let result = add(5, 10);",
// 			want: []token.Token{
// 				{Type: token.Let, Literal: "let"},
// 				{Type: token.Ident, Literal: "result"},
// 				{Type: token.Assign, Literal: "="},
// 				{Type: token.Ident, Literal: "add"},
// 				{Type: token.LParen, Literal: "("},
// 				{Type: token.Int, Literal: "5"},
// 				{Type: token.Comma, Literal: ","},
// 				{Type: token.Int, Literal: "10"},
// 				{Type: token.RParen, Literal: ")"},
// 				{Type: token.SemiCol, Literal: ";"},
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			l := New(tt.input)
// 			for i, tok := range tt.want {
// 				got := l.NextToken()
// 				if got.Type != tok.Type {
// 					t.Errorf("NextToken().Type = (%v) at %d, want (%v)", got.Type, i, tok.Type)
// 				}
// 				if !strings.EqualFold(got.Literal, tok.Literal) {
// 					t.Errorf("NextToken().Literal = (%s) at %d, want (%s)", got.Literal, i, tok.Literal)
// 				}
// 			}
// 		})
// 	}
// }

func TestLexer(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []token.Token
	}{

		{
			name:  "basic",
			input: "let x = 12;",
			want: []token.Token{
				token.Let,
				token.Ident,
				token.Assign,
				token.Int,
				token.SemiCol,
			},
		},
		{
			name:  "import",
			input: `import "test"`,
			want:  []token.Token{token.Import, token.String},
		},
		{
			name:  "dot",
			input: "strings.join",
			want:  []token.Token{token.Ident, token.Dot, token.Ident},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := New(test.input)
			for _, w := range test.want {
				tok, _, _ := l.NextToken()
				if tok != w {
					t.Errorf("invalid token - got %s - want %s", tok, w)
				}
			}
		})
	}
}

func TestLexer_readNumber(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		token   token.Token
		literal string
	}{
		{
			name:    "integer",
			input:   "123",
			token:   token.Int,
			literal: "123",
		},
		{
			name:    "float",
			input:   "123.02",
			token:   token.Float,
			literal: "123.02",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.input)
			tok, lit := l.readNumber()
			if tok != tt.token {
				t.Errorf("invalid token. want %s - got %s", tt.token, tok)
			}
			if lit != tt.literal {
				t.Errorf("invalid literal. want %s - got %s", tt.literal, lit)
			}
		})
	}
}

func TestLexer_readIdent(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		token   token.Token
		literal string
	}{
		{
			name:    "func",
			input:   "fn",
			token:   token.Func,
			literal: "fn",
		},
		{
			name:    "break",
			input:   "break",
			token:   token.Break,
			literal: "break",
		},
		{
			name:    "camel var",
			input:   "someVar",
			token:   token.Ident,
			literal: "someVar",
		},
		{
			name:    "snake var",
			input:   "some_var",
			token:   token.Ident,
			literal: "some_var",
		},
		{
			name:    "camel with num",
			input:   "someVar1",
			token:   token.Ident,
			literal: "someVar1",
		},
		{
			name:    "snake with num",
			input:   "some_var_1",
			token:   token.Ident,
			literal: "some_var_1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.input)
			tok, lit := l.readIdent()
			if tok != tt.token {
				t.Errorf("invalid token. want %s - got %s", tt.token, tok)
			}
			if lit != tt.literal {
				t.Errorf("invalid literal. want %s - got %s", tt.literal, lit)
			}
		})
	}
}
