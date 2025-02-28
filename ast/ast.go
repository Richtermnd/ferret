package ast

import (
	"strings"
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
