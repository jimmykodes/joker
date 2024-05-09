package parser_test

import (
	"reflect"
	"testing"

	"github.com/jimmykodes/joker/ast"
	"github.com/jimmykodes/joker/lexer"
	"github.com/jimmykodes/joker/parser"
	"github.com/jimmykodes/joker/token"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    ast.Expr
		expectedErr error
	}{
		{
			name:        "simple addition",
			input:       "12 + 2",
			expectedErr: nil,
			expected: &ast.BinaryExpr{
				Left:     &ast.IntLitExpr{Token: token.Token{Type: token.Int, Line: 1, Value: []byte{'1', '2'}}, Value: 12},
				Right:    &ast.IntLitExpr{Token: token.Token{Type: token.Int, Line: 1, Value: []byte{'2'}}, Value: 2},
				Operator: token.Token{Type: token.Plus, Line: 1},
			},
		},
		{
			name:        "nested addition",
			input:       "12 + 2 + 1",
			expectedErr: nil,
			expected: &ast.BinaryExpr{
				Left: &ast.BinaryExpr{
					Left:     &ast.IntLitExpr{Token: token.Token{Type: token.Int, Line: 1, Value: []byte{'1', '2'}}, Value: 12},
					Right:    &ast.IntLitExpr{Token: token.Token{Type: token.Int, Line: 1, Value: []byte{'2'}}, Value: 2},
					Operator: token.Token{Type: token.Plus, Line: 1},
				},
				Right:    &ast.IntLitExpr{Token: token.Token{Type: token.Int, Line: 1, Value: []byte{'1'}}, Value: 1},
				Operator: token.Token{Type: token.Plus, Line: 1},
			},
		},
		{
			name:        "nested addition",
			input:       "12 + 2 + 1",
			expectedErr: nil,
			expected: &ast.BinaryExpr{
				Left: &ast.BinaryExpr{
					Left:     &ast.IntLitExpr{Token: token.Token{Type: token.Int, Line: 1, Value: []byte{'1', '2'}}, Value: 12},
					Right:    &ast.IntLitExpr{Token: token.Token{Type: token.Int, Line: 1, Value: []byte{'2'}}, Value: 2},
					Operator: token.Token{Type: token.Plus, Line: 1},
				},
				Right:    &ast.IntLitExpr{Token: token.Token{Type: token.Int, Line: 1, Value: []byte{'1'}}, Value: 1},
				Operator: token.Token{Type: token.Plus, Line: 1},
			},
		},
		{
			name:        "nested multiplication",
			input:       "12 + 2 * 1",
			expectedErr: nil,
			expected: &ast.BinaryExpr{
				Left: &ast.IntLitExpr{Token: token.Token{Type: token.Int, Line: 1, Value: []byte{'1', '2'}}, Value: 12},
				Right: &ast.BinaryExpr{
					Left:     &ast.IntLitExpr{Token: token.Token{Type: token.Int, Line: 1, Value: []byte{'2'}}, Value: 2},
					Right:    &ast.IntLitExpr{Token: token.Token{Type: token.Int, Line: 1, Value: []byte{'1'}}, Value: 1},
					Operator: token.Token{Type: token.Mult, Line: 1},
				},
				Operator: token.Token{Type: token.Plus, Line: 1},
			},
		},
		{
			name:        "grouped multiplication",
			input:       "(12 + 2) * 1",
			expectedErr: nil,
			expected: &ast.BinaryExpr{
				Left: &ast.GroupingExpr{
					Expr: &ast.BinaryExpr{
						Left:     &ast.IntLitExpr{Token: token.Token{Type: token.Int, Line: 1, Value: []byte{'1', '2'}}, Value: 12},
						Right:    &ast.IntLitExpr{Token: token.Token{Type: token.Int, Line: 1, Value: []byte{'2'}}, Value: 2},
						Operator: token.Token{Type: token.Plus, Line: 1},
					},
					Token: token.Token{Type: token.LPar, Line: 1},
				},
				Right:    &ast.IntLitExpr{Token: token.Token{Type: token.Int, Line: 1, Value: []byte{'1'}}, Value: 1},
				Operator: token.Token{Type: token.Mult, Line: 1},
			},
		},
		{
			name:     "hex int",
			input:    "0x12",
			expected: &ast.IntLitExpr{Token: token.Token{Type: token.Hex, Line: 1, Value: []byte("0x12")}, Value: 18},
		},
		{
			name:     "oct int",
			input:    "0o12",
			expected: &ast.IntLitExpr{Token: token.Token{Type: token.Oct, Line: 1, Value: []byte("0o12")}, Value: 10},
		},
		{
			name:     "int",
			input:    "12",
			expected: &ast.IntLitExpr{Token: token.Token{Type: token.Int, Line: 1, Value: []byte("12")}, Value: 12},
		},
		{
			name:     "float",
			input:    "12.0",
			expected: &ast.FloatLitExpr{Token: token.Token{Type: token.Float, Line: 1, Value: []byte("12.0")}, Value: 12.0},
		},
		{
			name:     "leading zero float",
			input:    "0.5",
			expected: &ast.FloatLitExpr{Token: token.Token{Type: token.Float, Line: 1, Value: []byte("0.5")}, Value: 0.5},
		},
		{
			name:     "no leading float",
			input:    ".5",
			expected: &ast.FloatLitExpr{Token: token.Token{Type: token.Float, Line: 1, Value: []byte(".5")}, Value: 0.5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New([]byte(tt.input))
			p, err := parser.New(l)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			expr, exprErr := p.Parse()
			if exprErr != tt.expectedErr {
				t.Errorf("expected error %v, got %v", tt.expectedErr, exprErr)
				return
			}
			if !reflect.DeepEqual(expr, tt.expected) {
				t.Errorf("expected\n%+v\ngot\n%+v", tt.expected, expr)
				return
			}
		})
	}
}
