package main

import (
	"bufio"
	"fmt"
	calc "github.com/nickwallen/quick-calc"
	"github.com/nickwallen/quick-calc/internal/io"
	"github.com/nickwallen/quick-calc/internal/tokenizer"
	"os"
)

const (
	debugMode = "debug"
)

// used to read input from the user
type inputReader interface {
	ReadString(delimiter byte) (string, error)
}

// used to write output for the user
type outputWriter interface {
	Write(p []byte) (n int, err error)
}

// tokenize the input string.
func tokenize(input string, writer outputWriter) {
	output := io.NewTokenChannel()
	go tokenizer.Tokenize(input, output)
	for token := range output {
		fmt.Fprintf(writer, "%v  ", token)
	}
}

// calculate the value of the input.
func calculate(input string, writer outputWriter) {
	result, err := calc.Calculate(input)
	if err != nil {
		fmt.Printf("%s \n", err)
		return
	}
	fmt.Fprintf(writer, "%s \n", result)
}

// prompt the user for input
func prompt(reader inputReader, writer outputWriter, mode string) {
	// prompt for input
	fmt.Fprintf(writer, "\n > ")
	input, _ := reader.ReadString('\n')

	// start-up either the calculator or the tokenizer
	switch mode {
	case debugMode:
		// output just the tokens
		tokenize(input, writer)
	default:
		// evaluate each expression
		calculate(input, writer)
	}
}

func main() {
	mode := ""
	if len(os.Args) > 1 {
		mode = os.Args[1]
	}
	reader := bufio.NewReader(os.Stdin)
	for {
		prompt(reader, os.Stdout, mode)
	}
}
