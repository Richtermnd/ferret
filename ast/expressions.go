package ast

import (
	"strings"

	"github.com/Richtermnd/ferret/token"
)

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) Literal() string { return i.Token.Literal }
func (i *Identifier) String() string  { return i.Value }
func (i *Identifier) exprNode()       {}

type ExpressionStatement struct {
	Token token.Token
	Expr  Expression
}

func (s *ExpressionStatement) Literal() string { return s.Token.Literal }
func (s *ExpressionStatement) String() string  { return s.Expr.String() }
func (s *ExpressionStatement) stmtNode()       {}

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) exprNode()       {}
func (pe *PrefixExpression) Literal() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out strings.Builder
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (pe *InfixExpression) exprNode()       {}
func (pe *InfixExpression) Literal() string { return pe.Token.Literal }
func (pe *InfixExpression) String() string {
	var out strings.Builder
	out.WriteString("(")
	out.WriteString(pe.Left.String())
	out.WriteString(" ")
	out.WriteString(pe.Operator)
	out.WriteString(" ")
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}
