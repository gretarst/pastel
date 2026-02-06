package main

import (
	"fmt"
	"os"
	"pastel/interpreter"
	"pastel/lexer"
	"pastel/parser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: pastel <source-file>")
		os.Exit(1)
	}

	filename := os.Args[1]
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	input := string(data)

	// Step 1: Lexical analysis
	l := lexer.New(input)

	// Step 2: Parsing
	p := parser.New(l)
	prog := p.ParseProgram()

	// Step 3: Check for parsing errors
	if p.HasErrors() {
		fmt.Println("Parsing errors encountered:")
		for _, err := range p.Errors() {
			fmt.Println(err.Error())
		}
		return
	}

	// Step 4: Create interpreter and run the program
	interp := interpreter.New()
	if err := interp.Run(prog); err != nil {
		fmt.Println("Runtime error encountered:")
		fmt.Println(err.Error())
		return
	}

	// Step 5: Successful execution
	fmt.Println("Program executed successfully.")
}
