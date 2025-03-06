package ast

import (
	"fmt"
	"strings"

	"github.com/Richtermnd/ferret/token"
)

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (b *BlockStatement) String() string {
	sb := strings.Builder{}
	sb.WriteString("{")
	for _, stmt := range b.Statements {
		sb.WriteString(" ")
		sb.WriteString(stmt.String())
		sb.WriteString(";")
	}
	sb.WriteString(" }")
	return sb.String()
}

func (b *BlockStatement) stmtNode() {}

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) String() string { return "let " + ls.Name.String() + " = " + ls.Value.String() }
func (ls *LetStatement) stmtNode()      {}

type IfStatement struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative Statement
}

func (is *IfStatement) String() string {
	sb := strings.Builder{}
	fmt.Fprint(&sb, "if ")
	fmt.Fprint(&sb, is.Condition.String())
	fmt.Fprint(&sb, " ")
	fmt.Fprint(&sb, is.Consequence.String())
	if is.Alternative != nil {
		fmt.Fprint(&sb, " else ")
		fmt.Fprint(&sb, is.Alternative.String())
	}
	return sb.String()
}

func (is *IfStatement) stmtNode() {}
