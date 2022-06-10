package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/jimmykodes/jk/evaluator"
	"github.com/jimmykodes/jk/lexer"
	"github.com/jimmykodes/jk/parser"
)

const Prompt = ">> "

func Start(in io.Reader, out io.Writer) {
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
		// fmt.Fprintln(out, prog.String())
		evaluated := evaluator.Eval(prog)
		if evaluated != nil {
			fmt.Fprintln(out, evaluated.Inspect())
		}
	}
}
