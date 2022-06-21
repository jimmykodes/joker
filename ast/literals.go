package ast

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jimmykodes/jk/token"
)

type StringLiteral struct {
	Token token.Token
	Value string
}

func (l *StringLiteral) expressionNode()      {}
func (l *StringLiteral) TokenLiteral() string { return l.Token.String() }
func (l *StringLiteral) String() string       { return fmt.Sprintf(`"%s"`, l.Value) }

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (l *IntegerLiteral) expressionNode()      {}
func (l *IntegerLiteral) TokenLiteral() string { return l.Token.String() }
func (l *IntegerLiteral) String() string       { return strconv.FormatInt(l.Value, 10) }

type FloatLiteral struct {
	Token token.Token
	Value float64
}

func (l *FloatLiteral) expressionNode()      {}
func (l *FloatLiteral) TokenLiteral() string { return l.Token.String() }
func (l *FloatLiteral) String() string       { return strconv.FormatFloat(l.Value, 'g', -1, 64) }

type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (l *BooleanLiteral) expressionNode()      {}
func (l *BooleanLiteral) TokenLiteral() string { return l.Token.String() }
func (l *BooleanLiteral) String() string       { return strconv.FormatBool(l.Value) }

type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (f *FunctionLiteral) expressionNode()      {}
func (f *FunctionLiteral) TokenLiteral() string { return f.Token.String() }
func (f *FunctionLiteral) String() string {
	var sb strings.Builder
	sb.WriteString(f.TokenLiteral() + " (")
	params := make([]string, len(f.Parameters))
	for i, parameter := range f.Parameters {
		params[i] = parameter.String()
	}
	sb.WriteString(strings.Join(params, ", "))
	sb.WriteString(") {\n" + f.Body.String() + "}\n")
	return sb.String()
}

type ArrayLiteral struct {
	Token    token.Token
	Elements []Expression
}

func (a *ArrayLiteral) expressionNode()      {}
func (a *ArrayLiteral) TokenLiteral() string { return a.Token.String() }
func (a *ArrayLiteral) String() string {
	var sb strings.Builder
	sb.WriteString("[")
	elems := make([]string, len(a.Elements))
	for i, element := range a.Elements {
		elems[i] = element.String()
	}
	sb.WriteString(strings.Join(elems, ", "))
	sb.WriteString("]")
	return sb.String()
}

type MapLiteral struct {
	Token token.Token
	Pairs map[Expression]Expression
}

func (m *MapLiteral) expressionNode()      {}
func (m *MapLiteral) TokenLiteral() string { return m.Token.String() }
func (m *MapLiteral) String() string {
	var sb strings.Builder
	sb.WriteString("{")
	var final bool
	for key, val := range m.Pairs {
		final = true
		fmt.Fprintf(&sb, "\n\t%s: %s,", key, val)
	}
	if final {
		sb.WriteString("\n")
	}
	sb.WriteString("}")
	return sb.String()
}
