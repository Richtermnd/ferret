package ast

import (
	"strings"
)

type Node interface {
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

func (p *Program) String() string {
	sb := strings.Builder{}
	for _, stmt := range p.Statements {
		sb.WriteString(stmt.String())
	}
	return sb.String()
}
