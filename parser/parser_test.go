package parser

import (
	"fmt"
	"testing"

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
		{
			name:          "expression statement - ident",
			input:         "foobar;",
			numStatements: 1,
		},
		{
			name:          "expression statement - int",
			input:         "-5;",
			numStatements: 1,
		},
		{
			name:          "expression statement - int",
			input:         "-5 + 5;",
			numStatements: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(lexer.New(tt.input))
			got := p.ParseProgram()
			fmt.Println(got)
			if l := len(got.Statements); l != tt.numStatements {
				t.Errorf("incorrect number of statements returned: got %d - want %d", l, tt.numStatements)
			}
			for _, err := range p.errors {
				t.Errorf("parser error: %s", err)
			}
			// for _, stmt := range got.Statements {
			//
			// }
		})
	}
}
