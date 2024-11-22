// main.go

package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"renelle/evaluator"
	"renelle/lexer"
	"renelle/object"
	"renelle/parser"
	"renelle/repl"
	"strings"
)

func main() {
	flag.Parse()

	args := flag.Args()

	if len(args) >= 1 {
		switch args[0] {
		case "new":
			if len(args) < 2 {
				fmt.Println("Usage: renelle new <project_name>")
				os.Exit(1)
			}

			projectName := args[1]
			err := createProject(projectName)
			if err != nil {
				fmt.Println("Error creating project:", err)
				os.Exit(1)
			}

			fmt.Println("Project created successfully:", projectName)
		case "run":
			dir, err := findProjectDir()
			if err != nil {
				fmt.Println("Error finding project directory:", err)
				os.Exit(1)
			}

			moduleName, err := getModuleName(dir, args)
			if err != nil {
				fmt.Println("Error getting module name:", err)
				os.Exit(1)
			}

			filename := filepath.Join(dir, "src", "main.rnl")
			runFile(filename, moduleName, args[1:])
		case "test":
			fmt.Printf("Test command not implemented yet\n")
			var dir string
			if len(args) > 1 {
				dir = args[1]
			} else {
				dir = "./test"
			}

			err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if !info.IsDir() && strings.HasSuffix(path, "_test.rnl") {
					// Run the tests in the file
					fmt.Printf("Running tests in %s\n", path)
					// TODO: Implement the function to run the tests
					// runTests(path)
				}

				return nil
			})

			if err != nil {
				fmt.Printf("Error walking the path %v: %v\n", dir, err)
				os.Exit(1)
			}
		default:
			filename := args[0]
			content, err := os.ReadFile(filename)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading file %s: %s\n", filename, err)
				os.Exit(1)
			}

			l := lexer.New(string(content), filename)
			p := parser.New(l)
			program := p.ParseProgram()
			if len(p.Errors()) != 0 {
				printParserErrors(os.Stderr, p.Errors())
				os.Exit(1)
			}

			env := object.NewEnvironment()
			ctx := object.NewEvalContext()
			(*ctx.MetaData)["args"] = args[1:]

			evaluator.Eval(program, env, ctx)
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

func findProjectDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "rnl.rnl")); err == nil {
			return dir, nil
		}

		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			return "", fmt.Errorf("no rnl.rnl file found")
		}

		dir = parentDir
	}
}

func runFile(filename string, moduleName string, args []string) {
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file %s: %s\n", filename, err)
		os.Exit(1)
	}

	l := lexer.New(string(content), filename)
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		printParserErrors(os.Stderr, p.Errors())
		os.Exit(1)
	}

	env := object.NewEnvironment()
	ctx := object.NewEvalContext()
	(*ctx.MetaData)["args"] = args

	evaluator.Eval(program, env, ctx)
	module, ok := env.GetModule(moduleName)
	if !ok {
		fmt.Println("Module not found:", moduleName)
		os.Exit(1)
	}

	// Retrieve the main function from the module
	mainFunc, ok := module.Environment.Get("main")
	if !ok {
		fmt.Println("main function not found")
		os.Exit(1)
	}

	ret := evaluator.ApplyFunction(mainFunc, []object.Object{}, ctx)
	if e, ok := ret.(*object.Error); ok {
		fmt.Println(e.Message)
		os.Exit(1)
	}

}

func getModuleName(dir string, args []string) (string, error) {
	// Evaluate the rnl.rnl file to get the module name
	rnlFilename := filepath.Join(dir, "rnl.rnl")
	rnlContent, err := os.ReadFile(rnlFilename)
	if err != nil {
		return "", err
	}

	l := lexer.New(string(rnlContent), rnlFilename)
	p := parser.New(l)
	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		printParserErrors(os.Stderr, p.Errors())
		os.Exit(1)
	}

	env := object.NewEnvironment()
	ctx := object.NewEvalContext()
	(*ctx.MetaData)["args"] = args

	evaluator.Eval(program, env, ctx)

	// Get the properties from the environment
	properties, ok := env.Get("properties")
	if !ok {
		return "", errors.New("properties not found in rnl.rnl file")
	}

	// Get the module name from the properties
	mstring := &object.Atom{Value: "moduleName"}
	moduleName, ok := properties.(*object.Map).Get(mstring)
	if !ok {
		return "", errors.New("module name not found in properties")
	}

	moduleString := moduleName.(*object.String).Value
	// Attach the module name to the metadata
	(*ctx.MetaData)["moduleName"] = moduleString

	return moduleString, nil
}
