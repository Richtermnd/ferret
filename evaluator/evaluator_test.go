package evaluator_test

import (
	"strings"
	"testing"

	"github.com/Richtermnd/ferret/evaluator"
	"github.com/Richtermnd/ferret/lexer"
	"github.com/Richtermnd/ferret/object"
	"github.com/Richtermnd/ferret/parser"
)

func TestEvalIntegerExpression(t *testing.T) {
	testCases := []struct {
		desc     string
		source   string
		expected int64
	}{
		{
			desc:     "simple",
			source:   "69",
			expected: 69,
		},
		{
			desc:     "separated",
			source:   "6_9",
			expected: 69,
		},
		{
			desc:     "sum",
			source:   "34 + 35",
			expected: 69,
		},
		{
			desc:     "sub",
			source:   "100 - 31",
			expected: 69,
		},
		{
			desc:     "mul",
			source:   "23 * 3",
			expected: 69,
		},
		{
			desc:     "div",
			source:   "4761 / 69",
			expected: 69,
		},
		{
			desc:     "few operands",
			source:   "35 + 17 * 2",
			expected: 69,
		},
		{
			desc:     "paranthesis",
			source:   "35 + (20 - 3) * 2",
			expected: 69,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.desc, func(t *testing.T) {
			evaluated := testEval(t, tt.source)
			testIntegerObject(t, evaluated, tt.expected)
		})
	}
}

func TestEvalFloatExpression(t *testing.T) {
	testCases := []struct {
		desc     string
		source   string
		expected float64
	}{
		{
			desc:     "simple",
			source:   "69.69",
			expected: 69.69,
		},
		{
			desc:     "scientific",
			source:   "0.6969e2",
			expected: 0.6969e2,
		},
		{
			desc:     "sum",
			source:   "34.5 + 34.5",
			expected: 69,
		},
		{
			desc:     "sub",
			source:   "71.5 - 2.5",
			expected: 69.0,
		},
		{
			desc:     "mul",
			source:   "23.23 * 3",
			expected: 69.69,
		},
		{
			desc:     "div",
			source:   "103.5 / 1.5",
			expected: 69,
		},
		{
			desc:     "few operands",
			source:   "34.5 + 17.25 * 2 + 0.69",
			expected: 69.69,
		},
		{
			desc:     "paranthesis",
			source:   "34.5 + (20 - 2.5) * 2 - 0.5",
			expected: 69.0,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.desc, func(t *testing.T) {
			evaluated := testEval(t, tt.source)
			testFloatObject(t, evaluated, tt.expected)
		})
	}
}

func TestEvalBoolExpression(t *testing.T) {
	const source = `
    true true == true 
    true false == false
    false true == false
    false true != true 
    false false != false
    true true != false
    true true or true
    true true and true
    false false and true
    true false or true

    true 1 == 1
    false 1 != 1
    false 1 == 0
    true 1 != 0
    true 1 > 0 
    true 1 >= 0
    false 1 < 0 
    false 1 <= 0
    true 1.0 == 1.0
    false 1.0 != 1.0
    false 1.0 == 0.0
    true 1.0 != 0.0
    true 1.0 > 0.0
    true 1.0 >= 0.0
    false 1.0 < 0.0
    false 1.0 <= 0.0

    true 1 == true
    false 1 != true
    false 1 == false
    true 1 != false
    true 1 > false
    true 1 >= false
    false 1 < false
    false 1 <= false

    true 1.0 == true 
    false 1.0 != true 
    false 1.0 == false
    true 1.0 != false
    true 1.0 > false
    true 1.0 >= false
    false 1.0 < false
    false 1.0 <= false`
	for _, line := range strings.Split(source, "\n") {
		t.Run(line, func(t *testing.T) {
			if line == "" {
				return
			}
			line := strings.TrimSpace(line)
			expectedLine, expr, _ := strings.Cut(line, " ")
			t.Log("expectedLine:", expectedLine)
			t.Log("expr:", expr)

			expected := testEval(t, expectedLine)
			res := testEval(t, expr)
			expBool, ok := expected.(*object.Bool)
			if !ok {
				t.Fatalf("wrong expected type: %T inspect: %s", expected, expected.Inspect())
			}
			resBool, ok := res.(*object.Bool)
			if !ok {
				t.Fatalf("wrong expr type: %T inspect: %s", res, res.Inspect())
			}
			if expBool.Value != resBool.Value {
				t.Errorf("wrong comprassion: '%s' expected: %t got %t\n", expr, expBool.Value, resBool.Value)
			}
		})
	}
}

func TestLetStatement(t *testing.T) {
	testCases := []struct {
		desc    string
		source  string
		value   int64
		wantErr bool
	}{
		{
			desc:   "simple let statement",
			source: "let a = 1; a",
			value:  1,
		},
		{
			desc:   "complex let statement",
			source: "let a = 10 - (3 + 1) * 2; a",
			value:  2,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.desc, func(t *testing.T) {
			testIntegerObject(t, testEval(t, tt.source), tt.value)
		})
	}
}

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

func testEval(t *testing.T, source string) object.Object {
	l := lexer.New(source)
	p := parser.New(l)
	program := p.Parse()
	checkParserErrors(t, p)
	env := object.NewEnv()
	return evaluator.Eval(env, program)
}

func testIntegerObject(t *testing.T, obj object.Object, value int64) bool {
	intObj, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("Not a object.Integer: %T\n", obj)
		return false
	}
	if intObj.Value != value {
		t.Errorf("Wrong value, expected: %d got: %d\n", value, intObj.Value)
		return false
	}
	return true
}

func testFloatObject(t *testing.T, obj object.Object, value float64) bool {
	floatObj, ok := obj.(*object.Float)
	if !ok {
		t.Errorf("Not a object.Float: %T\n", obj)
		return false
	}
	if floatObj.Value != value {
		t.Errorf("Wrong value, expected: %f got: %f\n", floatObj.Value, value)
		return false
	}
	return true
}
