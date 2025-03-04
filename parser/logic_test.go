package parser_test

import (
	"fmt"
	"testing"

	"github.com/Richtermnd/ferret/lexer"
	"github.com/Richtermnd/ferret/parser"
)

type node interface {
	String() string
}

// ns node string
type ns string

func (s ns) String() string { return string(s) }

type and struct {
	left  node
	right node
}

func (op *and) String() string {
	return fmt.Sprintf("(%s %s %s)", op.left.String(), "and", op.right.String())
}

type or struct {
	left  node
	right node
}

func (op *or) String() string {
	return fmt.Sprintf("(%s %s %s)", op.left.String(), "or", op.right.String())
}

func TestLogicOperands(t *testing.T) {
	testCases := []struct {
		desc   string
		source string
		ouput  node
	}{
		{
			desc:   "and",
			source: "foo and bar",
			ouput: &and{
				left:  ns("foo"),
				right: ns("bar"),
			},
		},
		{
			desc:   "or",
			source: "foo or bar",
			ouput: &or{
				left:  ns("foo"),
				right: ns("bar"),
			},
		},
		{
			desc:   "operators precedence",
			source: "foo or bar and baz",
			ouput: &or{
				left: ns("foo"),
				right: &and{
					left:  ns("bar"),
					right: ns("baz"),
				},
			},
		},
		{
			desc:   "operators precedence with parenthesis",
			source: "(foo or bar) and baz",
			ouput: &and{
				left: &or{
					left:  ns("foo"),
					right: ns("bar"),
				},
				right: ns("baz"),
			},
		},
		{
			desc:   "or or",
			source: "foo or bar or baz",
			ouput: &or{
				left: &or{
					left:  ns("foo"),
					right: ns("bar"),
				},
				right: ns("baz"),
			},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.desc, func(t *testing.T) {
			t.Log("expected:", tt.ouput.String())
			p := parser.New(lexer.New(tt.source))
			program := p.Parse()
			checkParserErrors(t, p)
			if len(program.Statements) != 1 {
				t.Fatalf("wrong num of statements expected: %d got: %d", 1, len(program.Statements))
			}
			res, expected := program.Statements[0].String(), tt.ouput.String()
			if res != expected {
				t.Errorf("expected %s got %s", expected, res)
			}
		})
	}
}
