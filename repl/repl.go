package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/butlermatt/monlox/lexer"
	"github.com/butlermatt/monlox/token"
)

const prompt = ">> "

// Start will begin a very simple REPL which reads from in and outputs response to out.
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprintf(out, prompt)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Fprintf(out, "%+v\n", tok)
		}
	}
}
