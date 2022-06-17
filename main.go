package main

import (
	"fmt"
	"io"
	"os"

	"github.com/jimmykodes/jk/evaluator"
	"github.com/jimmykodes/jk/lexer"
	"github.com/jimmykodes/jk/object"
	"github.com/jimmykodes/jk/parser"
	"github.com/jimmykodes/jk/repl"
)

func main() {
	if len(os.Args) == 1 {
		repl.Start(os.Stdin, os.Stdout)
	} else {
		file := os.Args[1]
		f, err := os.Open(file)
		if err != nil {
			fmt.Println("error opening file:", err)
			return
		}
		progText, err := io.ReadAll(f)
		if err != nil {
			fmt.Println("error reading file:", err)
			return
		}
		l := lexer.New(string(progText))
		p := parser.New(l)
		prog := p.ParseProgram()
		env := object.NewEnvironment()
		obj := evaluator.Eval(prog, env)
		if obj.Type() == object.ErrorType {
			fmt.Println(obj.Inspect())
		}
		for _, e := range p.Errors() {
			fmt.Println(e)
		}
	}
}
