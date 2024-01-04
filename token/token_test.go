package token

import (
	"testing"
)

func TestIdentType(t *testing.T) {
	tests := []struct {
		name  string
		ident string
		want  Type
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
			name:  "parses idents",
			ident: "my_var",
			want:  Ident,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IdentType(tt.ident); got != tt.want {
				t.Errorf("LookupIdent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumericType(t *testing.T) {
	tests := []struct {
		name  string
		ident string
		want  Type
	}{
		{
			name:  "parses int",
			ident: "123",
			want:  Int,
		},
		{
			name:  "parses float",
			ident: "123.5",
			want:  Float,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NumericType(tt.ident); got != tt.want {
				t.Errorf("NumericType() = %v, want %v", got, tt.want)
			}
		})
	}
}
