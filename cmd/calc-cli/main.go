package main

import (
	"bufio"
	"fmt"
	"github.com/nickwallen/toks"
	"io"
	"os"
)

const (
	tokenMode = "tokens"
)

func tokenize(input string, writer io.Writer) {
	tok := toks.NewTokenizer(input)
	for token := range tok.Tokens() {
		fmt.Fprintf(writer, "%v  ", token)
	}
}

func calculate(input string, writer io.Writer) {
	result, err := toks.Calculate(input)
	if err != nil {
		fmt.Printf("%s \n", err)
		return
	}
	fmt.Fprintf(writer, "%s \n", result)
}

func prompt(stdin io.Reader, writer io.Writer, mode string) {
	// prompt for input
	reader := bufio.NewReader(stdin)
	fmt.Fprintf(writer, "\n > ")
	input, _ := reader.ReadString('\n')

	switch mode {
	case tokenMode:
		tokenize(input, writer)
	default:
		calculate(input, writer)
	}
}

func main() {
	mode := ""
	if len(os.Args) > 1 {
		mode = os.Args[1]
	}
	for {
		prompt(os.Stdin, os.Stdout, mode)
	}
}
