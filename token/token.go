package token

import (
	"fmt"
	"strconv"
)

type TokenType int

func (t TokenType) IsLiteral() bool {
	return literal_begin < t && t < literal_end
}

func (t TokenType) IsOperator() bool {
	return operators_begin < t && t < operators_end
}

func (t TokenType) IsKeyword() bool {
	return keywords_begin < t && t < keywords_end
}

func (t TokenType) Is(t2 TokenType) bool {
	return t == t2
}

func (t TokenType) String() string {
	return tokens[t]
}

const (
	ILLEGAL TokenType = iota
	EOF
	LF

	// cool idea, that i stole from go source code
	literal_begin
	IDENT // a
	INT   // 2
	FLOAT // 2.5
	BOOL  // true | false
	literal_end

	operators_begin
	ADD       // +
	SUB       // -
	MUL       // *
	DIV       // /
	REM       // %
	LPAREN    // (
	RPAREN    // )
	LBRACE    // {
	RBRACE    // }
	SEMICOLON // ;
	ASSIGN    // =
	EQ        // ==
	NOT       // !
	NEQ       // !=
	GT        // >
	GEQ       // >=
	LT        // <
	LEQ       // <=
	operators_end

	keywords_begin
	LET   // let
	IF    // if
	ELSE  // else
	TRUE  // true
	FALSE // false
	AND   // and
	OR    // or
	keywords_end
)

var tokens = [...]string{
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",
	LF:      "\\n",

	IDENT: "ident",
	INT:   "int",
	FLOAT: "float",
	BOOL:  "bool",

	ADD:       "+",
	SUB:       "-",
	MUL:       "*",
	DIV:       "/",
	REM:       "%",
	LPAREN:    "(",
	RPAREN:    ")",
	LBRACE:    "{",
	RBRACE:    "}",
	SEMICOLON: ";",
	ASSIGN:    "=",
	EQ:        "==",
	NOT:       "!",
	NEQ:       "!=",
	GT:        ">",
	GEQ:       ">=",
	LT:        "<",
	LEQ:       "<=",

	LET:   "let",
	IF:    "if",
	ELSE:  "else",
	TRUE:  "true",
	FALSE: "false",
	AND:   "and",
	OR:    "or",
}

// vim replace command for <TokenType> // <litetal> -> "<literal>": <TokenType>
// s/\(\w\+\)\s\+\/\/\s\+\(.\+\)/"\2": \1,

var keywords = map[string]TokenType{
	"let":   LET,
	"if":    IF,
	"else":  ELSE,
	"true":  TRUE,
	"false": FALSE,
	"and":   AND,
	"or":    OR,
}

// LookupKeyword lookup in keywords table
// and return the appropriate TokenType on found and IDENT otherwise
func LookupKeyword(literal string) TokenType {
	if keyword, ok := keywords[literal]; ok {
		return keyword
	}
	return IDENT
}

type Token struct {
	Type    TokenType
	Literal string
}

func (t Token) IsLiteral() bool {
	return literal_begin < t.Type && t.Type < literal_end
}

func (t Token) IsOperator() bool {
	return operators_begin < t.Type && t.Type < operators_end
}

func (t Token) IsKeyword() bool {
	return keywords_begin < t.Type && t.Type < keywords_end
}

func (t Token) Is(t2 TokenType) bool {
	return t.Type == t2
}

func (t Token) String() string {
	return fmt.Sprintf("[%s] %s", tokens[t.Type], t.Literal)
}

// TODO: use smaller numbers for precedence
// https://en.cppreference.com/w/c/language/operator_precedence
const (
	LOWEST  = 0
	UNARY   = 90
	HIGHEST = 100
)

func (t Token) Precedence() int {
	switch t.Type {
	case OR:
		return 1
	case AND:
		return 2
	case EQ, NEQ:
		return 3
	case GT, GEQ, LT, LEQ:
		return 4
	case ADD, SUB:
		return 5
	case MUL, DIV, REM:
		return 6
	case NOT:
		return UNARY
	default:
		return LOWEST
	}
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
