package ast

import (
	"fmt"
	"strings"

	"github.com/jimmykodes/jk/token"
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
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string {
	return i.Value
}

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (p *PrefixExpression) expressionNode()      {}
func (p *PrefixExpression) TokenLiteral() string { return p.Token.Literal }
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
func (e *InfixExpression) TokenLiteral() string { return e.Token.Literal }
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
func (i *IfExpression) TokenLiteral() string { return i.Token.Literal }
func (i *IfExpression) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "if %s {\n", i.Condition)
	fmt.Fprintf(&sb, "%s", i.Consequence.String())

	if i.Alternative != nil {
		fmt.Fprintf(&sb, "} else {\n%s}\n", i.Alternative.String())
	} else {
		sb.WriteString("}\n")
	}
	return sb.String()
}

type WhileExpression struct {
	Token     token.Token
	Condition Expression
	Body      *BlockStatement
}

func (w *WhileExpression) expressionNode()      {}
func (w *WhileExpression) TokenLiteral() string { return w.Token.Literal }
func (w *WhileExpression) String() string {
	var sb strings.Builder
	sb.WriteString("while (")
	sb.WriteString(w.Condition.String())
	sb.WriteString(") {\n")
	sb.WriteString(w.Body.String())
	sb.WriteString("};\n")
	return sb.String()
}

type CallExpression struct {
	Token     token.Token
	Function  Expression
	Arguments []Expression
}

func (c *CallExpression) expressionNode()      {}
func (c *CallExpression) TokenLiteral() string { return c.Token.Literal }
func (c *CallExpression) String() string {
	var sb strings.Builder
	sb.WriteString(c.Function.String() + c.Token.Literal)
	args := make([]string, len(c.Arguments))
	for i, arg := range c.Arguments {
		args[i] = arg.String()
	}
	sb.WriteString(strings.Join(args, ", ") + ")")
	return sb.String()
}

type IndexExpression struct {
	Token token.Token
	Left  Expression
	Index Expression
}

func (i *IndexExpression) expressionNode()      {}
func (i *IndexExpression) TokenLiteral() string { return i.Token.Literal }
func (i *IndexExpression) String() string {
	return "(" + i.Left.String() + "[" + i.Index.String() + "])"
}
