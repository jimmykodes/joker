package run

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jimmykodes/joker/compiler"
	"github.com/jimmykodes/joker/lexer"
	"github.com/jimmykodes/joker/object"
	"github.com/jimmykodes/joker/parser"
	"github.com/jimmykodes/joker/vm"
)

func Cmd() func(args []string) error {
	return func(args []string) error {
		filename := "main.jkb"
		if len(args) > 0 {
			filename = args[0]
			args = args[1:]
		}
		switch ext := filepath.Ext(filename); ext {
		case ".jkb":
			data, err := os.ReadFile(filename)
			if err != nil {
				return err
			}

			var bc compiler.Bytecode
			if err := bc.UnmarshalBinary(data); err != nil {
				return err
			}

			machine := vm.New(&bc)
			return machine.Run()
		case ".jk":
			data, err := os.ReadFile(filename)
			if err != nil {
				return err
			}

			l := lexer.New(string(data))
			p := parser.New(l)
			prog := p.ParseProgram()
			if err := errors.Join(p.Errors()...); err != nil {
				return err
			}

			c := compiler.New()
			if err := c.Compile(prog); err != nil {
				return err
			}
			machine := vm.New(c.Bytecode())
			if err := machine.Run(); err != nil {
				return err
			}
			st := machine.StackTop()
			if st != nil && st.Type() == object.ErrorType {
				errOb, ok := st.(*object.Error)
				if ok {
					return fmt.Errorf("runtime error: %s", errOb)
				}
			}
			return nil

		default:
			return fmt.Errorf("invalid filetype: %s", ext)
		}
	}
}
