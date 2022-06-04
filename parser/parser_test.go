package parser

import (
	"strings"
	"testing"

	"github.com/jimmykodes/jk/ast"
	"github.com/jimmykodes/jk/lexer"
)

func TestParser_ParseProgram(t *testing.T) {

	tests := []struct {
		name          string
		input         string
		numStatements int
		wantIdents    []string
	}{
		{
			name:          "simple let",
			input:         "let x = 5;",
			numStatements: 1,
			wantIdents:    []string{"x"},
		},
		{
			name:          "simple let",
			input:         `let x = 5; let why = 12; let zed = 22;`,
			numStatements: 3,
			wantIdents:    []string{"x", "why", "zed"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(lexer.New(tt.input))
			got := p.ParseProgram()
			if l := len(got.Statements); l != tt.numStatements {
				t.Errorf("incorrect number of statements returned: got %d - want %d", l, tt.numStatements)
			}
			for i, stmt := range got.Statements {
				s := stmt.(*ast.LetStatement)
				if !strings.EqualFold(s.Name.Value, tt.wantIdents[i]) {
					t.Errorf("incorrect identifier: got %s - want %s", s.Name.Value, tt.wantIdents[i])
				}
			}
		})
	}
}
