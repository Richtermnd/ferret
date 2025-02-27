package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/Richtermnd/ferret/evaluator"
	"github.com/Richtermnd/ferret/lexer"
	"github.com/Richtermnd/ferret/parser"
)

const prompt = ">> "

func interactive(in io.Reader, out io.Writer) {
	s := bufio.NewScanner(in)
	fmt.Print(prompt)
	for s.Scan() {
		l := lexer.New(s.Text())
		p := parser.New(l)
		program := p.Parse()
		if len(p.Errors()) != 0 {
			for _, err := range p.Errors() {
				fmt.Fprintln(out, err)
			}
			fmt.Fprint(out, prompt)
			continue
		}
		evaluated := evaluator.Eval(program)
		if evaluated != nil {
			fmt.Fprintln(out, evaluated.Inspect())
		} else {
			fmt.Fprintln(out, "failed to avaluate")
		}
		fmt.Fprint(out, prompt)
	}
}

func main() {
	if len(os.Args) >= 2 && os.Args[1] == "i" {
		interactive(os.Stdin, os.Stdout)
	}
}
