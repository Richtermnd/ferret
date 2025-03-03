package ast

import "github.com/Richtermnd/ferret/token"

type FloatLiteral struct {
	Token token.Token
	Value float64
}

func (s *FloatLiteral) Literal() string { return s.Token.Literal }
func (s *FloatLiteral) String() string  { return s.Token.Literal }
func (s *FloatLiteral) exprNode()       {}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (s *IntegerLiteral) Literal() string { return s.Token.Literal }
func (s *IntegerLiteral) String() string  { return s.Token.Literal }
func (s *IntegerLiteral) exprNode()       {}

type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (s *BooleanLiteral) Literal() string { return s.Token.Literal }
func (s *BooleanLiteral) String() string  { return s.Token.Literal }
func (s *BooleanLiteral) exprNode()       {}

// TODO: string literal
