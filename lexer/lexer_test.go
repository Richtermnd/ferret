package lexer_test

import (
	"testing"

	"github.com/Richtermnd/ferret/lexer"
	"github.com/Richtermnd/ferret/token"
)

func TestOperandsRecognizing(t *testing.T) {
	source := "+-*/()"
	expected := []token.Token{
		token.NoLiteralToken(token.ADD),
		token.NoLiteralToken(token.SUB),
		token.NoLiteralToken(token.MUL),
		token.NoLiteralToken(token.DIV),
		token.NoLiteralToken(token.LPAREN),
		token.NoLiteralToken(token.RPAREN),
	}
	l := lexer.New(source)
	for i, expectedToken := range expected {
		tok := l.NextToken()
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
				Type:    token.FLOAT,
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
				Type:    token.FLOAT,
				Literal: "123456",
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
	source := "1 + (10 - 2) * 3.5 / 4"
	expected := []token.Token{
		{
			Type:    token.FLOAT,
			Literal: "1",
		},
		token.NoLiteralToken(token.ADD),
		token.NoLiteralToken(token.LPAREN),
		{
			Type:    token.FLOAT,
			Literal: "10",
		},
		token.NoLiteralToken(token.SUB),
		{
			Type:    token.FLOAT,
			Literal: "2",
		},
		token.NoLiteralToken(token.RPAREN),
		token.NoLiteralToken(token.MUL),
		{
			Type:    token.FLOAT,
			Literal: "3.5",
		},
		token.NoLiteralToken(token.DIV),
		{
			Type:    token.FLOAT,
			Literal: "4",
		},
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
