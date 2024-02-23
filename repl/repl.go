package repl

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"renelle/lexer"
	"renelle/parser"
)

func Start() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print(">> ")
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParserErrors(os.Stdout, p.Errors())
			continue
		}

		fmt.Println(program.String())
	}
}

func printParserErrors(out io.Writer, errors []parser.ParseError) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg.Message+"\n")
	}
}
