package ast

import (
	"fmt"
	"strings"

	"github.com/jimmykodes/joker/token"
)

type Expression interface {
	Node
	expressionNode()
}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.String() }
func (i *Identifier) String() string {
	return i.Value
}

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (p *PrefixExpression) expressionNode()      {}
func (p *PrefixExpression) TokenLiteral() string { return p.Token.String() }
func (p *PrefixExpression) String() string {
	return fmt.Sprintf("(%s%s)", p.Operator, p.Right.String())
}

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (e *InfixExpression) expressionNode()      {}
func (e *InfixExpression) TokenLiteral() string { return e.Token.String() }
func (e *InfixExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", e.Left.String(), e.Operator, e.Right.String())
}

// todo: postfix Expression

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (i *IfExpression) expressionNode()      {}
func (i *IfExpression) TokenLiteral() string { return i.Token.String() }
func (i *IfExpression) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "if %s {\n", i.Condition)
	fmt.Fprintf(&sb, "%s", i.Consequence.String())

	if i.Alternative != nil {
		fmt.Fprintf(&sb, "} else {\n%s}", i.Alternative.String())
	} else {
		sb.WriteString("}")
	}
	return sb.String()
}

type WhileExpression struct {
	Token     token.Token
	Condition Expression
	Body      *BlockStatement
}

func (w *WhileExpression) expressionNode()      {}
func (w *WhileExpression) TokenLiteral() string { return w.Token.String() }
func (w *WhileExpression) String() string {
	var sb strings.Builder
	sb.WriteString("while (")
	sb.WriteString(w.Condition.String())
	sb.WriteString(") {\n")
	sb.WriteString(w.Body.String())
	sb.WriteString("};")
	return sb.String()
}

type ForExpression struct {
	Token     token.Token
	Init      Statement
	Condition Statement
	Increment Statement
	Body      *BlockStatement
}

func (f *ForExpression) expressionNode()      {}
func (f *ForExpression) TokenLiteral() string { return f.Token.String() }
func (f *ForExpression) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "for %s %s; %s {\n", f.Init, f.Condition, f.Increment)
	sb.WriteString(f.Body.String())
	sb.WriteString("}")
	return sb.String()
}

type CallExpression struct {
	Token     token.Token
	Function  Expression
	Arguments []Expression
}

func (c *CallExpression) expressionNode()      {}
func (c *CallExpression) TokenLiteral() string { return c.Token.String() }
func (c *CallExpression) String() string {
	var sb strings.Builder
	sb.WriteString(c.Function.String())
	sb.WriteRune('(')
	args := make([]string, len(c.Arguments))
	for i, arg := range c.Arguments {
		args[i] = arg.String()
	}
	sb.WriteString(strings.Join(args, ", "))
	sb.WriteString(");")
	return sb.String()
}

type IndexExpression struct {
	Token token.Token
	Left  Expression
	Index Expression
}

func (i *IndexExpression) expressionNode()      {}
func (i *IndexExpression) TokenLiteral() string { return i.Token.String() }
func (i *IndexExpression) String() string {
	return "(" + i.Left.String() + "[" + i.Index.String() + "])"
}
