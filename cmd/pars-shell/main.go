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

	expr := toks.NewParser(input).Parse()
	amount, err := expr.Evaluate()
	if err == nil {
		fmt.Fprintf(writer, "%d %s \n", amount.Value, amount.Units)
	} else {
		fmt.Fprintf(writer, "Error: %s \n", err.Error())
	}

}

func main() {
	for {
		parse(os.Stdin, os.Stdout)
	}
}
