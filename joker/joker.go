package joker

import (
	"bufio"
	"fmt"
	"os"

	"github.com/jimmykodes/joker/eval"
	"github.com/jimmykodes/joker/lexer"
	"github.com/jimmykodes/joker/parser"
)

type Joker struct{}

func (j *Joker) REPL() error {
	reader := bufio.NewReader(os.Stdin)
	env := eval.NewEnv(nil)
	for {
		line, err := j.Read(reader)
		if err != nil {
			return err
		}
		if err := j.run(line, env); err != nil {
			return err
		}
	}
}

func (j *Joker) RunFile(file string) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	return j.run(data, eval.NewEnv(nil))
}

func (j *Joker) Read(reader *bufio.Reader) ([]byte, error) {
	fmt.Print(">> ")
	return reader.ReadBytes('\n')
}

func (j *Joker) run(code []byte, env *eval.Env) error {
	l := lexer.New(code)
	p, err := parser.New(l)
	if err != nil {
		return err
	}
	expr, err := p.Parse()
	if err != nil {
		return err
	}
	res := eval.Eval(expr, env)
	fmt.Println(res)
	return nil
}
