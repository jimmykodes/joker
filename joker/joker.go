package joker

import (
	"bufio"
	"fmt"
	"os"
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
	fmt.Println(string(code))
	return nil
}
