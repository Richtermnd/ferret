package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/Richtermnd/ferret/evaluator"
	"github.com/Richtermnd/ferret/lexer"
	"github.com/Richtermnd/ferret/object"
	"github.com/Richtermnd/ferret/parser"
)

func init() {
	flag.Parse()
}

func eval(env *object.Environment, source string) object.Object {
	l := lexer.New(source)
	p := parser.New(l)
	program := p.Parse()
	if p.HasErrors() {
		p.PrintErrors(os.Stderr)
		return nil
	}
	fmt.Printf("program.Statements: %v\n", program.Statements)
	return evaluator.Eval(env, program)
}

func repl() {
	const prompt = ">> "
	s := bufio.NewScanner(os.Stdin)
	env := object.NewEnv()
	fmt.Print(prompt)
	for s.Scan() {
		evaluated := eval(env, s.Text())
		if evaluated != nil {
			fmt.Println(evaluated.Inspect())
		}
		fmt.Print(prompt)
	}
}

func fatalf(msg string, args ...any) {
	fmt.Fprintf(os.Stderr, msg, args...)
	os.Exit(1)
}

func main() {
	if flag.NArg() == 0 {
		repl()
	} else {
		source, err := os.ReadFile(flag.Arg(0))
		if err != nil {
			fatalf("failed to read %s: %v\n", flag.Arg(0), err)
		}
		env := object.NewEnv()
		eval(env, string(source))
	}
}
