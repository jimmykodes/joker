package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/jimmykodes/joker/compiler"
	"github.com/jimmykodes/joker/evaluator"
	"github.com/jimmykodes/joker/lexer"
	"github.com/jimmykodes/joker/object"
	"github.com/jimmykodes/joker/parser"
	"github.com/jimmykodes/joker/vm"
)

const Prompt = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()
	for {
		fmt.Fprint(out, Prompt)
		if scanned := scanner.Scan(); !scanned {
			return
		}

		l := lexer.New(scanner.Text())
		p := parser.New(l)
		prog := p.ParseProgram()
		if len(p.Errors()) > 0 {
			for _, err := range p.Errors() {
				fmt.Fprintf(out, "\t%s\n", err)
			}
			continue
		}
		evaluated := evaluator.Eval(prog, env)
		if evaluated != nil {
			fmt.Fprintln(out, evaluated.Inspect())
		}
	}
}

func StartVM(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	for {
		fmt.Fprint(out, Prompt)
		if scanned := scanner.Scan(); !scanned {
			return
		}

		l := lexer.New(scanner.Text())
		p := parser.New(l)
		prog := p.ParseProgram()
		if len(p.Errors()) > 0 {
			for _, err := range p.Errors() {
				fmt.Fprintf(out, "\t%s\n", err)
			}
			continue
		}
		comp := compiler.New()
		if err := comp.Compile(prog); err != nil {
			fmt.Fprintf(out, "Compile Error: %s\n", err)
			continue
		}
		machine := vm.New(comp.Bytecode())
		if err := machine.Run(); err != nil {
			fmt.Fprintf(out, "VM Error: %s\n", err)
			continue
		}
		stackTop := machine.StackTop()
		if stackTop != nil {
			fmt.Fprintln(out, stackTop.Inspect())
		}
	}
}
