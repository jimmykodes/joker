package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jimmykodes/joker/cmd/joker/internal/build"
	"github.com/jimmykodes/joker/cmd/joker/internal/bytecode"
	"github.com/jimmykodes/joker/cmd/joker/internal/run"
	"github.com/jimmykodes/joker/repl"
)

func main() {
	if err := runner(); err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}

func runner() error {
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		repl.Start(os.Stdin, os.Stdout)
		return nil
	}

	var f func([]string) error
	switch cmd := args[0]; cmd {
	case "build":
		f = build.Cmd()
	case "run":
		f = run.Cmd()
	case "bytecode", "bc":
		f = bytecode.Cmd()
	}
	return f(args[1:])
}
