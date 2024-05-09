package joker

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/jimmykodes/joker/eval"
	"github.com/jimmykodes/joker/lexer"
	"github.com/jimmykodes/joker/parser"
)

type Joker struct{}

func (j *Joker) REPL() error {
	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := j.Read(reader)
		if err != nil {
			return err
		}
		if err := j.run(line); err != nil {
			return err
		}
	}
}

func (j *Joker) RunFile(file string) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	return j.run(data)
}

func (j *Joker) Read(reader *bufio.Reader) ([]byte, error) {
	fmt.Print(">> ")
	return reader.ReadBytes('\n')
}

func (j *Joker) run(code []byte) error {
	l := lexer.New(code)
	p, err := parser.New(l)
	if err != nil {
		return err
	}
	expr, err := p.Parse()
	var pErr parser.ParserError
	if errors.As(err, &pErr) {
		fmt.Println(pErr.Error())
		return nil
	} else if err != nil {
		return err
	}
	res := eval.Eval(expr)
	fmt.Println(res)
	return nil
}
