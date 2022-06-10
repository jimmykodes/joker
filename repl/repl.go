package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/jimmykodes/jk/evaluator"
	"github.com/jimmykodes/jk/lexer"
	"github.com/jimmykodes/jk/object"
	"github.com/jimmykodes/jk/parser"
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
