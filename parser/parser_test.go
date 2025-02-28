package parser_test

import (
	"fmt"
	"testing"

	"github.com/Richtermnd/ferret/ast"
	"github.com/Richtermnd/ferret/lexer"
	"github.com/Richtermnd/ferret/parser"
	"github.com/Richtermnd/ferret/token"
)

func checkParserErrors(t *testing.T, p *parser.Parser) {
	errs := p.Errors()
	if !p.HasErrors() {
		return
	}
	for _, err := range errs {
		t.Error(err)
	}
	t.FailNow()
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}
	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}
	if integ.Literal() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d. got=%s", value, integ.Literal())
		return false
	}
	return true
}

func TestIntegerLiteralExpression(t *testing.T) {
	source := "1 23 4_5_6"
	expected := []ast.IntegerLiteral{
		{
			Token: token.Token{Type: token.INT, Literal: "1"},
			Value: 1,
		},
		{
			Token: token.Token{Type: token.INT, Literal: "23"},
			Value: 23,
		},
		{
			Token: token.Token{Type: token.INT, Literal: "456"},
			Value: 456,
		},
	}

	l := lexer.New(source)
	p := parser.New(l)
	program := p.Parse()
	checkParserErrors(t, p)
	if len(program.Statements) != len(expected) {
		t.Log(program.Statements)
		t.Fatalf("Expected num of expressions %d got %d\n", len(expected), len(program.Statements))
	}

	for i, stmt := range program.Statements {
		stmt, ok := stmt.(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("stmt not a ast.ExpressionStatement: %s", stmt.String())
		}
		literal, ok := stmt.Expr.(*ast.IntegerLiteral)
		if !ok {
			t.Fatalf("stmt exp not a ast.IntegerLiteral")
		}
		if literal.Token != expected[i].Token {
			t.Errorf("mismatch tokens, expected: %+v got: %+v\n", expected[i].Token, literal.Token)
		}
		if literal.Value != expected[i].Value {
			t.Errorf("mismatch tokens, expected: %+v got: %+v\n", expected[i].Token, literal.Token)
		}
	}
}

func TestFloatLiteralExpression(t *testing.T) {
	source := "1.0 2.3"
	expected := []ast.FloatLiteral{
		{
			Token: token.Token{Type: token.FLOAT, Literal: "1.0"},
			Value: 1.0,
		},
		{
			Token: token.Token{Type: token.FLOAT, Literal: "2.3"},
			Value: 2.3,
		},
	}

	l := lexer.New(source)
	p := parser.New(l)
	program := p.Parse()
	checkParserErrors(t, p)
	if len(program.Statements) != len(expected) {
		t.Log(program.Statements)
		t.Fatalf("Expected num of expressions %d got %d\n", len(expected), len(program.Statements))
	}

	for i, stmt := range program.Statements {
		stmt, ok := stmt.(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("stmt not a ast.ExpressionStatement: %s", stmt.String())
		}
		literal, ok := stmt.Expr.(*ast.FloatLiteral)
		if !ok {
			t.Fatalf("stmt exp not a ast.FloatLiteral")
		}
		if literal.Token != expected[i].Token {
			t.Errorf("mismatch tokens, expected: %+v got: %+v\n", expected[i].Token, literal.Token)
		}
		if literal.Value != expected[i].Value {
			t.Errorf("mismatch tokens, expected: %+v got: %+v\n", expected[i].Token, literal.Token)
		}
	}
}

func TestPrefixExpressions(t *testing.T) {
	testCases := []struct {
		desc     string
		source   string
		operator string
		value    int64
	}{
		// For future
		// {
		// 	desc:     "not",
		// 	source:   "!5",
		// 	operator: "!",
		// 	value:    5,
		// },
		{
			desc:     "unary minus",
			source:   "-5",
			operator: "-",
			value:    5,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.desc, func(t *testing.T) {
			l := lexer.New(tt.source)
			p := parser.New(l)
			program := p.Parse()
			checkParserErrors(t, p)
			t.Log(program.Statements)
			if len(program.Statements) != 1 {
				t.Fatalf("wrong number of statements, expected: 1 got: %d\n", len(program.Statements))
			}
			stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf("stmt not a ast.ExpressionStatement: %T\n", stmt)
			}

			exp, ok := stmt.Expr.(*ast.PrefixExpression)
			if !ok {
				t.Fatalf("stmt exp not a ast.PrefixExpression: %T\n", stmt.Expr)
			}
			if exp.Operator != tt.operator {
				t.Errorf("mismatch operator expected: %s got: %s\n", tt.operator, exp.Operator)
			}
			testIntegerLiteral(t, exp.Right, tt.value)
		})
	}
}

func TestInfixExpressions(t *testing.T) {
	testCases := []struct {
		desc     string
		source   string
		left     int64
		operator string
		right    int64
	}{
		{
			desc:     "add",
			source:   "1 + 1",
			left:     1,
			operator: "+",
			right:    1,
		},
		{
			desc:     "sub",
			source:   "1 - 1",
			left:     1,
			operator: "-",
			right:    1,
		},
		{
			desc:     "mul",
			source:   "1 * 1",
			left:     1,
			operator: "*",
			right:    1,
		},
		{
			desc:     "div",
			source:   "1 / 1",
			left:     1,
			operator: "/",
			right:    1,
		},
		{
			desc:     "rem",
			source:   "1 % 1",
			left:     1,
			operator: "%",
			right:    1,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.desc, func(t *testing.T) {
			l := lexer.New(tt.source)
			p := parser.New(l)
			program := p.Parse()
			checkParserErrors(t, p)
			t.Log(program.Statements)
			if len(program.Statements) != 1 {
				t.Fatalf("wrong number of statements, expected: 1 got: %d\n", len(program.Statements))
			}
			stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf("stmt not a ast.ExpressionStatement: %T\n", stmt)
			}

			exp, ok := stmt.Expr.(*ast.InfixExpression)
			if !ok {
				t.Fatalf("stmt exp not a ast.InfixExpression: %T\n", stmt.Expr)
			}
			if exp.Operator != tt.operator {
				t.Errorf("mismatch operator expected: %s got: %s\n", tt.operator, exp.Operator)
			}
			testIntegerLiteral(t, exp.Left, tt.left)
			testIntegerLiteral(t, exp.Right, tt.right)
		})
	}
}

func TestParenthesisExpressions(t *testing.T) {
	testCases := []struct {
		desc   string
		input  string
		output string
	}{
		{
			desc:   "no parenthesis",
			input:  "1 + 2 * 3",
			output: "(1+(2*3))",
		},
		{
			desc:   "no matter parenthesis",
			input:  "1 + (2 * 3)",
			output: "(1+(2*3))",
		},
		{
			desc:   "simple parenthesis",
			input:  "(1 + 2) * 3",
			output: "((1+2)*3)",
		},
		{
			desc:   "few parenthesis",
			input:  "(1 + 2) * (3 + 4)",
			output: "((1+2)*(3+4))",
		},
		{
			desc:   "nested parenthesis",
			input:  "(1 - (2 + 3)) * 4",
			output: "((1-(2+3))*4)",
		},
	}
	for _, tt := range testCases {
		t.Run(tt.desc, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			program := p.Parse()
			checkParserErrors(t, p)
			s := program.String()
			if tt.output != s {
				t.Errorf("expected: %s got: %s\n", tt.output, s)
			}
		})
	}
}

func TestLetStatements(t *testing.T) {
	testCases := []struct {
		desc       string
		input      string
		identifier string
		value      string
	}{
		{
			desc:       "simple assignment",
			input:      "let a = 5",
			identifier: "a",
			value:      "5",
		},
		{
			desc:       "complex assignment",
			input:      "let a = -5 * (1 + 2)",
			identifier: "a",
			value:      "((-5)*(1+2))",
		},
	}
	for _, tt := range testCases {
		t.Run(tt.desc, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			program := p.Parse()
			checkParserErrors(t, p)
			if len(program.Statements) != 1 {
				t.Log(program.Statements)
				t.Fatalf("Expected num of expressions %d got %d\n", 1, len(program.Statements))
			}
			testLetStatement(t, program.Statements[0], tt.identifier, tt.value)
		})
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name, value string) bool {
	if s.Literal() != "let" {
		t.Errorf("mismatch literals expected let got: %s", s.Literal())
		return false
	}
	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("mismatch statement type expected: *ast.LetStatement got: %T\n", s)
		return false
	}
	if letStmt.Name.Value != name {
		t.Errorf("mismatch identifier value expected: %s got: %s\n", name, letStmt.Name.Value)
		return false
	}
	if letStmt.Name.Literal() != name {
		t.Errorf("mismatch identifier literal expected: %s got: %s\n", name, letStmt.Name.Literal())
		return false
	}
	if letStmt.Value.String() != value {
		t.Errorf("mismatch value expected: %s got %s\n", value, letStmt.Value.String())
	}
	return true
}
