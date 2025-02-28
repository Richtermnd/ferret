package ast

import "github.com/Richtermnd/ferret/token"

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) Literal() string { return ls.Token.Literal }
func (ls *LetStatement) String() string  { return "let " + ls.Name.String() + " = " + ls.Value.String() }
func (ls *LetStatement) stmtNode()       {}
