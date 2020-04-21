package toks

import (
	"fmt"
)

// Calculate Calculates the value of an input Expression.
func Calculate(input string) (string, error) {
	var result string

	// the tokenizer runs in the background populating the tokens channel
	tokens := make(chan Token, 2)
	go Tokenize(input, tokens)

	// parse the input
	expr, err := Parse(tokens)
	if err != nil {
		return result, fmt.Errorf("parse error: %s", err.Error())
	}

	// evaluate the expression
	amount, err := expr.Evaluate()
	if err != nil {
		return result, fmt.Errorf("execution error: %s", err.Error())
	}

	// output the result
	return fmt.Sprintf("%.2f %s", amount.value, amount.units), nil
}
