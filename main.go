package main

import (
	"os"

	"github.com/ManuelGarciaF/go-interpreter/repl"
)

func main() {
	repl.Start(os.Stdin, os.Stdout)
}
