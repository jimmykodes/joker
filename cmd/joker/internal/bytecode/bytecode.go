package bytecode

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jimmykodes/joker/compiler"
	"github.com/jimmykodes/joker/lexer"
	"github.com/jimmykodes/joker/parser"
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
			fmt.Println(bc)

		case ".jk":
			data, err := os.ReadFile(filename)
			if err != nil {
				return err
			}

			l := lexer.New(string(data))
			p := parser.New(l)
			prog := p.ParseProgram()
			c := compiler.New()
			if err := c.Compile(prog); err != nil {
				return err
			}
			fmt.Println(c.Bytecode())
		default:
			return fmt.Errorf("invalid filetype: %s", ext)
		}
		return nil
	}
}
