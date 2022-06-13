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
func (l *StringLiteral) TokenLiteral() string { return l.Token.Literal }
func (l *StringLiteral) String() string       { return fmt.Sprintf(`"%s"`, l.Value) }

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (l *IntegerLiteral) expressionNode()      {}
func (l *IntegerLiteral) TokenLiteral() string { return l.Token.Literal }
func (l *IntegerLiteral) String() string       { return l.Token.Literal }

type FloatLiteral struct {
	Token token.Token
	Value float64
}

func (l *FloatLiteral) expressionNode()      {}
func (l *FloatLiteral) TokenLiteral() string { return l.Token.Literal }
func (l *FloatLiteral) String() string       { return l.Token.Literal }

type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (l *BooleanLiteral) expressionNode()      {}
func (l *BooleanLiteral) TokenLiteral() string { return l.Token.Literal }
func (l *BooleanLiteral) String() string       { return strconv.FormatBool(l.Value) }

type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (f *FunctionLiteral) expressionNode()      {}
func (f *FunctionLiteral) TokenLiteral() string { return f.Token.Literal }
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
func (a *ArrayLiteral) TokenLiteral() string { return a.Token.Literal }
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
