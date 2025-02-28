package token

import "strconv"

type TokenType int

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
	literal_end

	operators_begin
	ADD       // +
	SUB       // -
	MUL       // *
	DIV       // /
	REM       // %
	LPAREN    // (
	RPAREN    // )
	ASSIGN    // =
	SEMICOLON // ;
	operators_end

	keywords_begin
	LET // let
	keywords_end
)

var tokens = [...]string{
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",
	LF:      "\\n",

	INT:   "int",
	FLOAT: "float",

	ADD:       "+",
	SUB:       "-",
	MUL:       "*",
	DIV:       "/",
	REM:       "%",
	LPAREN:    "(",
	RPAREN:    ")",
	ASSIGN:    "=",
	SEMICOLON: ";",
}

// vim replace command for <TokenType> // <litetal> -> "<literal>": <TokenType>
// s/\(\w\+\)\s\+\/\/\s\+\(.\+\)/"\2": \1,

var keywords = map[string]TokenType{
	"let": LET,
}

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
	return t.Literal
}

// TODO: use smaller numbers for precedence
const (
	LOWEST  = 0
	UNARY   = 90
	HIGHEST = 100
)

func (t Token) Precedence() int {
	switch t.Type {
	case ADD, SUB:
		return 10
	case MUL, DIV, REM:
		return 20
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
