package main

import (
	"bufio"
	"fmt"
	"github.com/nickwallen/toks"
	"io"
	"os"
)

func parse(stdin io.Reader, writer io.Writer) {
	reader := bufio.NewReader(stdin)
	fmt.Fprintf(writer, "\n > ")
	input, _ := reader.ReadString('\n')

	// parse the input
	expr, err := toks.NewParser(input).Parse()
	if err != nil {
		fmt.Fprintf(writer, "parse error: %s \n", err.Error())
		return
	}

	// evaluate the expression
	amount, err := expr.Evaluate()
	if err != nil {
		fmt.Fprintf(writer, "Error: %s \n", err.Error())
		return
	}

	// output the result
	fmt.Fprintf(writer, "%d %s \n", amount.Value, amount.Units)
}

func main() {
	for {
		parse(os.Stdin, os.Stdout)
	}
}
