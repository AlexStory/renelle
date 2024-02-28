// main.go

package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"renelle/evaluator"
	"renelle/lexer"
	"renelle/object"
	"renelle/parser"
	"renelle/repl"
)

func main() {
	flag.Parse()

	args := flag.Args()

	if len(args) == 1 {
		filename := args[0]
		content, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file %s: %s\n", filename, err)
			os.Exit(1)
		}

		l := lexer.New(string(content))
		p := parser.New(l)
		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(os.Stderr, p.Errors())
			os.Exit(1)
		}

		env := object.NewEnvironment()
		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			fmt.Println(evaluated.Inspect())
		}
	} else if len(args) == 0 {
		repl.Start()
	} else {
		fmt.Println("Usage: renelle [file]")
		os.Exit(1)
	}
}

func printParserErrors(out io.Writer, errors []parser.ParseError) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg.Message+"\n")
	}
}
