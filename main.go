package main

import (
	"fmt"
	"io"
	"os"

	"github.com/jimmykodes/joker/compiler"
	"github.com/jimmykodes/joker/lexer"
	"github.com/jimmykodes/joker/parser"
	"github.com/jimmykodes/joker/repl"
	"github.com/jimmykodes/joker/vm"
)

func main() {
	if len(os.Args) == 1 {
		// repl.Start(os.Stdin, os.Stdout)
		repl.StartVM(os.Stdin, os.Stdout)
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
		c := compiler.New()
		if err := c.Compile(prog); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		v := vm.New(c.Bytecode())
		if err := v.Run(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// env := object.NewEnvironment()
		// obj := evaluator.Eval(prog, env)
		// if obj.Type() == object.ErrorType {
		// 	fmt.Println(obj.Inspect())
		// }
		// for _, e := range p.Errors() {
		// 	fmt.Println(e)
		// }
	}
}
