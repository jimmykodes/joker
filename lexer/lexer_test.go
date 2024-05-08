package lexer_test

import (
	"reflect"
	"testing"

	"github.com/jimmykodes/joker/lexer"
	"github.com/jimmykodes/joker/token"
)

func TestLexer(t *testing.T) {
	tests := []struct {
		input    string
		expected []token.Token
	}{
		{
			input: "+-*%/;",
			expected: []token.Token{
				{Type: token.Plus, Line: 1},
				{Type: token.Minus, Line: 1},
				{Type: token.Mult, Line: 1},
				{Type: token.Mod, Line: 1},
				{Type: token.Div, Line: 1},
				{Type: token.SemiColon, Line: 1},
			},
		},
		{
			input: `+ // this is a comment
// and another comment`,
			expected: []token.Token{
				{Type: token.Plus, Line: 1},
				{Type: token.Comment, Line: 1, Value: []byte("// this is a comment")},
				{Type: token.Comment, Line: 2, Value: []byte("// and another comment")},
			},
		},
		{
			input: `()
			{}
			[]`,
			expected: []token.Token{
				{Type: token.LPar, Line: 1},
				{Type: token.RPar, Line: 1},
				{Type: token.LBrace, Line: 2},
				{Type: token.RBrace, Line: 2},
				{Type: token.LBrack, Line: 3},
				{Type: token.RBrack, Line: 3},
			},
		},
		{
			input: `"this is a string"`,
			expected: []token.Token{
				{Type: token.String, Line: 1, Value: []byte("this is a string")},
			},
		},
		{
			input: `0x1fa
0X1FA
0o123
0O123
0b0101
0B0101
0.123
0123
0xaa_ff_11
0XAA_FF_11
0o12_34
0O12_34
0b01_01
0B01_01
0.12_34
01_23
123
12_34
`,
			expected: []token.Token{
				{Type: token.Hex, Line: 1, Value: []byte("0x1fa")},
				{Type: token.Hex, Line: 2, Value: []byte("0X1FA")},
				{Type: token.Oct, Line: 3, Value: []byte("0o123")},
				{Type: token.Oct, Line: 4, Value: []byte("0O123")},
				{Type: token.Bin, Line: 5, Value: []byte("0b0101")},
				{Type: token.Bin, Line: 6, Value: []byte("0B0101")},
				{Type: token.Float, Line: 7, Value: []byte("0.123")},
				{Type: token.Int, Line: 8, Value: []byte("0123")},

				{Type: token.Hex, Line: 9, Value: []byte("0xaa_ff_11")},
				{Type: token.Hex, Line: 10, Value: []byte("0XAA_FF_11")},
				{Type: token.Oct, Line: 11, Value: []byte("0o12_34")},
				{Type: token.Oct, Line: 12, Value: []byte("0O12_34")},
				{Type: token.Bin, Line: 13, Value: []byte("0b01_01")},
				{Type: token.Bin, Line: 14, Value: []byte("0B01_01")},
				{Type: token.Float, Line: 15, Value: []byte("0.12_34")},
				{Type: token.Int, Line: 16, Value: []byte("01_23")},
				{Type: token.Int, Line: 17, Value: []byte("123")},
				{Type: token.Int, Line: 18, Value: []byte("12_34")},
			},
		},

		{
			input: `let twelve_12 = 12;`,
			expected: []token.Token{
				{Type: token.Let, Line: 1},
				{Type: token.Ident, Line: 1, Value: []byte("twelve_12")},
				{Type: token.Eq, Line: 1},
				{Type: token.Int, Line: 1, Value: []byte("12")},
				{Type: token.SemiColon, Line: 1},
			},
		},
		{
			input: `fn one() {
				return 1; 
			}`,
			expected: []token.Token{
				{Type: token.Func, Line: 1},
				{Type: token.Ident, Line: 1, Value: []byte("one")},
				{Type: token.LPar, Line: 1},
				{Type: token.RPar, Line: 1},
				{Type: token.LBrace, Line: 1},
				{Type: token.Return, Line: 2},
				{Type: token.Int, Line: 2, Value: []byte("1")},
				{Type: token.SemiColon, Line: 2},
				{Type: token.RBrace, Line: 3},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := lexer.New([]byte(tt.input))
			for _, expected := range tt.expected {
				got, err := l.Next()
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}
				if !reflect.DeepEqual(got, expected) {
					t.Errorf("expected %+v, got %+v", expected, got)
					return
				}
			}
		})
	}
}
