package lexer_test

import (
	"testing"

	"github.com/Richtermnd/ferret/lexer"
	"github.com/Richtermnd/ferret/token"
)

func TestOperandsRecognizing(t *testing.T) {
	source := "+ - * / ( ) ; ="
	expected := []token.Token{
		{Type: token.ADD, Literal: "+"},
		{Type: token.SUB, Literal: "-"},
		{Type: token.MUL, Literal: "*"},
		{Type: token.DIV, Literal: "/"},
		{Type: token.LPAREN, Literal: "("},
		{Type: token.RPAREN, Literal: ")"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.ASSIGN, Literal: "="},
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
	source := "1 + (10 - 2) * 3.5 / 4 a foo"
	expected := []token.Token{
		{Type: token.INT, Literal: "1"},
		{Type: token.ADD, Literal: "+"},
		{Type: token.LPAREN, Literal: "("},
		{Type: token.INT, Literal: "10"},
		{Type: token.SUB, Literal: "-"},
		{Type: token.INT, Literal: "2"},
		{Type: token.RPAREN, Literal: ")"},
		{Type: token.MUL, Literal: "*"},
		{Type: token.FLOAT, Literal: "3.5"},
		{Type: token.DIV, Literal: "/"},
		{Type: token.INT, Literal: "4"},
		{Type: token.IDENT, Literal: "a"},
		{Type: token.IDENT, Literal: "foo"},
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
	source := "let"
	expected := []token.Token{
		{Type: token.LET, Literal: "let"},
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
