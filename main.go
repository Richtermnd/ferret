package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/Richtermnd/ferret/lexer"
	"github.com/Richtermnd/ferret/token"
)

const prompt = ">> "

func main() {
	s := bufio.NewScanner(os.Stdin)
	fmt.Print(prompt)
	for s.Scan() {
		l := lexer.New(s.Text())
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
