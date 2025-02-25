package token

import "strconv"

type TokenType int

const (
	ILLEGAL TokenType = iota
	EOF

	// cool idea, that i stole from go source code
	literal_begin
	INT   // 2
	FLOAT // 2.5
	literal_end

	operators_begin
	ADD    // +
	SUB    // -
	MUL    // *
	DIV    // /
	REM    // %
	LPAREN // (
	RPAREN // )
	operators_end
)

var tokens = [...]string{
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",

	INT:   "int",
	FLOAT: "float",

	ADD:    "+",
	SUB:    "-",
	MUL:    "*",
	DIV:    "/",
	REM:    "%",
	LPAREN: "(",
	RPAREN: ")",
}

type Token struct {
	Type    TokenType
	Literal string
}

func (t Token) isLiteral() bool {
	return literal_begin < t.Type && t.Type < literal_end
}

func (t Token) isOperator() bool {
	return operators_begin < t.Type && t.Type < operators_end
}

func (t Token) String() string {
	return t.Literal
}

func NoLiteralToken(t TokenType) Token {
	literal := ""
	if 0 <= t && t < TokenType(len(tokens)) {
		literal = tokens[t]
	}
	if literal == "" {
		literal = "unknown token: " + strconv.Itoa(int(t))
	}
	return Token{
		Type:    t,
		Literal: literal,
	}
}
