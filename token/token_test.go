package token

import (
	"testing"
)

func TestIdentType(t *testing.T) {
	tests := []struct {
		name  string
		ident string
		want  Token
	}{
		{
			name:  "parses let",
			ident: "let",
			want:  Let,
		},
		{
			name:  "parses fn",
			ident: "fn",
			want:  Func,
		},
		{
			name:  "import",
			ident: "import",
			want:  Import,
		},
		{
			name:  "parses idents",
			ident: "my_var",
			want:  Ident,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Lookup(tt.ident); got != tt.want {
				t.Errorf("LookupIdent() = %v, want %v", got, tt.want)
			}
		})
	}
}
