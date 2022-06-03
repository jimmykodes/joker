package main

import (
	"os"

	"github.com/jimmykodes/jk/repl"
)

func main() {
	repl.Start(os.Stdin, os.Stdout)
}
