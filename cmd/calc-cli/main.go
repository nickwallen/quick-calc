package main

import (
	"bufio"
	"fmt"
	"github.com/nickwallen/toks"
	"io"
	"os"
)

func calculate(stdin io.Reader, writer io.Writer) {
	reader := bufio.NewReader(stdin)
	fmt.Fprintf(writer, "\n > ")
	input, _ := reader.ReadString('\n')

	// calculate the result
	result, err := toks.Calculate(input)
	if err != nil {
		fmt.Printf("%s \n", err)
		return
	}

	// output the result
	fmt.Fprintf(writer, "%s \n", result)
}

func main() {
	for {
		calculate(os.Stdin, os.Stdout)
	}
}
