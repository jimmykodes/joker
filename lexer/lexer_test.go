package lexer

import (
	"strings"
	"testing"

	"github.com/jimmykodes/jk/token"
)

func TestLexer_NextToken(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []token.Token
	}{
		{
			name:  "basic tokens",
			input: "[](){};",
			want: []token.Token{
				newFixedToken(token.LBrack),
				newFixedToken(token.RBrack),
				newFixedToken(token.LParen),
				newFixedToken(token.RParen),
				newFixedToken(token.LBrace),
				newFixedToken(token.RBrace),
				newFixedToken(token.SemiCol),
			},
		},
		{
			name:  "multichar tokens",
			input: "!===>=<=",
			want: []token.Token{
				newFixedToken(token.NEQ),
				newFixedToken(token.EQ),
				newFixedToken(token.GTE),
				newFixedToken(token.LTE),
			},
		},
		{
			name:  "assignment of string",
			input: `let my_val = "test";`,
			want: []token.Token{
				newFixedToken(token.Let),
				{Type: token.Ident, Literal: "my_val"},
				newFixedToken(token.Assign),
				newFixedToken(token.QUOTE),
				{Type: token.Ident, Literal: "test"},
				newFixedToken(token.QUOTE),
				newFixedToken(token.SemiCol),
			},
		},
		{
			name:  "assignment of int",
			input: "let my_int = 5;",
			want: []token.Token{
				{Type: token.Let, Literal: "let"},
				{Type: token.Ident, Literal: "my_int"},
				{Type: token.Assign, Literal: "="},
				{Type: token.Int, Literal: "5"},
				{Type: token.SemiCol, Literal: ";"},
			},
		},
		{
			name:  "assignment of float",
			input: "let my_float = 5.0;",
			want: []token.Token{
				{Type: token.Let, Literal: "let"},
				{Type: token.Ident, Literal: "my_float"},
				{Type: token.Assign, Literal: "="},
				{Type: token.Float, Literal: "5.0"},
				{Type: token.SemiCol, Literal: ";"},
			},
		},
		{
			name:  "function definition",
			input: `let add = fn(a, b) { a + b };`,
			want: []token.Token{
				{Type: token.Let, Literal: "let"},
				{Type: token.Ident, Literal: "add"},
				{Type: token.Assign, Literal: "="},
				{Type: token.Func, Literal: "fn"},
				{Type: token.LParen, Literal: "("},
				{Type: token.Ident, Literal: "a"},
				{Type: token.Comma, Literal: ","},
				{Type: token.Ident, Literal: "b"},
				{Type: token.RParen, Literal: ")"},
				{Type: token.LBrace, Literal: "{"},
				{Type: token.Ident, Literal: "a"},
				{Type: token.Plus, Literal: "+"},
				{Type: token.Ident, Literal: "b"},
				{Type: token.RBrace, Literal: "}"},
				{Type: token.SemiCol, Literal: ";"},
			},
		},
		{
			name:  "function call",
			input: "let result = add(5, 10);",
			want: []token.Token{
				{Type: token.Let, Literal: "let"},
				{Type: token.Ident, Literal: "result"},
				{Type: token.Assign, Literal: "="},
				{Type: token.Ident, Literal: "add"},
				{Type: token.LParen, Literal: "("},
				{Type: token.Int, Literal: "5"},
				{Type: token.Comma, Literal: ","},
				{Type: token.Int, Literal: "10"},
				{Type: token.RParen, Literal: ")"},
				{Type: token.SemiCol, Literal: ";"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.input)
			for i, tok := range tt.want {
				got := l.NextToken()
				if got.Type != tok.Type {
					t.Errorf("NextToken().Type = %v at %d, want %v", got.Type, i, tok.Type)
				}
				if !strings.EqualFold(got.Literal, tok.Literal) {
					t.Errorf("NextToken().Literal = %s at %d, want %s", got.Literal, i, tok.Literal)
				}
			}
		})
	}
}
