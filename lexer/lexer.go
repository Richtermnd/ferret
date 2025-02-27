package lexer

import (
	"strings"

	"github.com/Richtermnd/ferret/token"
)

type Lexer struct {
	source  string
	ch      byte
	peek    byte
	pos     int
	readpos int
	line    int
	col     int
}

func New(source string) *Lexer {
	l := &Lexer{
		source: source,
		line:   1,
		col:    0,
	}
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.readChar()
	l.skipWhitespaces()
	switch l.ch {
	case '\000':
		tok = newToken(token.EOF, string(l.ch))
	case '+':
		tok = newToken(token.ADD, "+")
	case '-':
		tok = newToken(token.SUB, "-")
	case '*':
		tok = newToken(token.MUL, "*")
	case '/':
		tok = newToken(token.DIV, "/")
	case '%':
		tok = newToken(token.REM, "%")
	case '(':
		tok = newToken(token.LPAREN, "(")
	case ')':
		tok = newToken(token.RPAREN, ")")
	case '\n':
		tok = newToken(token.LF, "\\n")
	default:
		if isDigit(l.ch) {
			literal, t := l.readNumber()
			if literal == "" {
				tok = newToken(token.ILLEGAL, "ILLEGAL")
			} else {
				tok.Type = t
				tok.Literal = literal
			}
		} else {
			tok = newToken(token.ILLEGAL, "ILLEGAL")
		}
	}

	return tok
}

func (l *Lexer) readChar() {
	if l.readpos >= len(l.source) {
		l.ch = '\000'
	} else {
		l.ch = l.source[l.readpos]
	}

	if l.ch == '\n' {
		l.line++
		l.col = 0
	} else {
		l.col++
	}

	l.pos = l.readpos
	l.readpos++
}

func (l *Lexer) unreadChar() {
	l.readpos--
	l.pos--
	l.ch = l.source[l.pos]
}

func (l *Lexer) peekChar() byte {
	if l.readpos >= len(l.source) {
		return 0
	} else {
		return l.source[l.readpos]
	}
}

func (l *Lexer) skipWhitespaces() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// readNumber read number, ignore '_' (python like syntax)
// TODO: check for invalid number
func (l *Lexer) readNumber() (string, token.TokenType) {
	sb := strings.Builder{}
	t := token.INT
	for isDigit(l.ch) || l.ch == '_' || l.ch == '.' {
		if isDigit(l.ch) || l.ch == '.' {
			sb.WriteByte(l.ch)
		}
		// ugly, but okay
		if l.ch == '.' {
			t = token.FLOAT
		}
		l.readChar()
	}
	l.unreadChar()
	return sb.String(), t
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func newToken(t token.TokenType, lit string) token.Token {
	return token.Token{
		Type:    t,
		Literal: lit,
	}
}
