package ast

import (
	"strings"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var sb strings.Builder
	for _, statement := range p.Statements {
		if statement == nil {
			sb.WriteString("nil statement")
		} else {
			sb.WriteString(statement.String())
		}
		sb.WriteString("\n")
	}
	return sb.String()
}
