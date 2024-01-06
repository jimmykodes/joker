package lexer

import (
	"strings"
	"testing"

	"github.com/jimmykodes/joker/token"
)

func TestLexer_NextToken(t *testing.T) {
	type result struct {
		token token.Token
		line  int
		lit   string
	}
	tests := []struct {
		name  string
		input string
		want  []result
	}{
		{
			name:  "basic tokens",
			input: "[](){};",
			want: []result{
				{token.LBrack, 1, "["},
				{token.RBrack, 1, "]"},
				{token.LParen, 1, "("},
				{token.RParen, 1, ")"},
				{token.LBrace, 1, "{"},
				{token.RBrace, 1, "}"},
				{token.SemiCol, 1, ";"},
			},
		},
		{
			name:  "bools",
			input: "true false",
			want: []result{
				{token.True, 1, "true"},
				{token.False, 1, "false"},
			},
		},
		{
			name:  "multichar tokens",
			input: "!= == >= <=",
			want: []result{
				{token.NEQ, 1, "!="},
				{token.EQ, 1, "=="},
				{token.GTE, 1, ">="},
				{token.LTE, 1, "<="},
			},
		},
		{
			name:  "assignment of string",
			input: `let my_val = "test";`,
			want: []result{
				{token.Let, 1, "let"},
				{token.Ident, 1, "my_val"},
				{token.Assign, 1, "="},
				{token.String, 1, "test"},
				{token.SemiCol, 1, ";"},
			},
		},
		{
			name:  "assignment of int",
			input: "let my_int = 5;",
			want: []result{
				{token.Let, 1, "let"},
				{token.Ident, 1, "my_int"},
				{token.Assign, 1, "="},
				{token.Int, 1, "5"},
				{token.SemiCol, 1, ";"},
			},
		},
		{
			name:  "assignment of float",
			input: "let my_float = 5.0;",
			want: []result{
				{token.Let, 1, "let"},
				{token.Ident, 1, "my_float"},
				{token.Assign, 1, "="},
				{token.Float, 1, "5.0"},
				{token.SemiCol, 1, ";"},
			},
		},
		{
			name:  "line w/out semicolon",
			input: "a * b",
			want: []result{
				{token.Ident, 1, "a"},
				{token.Mult, 1, "*"},
				{token.Ident, 1, "b"},
			},
		},
		{
			name:  "anonymous function definition with let",
			input: `let add = fn(a, b) { a + b };`,
			want: []result{
				{token.Let, 1, "let"},
				{token.Ident, 1, "add"},
				{token.Assign, 1, "="},
				{token.Func, 1, "fn"},
				{token.LParen, 1, "("},
				{token.Ident, 1, "a"},
				{token.Comma, 1, ","},
				{token.Ident, 1, "b"},
				{token.RParen, 1, ")"},
				{token.LBrace, 1, "{"},
				{token.Ident, 1, "a"},
				{token.Plus, 1, "+"},
				{token.Ident, 1, "b"},
				{token.RBrace, 1, "}"},
				{token.SemiCol, 1, ";"},
			},
		},
		{
			name:  "function definition",
			input: `fn add(a, b) { a + b }`,
			want: []result{
				{token.Func, 1, "fn"},
				{token.Ident, 1, "add"},
				{token.LParen, 1, "("},
				{token.Ident, 1, "a"},
				{token.Comma, 1, ","},
				{token.Ident, 1, "b"},
				{token.RParen, 1, ")"},
				{token.LBrace, 1, "{"},
				{token.Ident, 1, "a"},
				{token.Plus, 1, "+"},
				{token.Ident, 1, "b"},
				{token.RBrace, 1, "}"},
			},
		},
		{
			name: "multiline function definition",
			input: `fn add(a, b) { 
        return a + b;
      }`,
			want: []result{
				{token.Func, 1, "fn"},
				{token.Ident, 1, "add"},
				{token.LParen, 1, "("},
				{token.Ident, 1, "a"},
				{token.Comma, 1, ","},
				{token.Ident, 1, "b"},
				{token.RParen, 1, ")"},
				{token.LBrace, 1, "{"},
				{token.Return, 2, "return"},
				{token.Ident, 2, "a"},
				{token.Plus, 2, "+"},
				{token.Ident, 2, "b"},
				{token.SemiCol, 2, ";"},
				{token.RBrace, 3, "}"},
			},
		},
		{
			name:  "function call",
			input: "let result = add(5, 10);",
			want: []result{
				{token.Let, 1, "let"},
				{token.Ident, 1, "result"},
				{token.Assign, 1, "="},
				{token.Ident, 1, "add"},
				{token.LParen, 1, "("},
				{token.Int, 1, "5"},
				{token.Comma, 1, ","},
				{token.Int, 1, "10"},
				{token.RParen, 1, ")"},
				{token.SemiCol, 1, ";"},
			},
		},
		{
			name: "while loop",
			input: `while i <= 10 {
        i = i + 1;
        continue;
        print("idx", i);
        break;
      }`,
			want: []result{
				{token.While, 1, "while"},
				{token.Ident, 1, "i"},
				{token.LTE, 1, "<="},
				{token.Int, 1, "10"},
				{token.LBrace, 1, "{"},
				{token.Ident, 2, "i"},
				{token.Assign, 2, "="},
				{token.Ident, 2, "i"},
				{token.Plus, 2, "+"},
				{token.Int, 2, "1"},
				{token.SemiCol, 2, ";"},
				{token.Continue, 3, "continue"},
				{token.SemiCol, 3, ";"},
				{token.Ident, 4, "print"},
				{token.LParen, 4, "("},
				{token.String, 4, "idx"},
				{token.Comma, 4, ","},
				{token.Ident, 4, "i"},
				{token.RParen, 4, ")"},
				{token.SemiCol, 4, ";"},
				{token.Break, 5, "break"},
				{token.SemiCol, 5, ";"},
				{token.RBrace, 6, "}"},
			},
		},
		{
			name: "if else",
			input: `if thing == "test" {
        print("yes");
      } else {
        print("no");
      }`,
			want: []result{
				{token.If, 1, "if"},
				{token.Ident, 1, "thing"},
				{token.EQ, 1, "=="},
				{token.String, 1, "test"},
				{token.LBrace, 1, "{"},
				{token.Ident, 2, "print"},
				{token.LParen, 2, "("},
				{token.String, 2, "yes"},
				{token.RParen, 2, ")"},
				{token.SemiCol, 2, ";"},
				{token.RBrace, 3, "}"},
				{token.Else, 3, "else"},
				{token.LBrace, 3, "{"},
				{token.Ident, 4, "print"},
				{token.LParen, 4, "("},
				{token.String, 4, "no"},
				{token.RParen, 4, ")"},
				{token.SemiCol, 4, ";"},
				{token.RBrace, 5, "}"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.input)
			for i, want := range tt.want {
				gotToken, gotLine, gotLit := l.NextToken()
				if gotToken != want.token {
					t.Errorf("invalid token %d: got %s - want %s", i, gotToken, want.token)
					return
				}
				if gotLine != want.line {
					t.Errorf("invalid line for token %d: got %d - want %d", i, gotLine, want.line)
					return
				}
				if !strings.EqualFold(gotLit, want.lit) {
					t.Errorf("invalid literal %d: got %s - want %s", i, gotLit, want.lit)
					return
				}
			}
			tok, _, _ := l.NextToken()
			if tok != token.EOF {
				t.Errorf("tokens remain: got %s", tok)
				return
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
