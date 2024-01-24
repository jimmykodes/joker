package build

import (
	"os"

	"github.com/jimmykodes/joker/compiler"
	"github.com/jimmykodes/joker/lexer"
	"github.com/jimmykodes/joker/parser"
)

func Cmd() func(args []string) error {
	return func(args []string) error {
		filename := "main.jk"
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
		c := compiler.New()

		if err := c.Compile(prog); err != nil {
			return err
		}

		outFilename := filename + "b"
		outFile, err := os.Create(outFilename)
		if err != nil {
			return err
		}
		defer outFile.Close()

		data, err = c.Bytecode().MarshalBinary()
		if err != nil {
			return err
		}

		if _, err := outFile.Write(data); err != nil {
			return err
		}

		return nil
	}
}
