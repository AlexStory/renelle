// main.go

package main

import (
	"fmt"
	"os"
	"renelle/repl"
)

func main() {
	fmt.Print("Welcome! To the renelle programming language!\n")
	fmt.Printf("Feel free to type in commands\n")
	repl.Start(os.Stdin, os.Stdout)
}
