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
	case '\n':
		tok = newToken(token.LF, "\\n")
	case ';':
		tok = newToken(token.SEMICOLON, ";")
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
	case '{':
		tok = newToken(token.LBRACE, "(")
	case '}':
		tok = newToken(token.RBRACE, ")")
	case '=':
		tok = l.switchSuffix(token.ASSIGN, token.EQ, '=')
	case '!':
		tok = l.switchSuffix(token.NOT, token.NEQ, '=')
	case '>':
		tok = l.switchSuffix(token.GT, token.GEQ, '=')
	case '<':
		tok = l.switchSuffix(token.LT, token.LEQ, '=')
	default:
		if isLetter(l.ch) {
			literal := l.readIdentifier()
			t := token.LookupKeyword(literal)
			tok = newToken(t, literal)

		} else if isDigit(l.ch) {
			literal, t := l.readNumber()
			if literal == "" {
				tok = newToken(token.ILLEGAL, literal)
			} else {
				tok = newToken(t, literal)
			}

		} else {
			tok = newToken(token.ILLEGAL, string(l.ch))
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

func (l *Lexer) readIdentifier() string {
	startPos := l.pos
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' {
		l.readChar()
	}
	l.unreadChar()
	return l.source[startPos : l.pos+1]
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

	if isLetter(l.ch) {
		sb.WriteByte(l.ch)
		if l.ch == 'e' {
			l.readChar()
			exp, _ := l.readNumber()
			return sb.String() + exp, token.FLOAT
		} else {
			return sb.String(), token.ILLEGAL
		}
	}

	l.unreadChar()
	return sb.String(), t
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z'
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

func (l *Lexer) switchSuffix(t1, t2 token.TokenType, suffixCh byte) token.Token {
	var tok token.Token
	if peekChar := l.peekChar(); peekChar == suffixCh {
		tok = token.Token{Type: t2, Literal: string(l.ch) + string(peekChar)}
		l.readChar()
		return tok

	}
	return token.Token{Type: t1, Literal: string(l.ch)}
}
