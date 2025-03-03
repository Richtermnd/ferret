package parser_test

import (
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

func TestIdentifier(t *testing.T) {
	testCases := []struct {
		desc   string
		source string
		token  token.Token
		value  string
	}{
		{
			desc:   "one char identifier",
			source: "a",
			token:  token.Token{Type: token.IDENT, Literal: "a"},
			value:  "a",
		},
		{
			desc:   "few chars identifier",
			source: "abc",
			token:  token.Token{Type: token.IDENT, Literal: "abc"},
			value:  "abc",
		},
		{
			desc:   "all valid chars identifier",
			source: "a1_b",
			token:  token.Token{Type: token.IDENT, Literal: "a1_b"},
			value:  "a1_b",
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

			ident, ok := stmt.Expr.(*ast.Identifier)
			if !ok {
				t.Fatalf("stmt exp not a *ast.Identifier: %T\n", stmt.Expr)
			}
			if ident.Token != tt.token {
				t.Errorf("mismatch tokens expected: %s got: %s\n", tt.token, ident.Token)
			}
			if ident.Value != tt.value {
				t.Errorf("mismatch values expected: %s got: %s\n", tt.value, ident.Value)
			}
		})
	}
}

func TestPrefixExpressions(t *testing.T) {
	testCases := []struct {
		desc     string
		source   string
		operator string
		output   string
	}{
		{
			desc:   "not",
			source: "!true",
			output: "(!true)",
		},
		{
			desc:   "unary minus",
			source: "-5",
			output: "(-5)",
		},
		{
			desc:   "not not",
			source: "!!true",
			output: "(!(!true))",
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
			if exp.String() != tt.output {
				t.Errorf("expected: %s got: %s\n", tt.output, exp.String())
			}
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
		{
			desc:     "eq",
			source:   "1 == 1",
			left:     1,
			operator: "==",
			right:    1,
		},
		{
			desc:     "neq",
			source:   "1 != 1",
			left:     1,
			operator: "!=",
			right:    1,
		},
		{
			desc:     "gt",
			source:   "1 > 1",
			left:     1,
			operator: ">",
			right:    1,
		},
		{
			desc:     "geq",
			source:   "1 >= 1",
			left:     1,
			operator: ">=",
			right:    1,
		},
		{
			desc:     "lt",
			source:   "1 < 1",
			left:     1,
			operator: "<",
			right:    1,
		},
		{
			desc:     "leq",
			source:   "1 <= 1",
			left:     1,
			operator: "<=",
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
		desc    string
		input   string
		output  string
		wantErr bool
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
		{
			desc:    "unclosed parenthesis",
			input:   "1 - (2 + 3",
			wantErr: true,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.desc, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			program := p.Parse()
			if tt.wantErr && p.HasErrors() {
				t.SkipNow()
			}
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
		wantErr    bool
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
		{
			desc:    "missed let",
			input:   "a = 5",
			wantErr: true,
		},
		{
			desc:    "missed identifier",
			input:   "let = 5",
			wantErr: true,
		},
		{
			desc:    "missed =",
			input:   "let a 5",
			wantErr: true,
		},
		{
			desc:    "missed expr",
			input:   "let a =",
			wantErr: true,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.desc, func(t *testing.T) {
			l := lexer.New(tt.input)
			p := parser.New(l)
			program := p.Parse()
			if tt.wantErr && p.HasErrors() {
				t.SkipNow()
			}
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
