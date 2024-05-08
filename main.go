package main

import (
	"fmt"
	"os"

	"github.com/jimmykodes/joker/joker"
)

func main() {
	j := joker.Joker{}
	if err := j.REPL(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
