package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jimmykodes/joker/cmd/joker/internal/build"
	"github.com/jimmykodes/joker/cmd/joker/internal/bytecode"
	"github.com/jimmykodes/joker/cmd/joker/internal/debugger"
	"github.com/jimmykodes/joker/cmd/joker/internal/interpreter"
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
	case "interpret", "i":
		f = interpreter.Cmd()
	case "debug", "d":
		f = debugger.Cmd()
	case "help", "h":
		fallthrough
	default:
		usage()
		return nil
	}
	return f(args[1:])
}

func usage() {
	fmt.Println("Usage:")
	fmt.Println("\njoker command file")
	fmt.Println("\nCommands:")
	fmt.Println("  build          build a .jkb file from a .jk file")
	fmt.Println("  run            run a .jk or .jkb file")
	fmt.Println("  debug, d       run a .jk or .jkb file using an interactive debugger")
	fmt.Println("  bytecode, bc   print the bytecode for a .jk file")
	fmt.Println("  interpret, i   run a .jk file using the interpreter instead of compiler")
	fmt.Println("  help, h        show this usage text")
}
