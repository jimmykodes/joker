package ast

import (
	"fmt"
	"strings"

	"github.com/jimmykodes/jk/token"
)

type Statement interface {
	Node
	statementNode()
}

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) String() string {
	return fmt.Sprintf("%s %s = %s;\n", ls.TokenLiteral(), ls.Name.Value, ls.Value)
}

type ReassignStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (rs *ReassignStatement) statementNode()       {}
func (rs *ReassignStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReassignStatement) String() string {
	return fmt.Sprintf("%s %s %s", rs.Name.Value, rs.Token.Literal, rs.Value.String())
}

type ReturnStatement struct {
	Token token.Token
	Value Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var sb strings.Builder
	sb.WriteString(rs.TokenLiteral() + " ")
	if rs.Value != nil {
		sb.WriteString(rs.Value.String())
	} else {
		sb.WriteString("<nil>")
	}
	sb.WriteString(";\n")
	return sb.String()
}

type ContinueStatement struct {
	Token token.Token
}

func (c *ContinueStatement) statementNode()       {}
func (c *ContinueStatement) TokenLiteral() string { return c.Token.Literal }
func (c *ContinueStatement) String() string       { return "continue" }

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string       { return es.Expression.String() }

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (b *BlockStatement) statementNode()       {}
func (b *BlockStatement) TokenLiteral() string { return b.Token.Literal }
func (b *BlockStatement) String() string {
	var sb strings.Builder
	for _, statement := range b.Statements {
		if statement == nil {
			sb.WriteString("nil statement")
		} else {
			sb.WriteString("\t" + statement.String())
		}
	}
	return sb.String()
}
