package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/ManuelGarciaF/go-interpreter/lexer"
	"github.com/ManuelGarciaF/go-interpreter/token"
)

const PROMPT = "> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	fmt.Fprintf(out, PROMPT)
	for scanner.Scan() {
		line := scanner.Text()
		l := lexer.New(line)

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Fprintf(out, "%+v\n", tok)
		}
		fmt.Fprintf(out, PROMPT)
	}
}
