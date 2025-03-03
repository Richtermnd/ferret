package lexer_test

import (
	"testing"

	"github.com/Richtermnd/ferret/lexer"
	"github.com/Richtermnd/ferret/token"
)

func TestOperandsRecognizing(t *testing.T) {
	source := "+ - * / ( ) ; = == ! != > >= < <= $"
	expected := []token.Token{
		{Type: token.ADD, Literal: "+"},
		{Type: token.SUB, Literal: "-"},
		{Type: token.MUL, Literal: "*"},
		{Type: token.DIV, Literal: "/"},
		{Type: token.LPAREN, Literal: "("},
		{Type: token.RPAREN, Literal: ")"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.ASSIGN, Literal: "="},
		{Type: token.EQ, Literal: "=="},
		{Type: token.NOT, Literal: "!"},
		{Type: token.NEQ, Literal: "!="},
		{Type: token.GT, Literal: ">"},
		{Type: token.GEQ, Literal: ">="},
		{Type: token.LT, Literal: "<"},
		{Type: token.LEQ, Literal: "<="},
		{Type: token.ILLEGAL, Literal: "$"},
	}
	l := lexer.New(source)
	for i, expectedToken := range expected {
		tok := l.NextToken()
		t.Log(tok)
		if expectedToken != tok {
			t.Errorf("[%d] expected: %+v got: %+v\n", i, expectedToken, tok)
		}
	}
}

func TestNumbersRecognizing(t *testing.T) {
	testCases := []struct {
		desc     string
		source   string
		expected token.Token
	}{
		{
			desc:   "integer",
			source: "123",
			expected: token.Token{
				Type:    token.INT,
				Literal: "123",
			},
		},
		{
			desc:   "float",
			source: "123.456",
			expected: token.Token{
				Type:    token.FLOAT,
				Literal: "123.456",
			},
		},
		{
			desc:   "delimeted integer",
			source: "123_456",
			expected: token.Token{
				Type:    token.INT,
				Literal: "123456",
			},
		},
		{
			desc:   "invalid number",
			source: "1a2",
			expected: token.Token{
				Type:    token.ILLEGAL,
				Literal: "1a",
			},
		},
		{
			desc:   "scientific notation",
			source: "1.2e3",
			expected: token.Token{
				Type:    token.FLOAT,
				Literal: "1.2e3",
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			l := lexer.New(tC.source)
			tok := l.NextToken()
			if tok != tC.expected {
				t.Errorf("expected: %+v got: %+v\n", tC.expected, tok)
			}
		})
	}
}

func TestExpression(t *testing.T) {
	source := "-1 + (10 - a) * 3.5 / foo1_b !true false"
	expected := []token.Token{
		{Type: token.SUB, Literal: "-"},
		{Type: token.INT, Literal: "1"},
		{Type: token.ADD, Literal: "+"},
		{Type: token.LPAREN, Literal: "("},
		{Type: token.INT, Literal: "10"},
		{Type: token.SUB, Literal: "-"},
		{Type: token.IDENT, Literal: "a"},
		{Type: token.RPAREN, Literal: ")"},
		{Type: token.MUL, Literal: "*"},
		{Type: token.FLOAT, Literal: "3.5"},
		{Type: token.DIV, Literal: "/"},
		{Type: token.IDENT, Literal: "foo1_b"},
		{Type: token.NOT, Literal: "!"},
		{Type: token.TRUE, Literal: "true"},
		{Type: token.FALSE, Literal: "false"},
	}
	l := lexer.New(source)
	for i, expectedToken := range expected {
		tok := l.NextToken()
		t.Logf("%s\n", tok.Literal)
		if expectedToken != tok {
			t.Errorf("[%d] expected: %s got: %s\n", i, expectedToken.Literal, tok.Literal)
		}
	}
}

func TestKeywords(t *testing.T) {
	source := "let true false and or"
	expected := []token.Token{
		{Type: token.LET, Literal: "let"},
		{Type: token.TRUE, Literal: "true"},
		{Type: token.FALSE, Literal: "false"},
		{Type: token.AND, Literal: "and"},
		{Type: token.OR, Literal: "or"},
	}
	l := lexer.New(source)
	for i, expectedToken := range expected {
		tok := l.NextToken()
		t.Logf("%s\n", tok.Literal)
		if expectedToken.Type != tok.Type {
			t.Errorf("[%d] mismatch type expected: %d got: %d\n", i, expectedToken.Type, tok.Type)
		}
		if expectedToken != tok {
			t.Errorf("[%d] mismatch literals expected: %s got: %s\n", i, expectedToken.Literal, tok.Literal)
		}
	}
}
