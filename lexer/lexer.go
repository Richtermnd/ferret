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
	case 0:
		tok = token.NoLiteralToken(token.EOF)
	case '+':
		tok = token.NoLiteralToken(token.ADD)
	case '-':
		tok = token.NoLiteralToken(token.SUB)
	case '*':
		tok = token.NoLiteralToken(token.MUL)
	case '/':
		tok = token.NoLiteralToken(token.DIV)
	case '%':
		tok = token.NoLiteralToken(token.REM)
	case '(':
		tok = token.NoLiteralToken(token.LPAREN)
	case ')':
		tok = token.NoLiteralToken(token.RPAREN)
	default:
		if isDigit(l.ch) {
			// TODO: for now it will recognize everything as float, add integer support
			literal := l.readNumber()
			if literal == "" {
				tok = token.NoLiteralToken(token.ILLEGAL)
			} else {
				tok.Type = token.FLOAT
				tok.Literal = literal
			}
		} else {
			tok = token.NoLiteralToken(token.ILLEGAL)
		}
	}

	return tok
}

func (l *Lexer) readChar() {
	if l.readpos >= len(l.source) {
		l.ch = 0
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
func (l *Lexer) readNumber() string {
	sb := strings.Builder{}
	for isDigit(l.ch) || l.ch == '_' || l.ch == '.' {
		if isDigit(l.ch) || l.ch == '.' {
			sb.WriteByte(l.ch)
		}
		l.readChar()
	}
	l.unreadChar()
	return sb.String()
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
