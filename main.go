package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/Richtermnd/ferret/lexer"
	"github.com/Richtermnd/ferret/parser"
	"github.com/Richtermnd/ferret/token"
)

const prompt = ">> "

func interactive() {
	s := bufio.NewScanner(os.Stdin)
	fmt.Print(prompt)
	for s.Scan() {
		l := lexer.New(s.Text())
		p := parser.New(l)
		program := p.Parse()
		if len(p.Errors()) != 0 {
			for _, err := range p.Errors() {
				fmt.Println(err)
			}
		} else {
			for _, stmt := range program.Statements {
				fmt.Println(stmt)
			}
		}
		for {
			tok := l.NextToken()
			if tok.Type == token.EOF {
				break
			}
			fmt.Println(tok)
		}
		fmt.Print(prompt)
	}
}

func main() {
	if len(os.Args) >= 2 && os.Args[1] == "i" {
		interactive()
	}
}
