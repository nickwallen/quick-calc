package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/nickwallen/toks"
)

func tokenize(stdin io.Reader, writer io.Writer) {
	reader := bufio.NewReader(stdin)
	fmt.Fprintf(writer, "\n > ")
	input, _ := reader.ReadString('\n')

	tok := toks.New(input)
	for token := range tok.Tokens() {
		fmt.Fprintf(writer, "%v  ", token)
	}
}

func main() {
	for {
		tokenize(os.Stdin, os.Stdout)
	}
}
