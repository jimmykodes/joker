package main

import (
	"fmt"
	"os"

	"github.com/jimmykodes/joker/joker"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	j := joker.Joker{}
	if len(os.Args) > 1 {
		return j.RunFile(os.Args[1])
	} else {
		return j.REPL()
	}
}
