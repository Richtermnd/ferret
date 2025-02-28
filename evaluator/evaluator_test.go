package evaluator_test

import (
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
			source:   "1.23",
			expected: 1.23,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.desc, func(t *testing.T) {
			evaluated := testEval(t, tt.source)
			testFloatObject(t, evaluated, tt.expected)
		})
	}
}

func TestLetStatement(t *testing.T) {
	testCases := []struct {
		desc   string
		source string
		value  int64
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
		t.Errorf("Wrong value, expected: %d got: %d\n", intObj.Value, value)
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
