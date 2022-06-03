package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/jimmykodes/jk/lexer"
	"github.com/jimmykodes/jk/token"
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
		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Fprintln(out, tok)
		}
	}
}
