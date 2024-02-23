package repl

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"renelle/evaluator"
	"renelle/lexer"
	"renelle/object"
	"renelle/parser"
)

func Start() {
	scanner := bufio.NewScanner(os.Stdin)
	env := object.NewEnvironment()

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

		evaluated := evaluator.Eval(program, env)

		if evaluated != nil {
			io.WriteString(os.Stdout, evaluated.Inspect())
			io.WriteString(os.Stdout, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []parser.ParseError) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg.Message+"\n")
	}
}
