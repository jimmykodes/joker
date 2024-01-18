package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/jimmykodes/joker/ast"
	"github.com/jimmykodes/joker/compiler"
	"github.com/jimmykodes/joker/evaluator"
	"github.com/jimmykodes/joker/lexer"
	"github.com/jimmykodes/joker/object"
	"github.com/jimmykodes/joker/parser"
	"github.com/jimmykodes/joker/repl"
	"github.com/jimmykodes/joker/vm"
)

var (
	compile = flag.Bool("compile", false, "use the compiler instead of the interpreter")
	dump    = flag.Bool("dump", false, "if compiling, dump the constants and instructions")
)

func main() {
	flag.Parse()
	if len(flag.Args()) == 0 {
		repl.Start(os.Stdin, os.Stdout)
		// repl.StartVM(os.Stdin, os.Stdout)
	} else {
		file := flag.Arg(0)
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

		var runner func(ast.Node)
		if *compile {
			runner = compiled()
		} else {
			runner = interpreted(p)
		}

		runner(prog)
	}
}

func compiled() func(ast.Node) {
	return func(prog ast.Node) {
		c := compiler.New()
		if err := c.Compile(prog); err != nil {
			fmt.Println("compiler error:", err)
			os.Exit(1)
		}
		bc := c.Bytecode()

		if *dump {
			fmt.Println(bc)
		}

		v := vm.New(bc)
		if err := v.Run(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if *dump {
			fmt.Println(v.LastPoppedStackElem())
		}
	}
}

func interpreted(p *parser.Parser) func(ast.Node) {
	return func(prog ast.Node) {
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
