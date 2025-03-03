package parser_test

import (
	"testing"

	"github.com/Richtermnd/ferret/ast"
	"github.com/Richtermnd/ferret/lexer"
	"github.com/Richtermnd/ferret/parser"
	"github.com/Richtermnd/ferret/token"
)

func testIntegerLiteral(t *testing.T, expr ast.Expression, value int64) bool {
	intLit, ok := expr.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("expr not *ast.IntegerLiteral. got=%T", expr)
		return false
	}
	if intLit.Value != value {
		t.Errorf("Value not %d. got=%d", value, intLit.Value)
		return false
	}
	return true
}

func testFloatLiteral(t *testing.T, expr ast.Expression, value float64) bool {
	intLit, ok := expr.(*ast.FloatLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", expr)
		return false
	}
	if intLit.Value != value {
		t.Errorf("Value not %f. got=%f", value, intLit.Value)
		return false
	}
	return true
}

func testBooleanLiteral(t *testing.T, expr ast.Expression, value bool) bool {
	boolLit, ok := expr.(*ast.BooleanLiteral)
	if !ok {
		t.Errorf("not a *ast.BooleanLiteral: %T\n", expr)
		return false
	}
	if boolLit.Value != value {
		t.Errorf("mismatch values expected: %t got: %t\n", value, boolLit.Value)
		return false
	}
	return true
}

func TestIntegerLiteral(t *testing.T) {
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
		testIntegerLiteral(t, stmt.Expr, expected[i].Value)
	}
}

func TestFloatLiteral(t *testing.T) {
	source := "1.0 2.3 1.2e3"
	expected := []ast.FloatLiteral{
		{
			Token: token.Token{Type: token.FLOAT, Literal: "1.0"},
			Value: 1.0,
		},
		{
			Token: token.Token{Type: token.FLOAT, Literal: "2.3"},
			Value: 2.3,
		},
		{
			Token: token.Token{Type: token.FLOAT, Literal: "2.3"},
			Value: 1.2e3,
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
		testFloatLiteral(t, stmt.Expr, expected[i].Value)
	}
}

func TestBooleanLiteral(t *testing.T) {
	source := "true false"
	expected := []bool{true, false}

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
		testBooleanLiteral(t, stmt.Expr, expected[i])
	}
}
