package parser

import (
	"testing"

	"github.com/jimmykodes/jk/lexer"
)

func TestParser_ParseProgram(t *testing.T) {

	tests := []struct {
		name          string
		input         string
		numStatements int
		programText   string
	}{
		{
			name:          "simple let",
			input:         "let x = 5;",
			numStatements: 1,
			programText:   "let x = 5;\n",
		},
		{
			name:          "simple let",
			input:         `let x = 5; let why = 12; let zed = 22;`,
			numStatements: 3,
			programText:   "let x = 5;\nlet why = 12;\nlet zed = 22;\n",
		},
		{
			name:          "let string",
			input:         `let x = "test";`,
			numStatements: 1,
			programText:   "let x = \"test\";\n",
		},
		{
			name:          "expression statement - ident",
			input:         "foobar;",
			numStatements: 1,
			programText:   "foobar",
		},
		{
			name:          "expression statement - int",
			input:         "-5;",
			numStatements: 1,
			programText:   "(-5)",
		},
		{
			name:          "expression statement - int",
			input:         "5",
			numStatements: 1,
			programText:   "5",
		},
		{
			name:          "expression statement - multiple",
			input:         "5 + 4 <= 2 + 12 * 2",
			numStatements: 1,
			programText:   "((5 + 4) <= (2 + (12 * 2)))",
		},
		{
			name:          "return int",
			input:         "return 43;",
			numStatements: 1,
			programText:   "return 43;\n",
		},
		{
			name:          "bool",
			input:         "true == false",
			numStatements: 1,
			programText:   "(true == false)",
		},
		{
			name:          "complex bool",
			input:         "5 < 3 == false",
			numStatements: 1,
			programText:   "((5 < 3) == false)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(lexer.New(tt.input))
			got := p.ParseProgram()
			if l := len(got.Statements); l != tt.numStatements {
				t.Errorf("incorrect number of statements returned: got %d - want %d", l, tt.numStatements)
			}
			for _, err := range p.errors {
				t.Errorf("parser error: %s", err)
			}
			if got.String() != tt.programText {
				t.Errorf("invalid program string. got %s - want %s", got, tt.programText)
			}
		})
	}
}
