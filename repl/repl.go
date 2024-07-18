package repl

import (
	"fmt"
	"io"

	"github.com/ManuelGarciaF/go-interpreter/lexer"
	"github.com/ManuelGarciaF/go-interpreter/parser"

	"github.com/chzyer/readline"
)

func Start(in io.ReadCloser, out io.Writer) {
	rl, err := readline.NewEx(&readline.Config{
		Prompt: "> ",
		Stdin: in,
		Stdout: out,
	})
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	for {
		line, err := rl.Readline()
		if err != nil { // EOF or interrupt
			break
		}
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) == 0 {
			fmt.Fprintln(out, program.String())
		} else {
			for _, msg := range p.Errors() {
				fmt.Fprintf(out, "\t%s\n", msg)
			}
		}
	}
}
