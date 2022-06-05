package ast

import (
	"fmt"
	"strconv"

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
