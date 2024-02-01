package interpreter

import (
	"errors"
	"fmt"
	"os"

	"github.com/jimmykodes/joker/evaluator"
	"github.com/jimmykodes/joker/lexer"
	"github.com/jimmykodes/joker/object"
	"github.com/jimmykodes/joker/parser"
)

func Cmd() func(args []string) error {
	return func(args []string) error {
		filename := "main.jkb"
		if len(args) > 0 {
			filename = args[0]
			args = args[1:]
		}

		data, err := os.ReadFile(filename)
		if err != nil {
			return err
		}

		l := lexer.New(string(data))
		p := parser.New(l)
		prog := p.ParseProgram()

		if errs := p.Errors(); len(errs) > 0 {
			return errors.Join(errs...)
		}

		env := object.NewEnvironment()
		res := evaluator.Eval(prog, env)

		fmt.Println(res)

		if res.Type() == object.ErrorType {
			return fmt.Errorf(res.Inspect())
		}

		return nil
	}
}
