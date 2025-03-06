package ast

import (
	"strings"

	"github.com/Richtermnd/ferret/token"
)

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) exprNode()      {}
func (i *Identifier) String() string { return i.Value }

type ExpressionStatement struct {
	Token token.Token
	Expr  Expression
}

func (s *ExpressionStatement) String() string { return s.Expr.String() }
func (s *ExpressionStatement) stmtNode()      {}

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) exprNode() {}
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

func (pe *InfixExpression) exprNode() {}
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
