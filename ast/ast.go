package ast

import (
	"strings"

	"github.com/Richtermnd/ferret/token"
)

type Node interface {
	Literal() string
	String() string
}

type Statement interface {
	Node
	stmtNode()
}

type Expression interface {
	Node
	exprNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) Literal() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].Literal()
	} else {
		return ""
	}
}

func (p *Program) String() string {
	sb := strings.Builder{}
	for _, stmt := range p.Statements {
		sb.WriteString(stmt.String())
	}
	return sb.String()
}

type ExpressionStatement struct {
	Token token.Token
	Expr  Expression
}

func (s *ExpressionStatement) Literal() string { return s.Token.Literal }
func (s *ExpressionStatement) String() string  { return s.Expr.String() }
func (s *ExpressionStatement) stmtNode()       {}

type FloatLiteral struct {
	Token token.Token
	Value float64
}

func (s *FloatLiteral) Literal() string { return s.Token.Literal }
func (s *FloatLiteral) String() string  { return s.Token.String() }
func (s *FloatLiteral) exprNode()       {}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (s *IntegerLiteral) Literal() string { return s.Token.Literal }
func (s *IntegerLiteral) String() string  { return s.Token.String() }
func (s *IntegerLiteral) exprNode()       {}

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
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}
