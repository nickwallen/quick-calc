package toks

import (
	"fmt"
)

// Calculate Calculates the value of an input expression.
func Calculate(input string) (string, error) {
	var result string

	// parse the input
	expr, err := NewParser(input).Parse()
	if err != nil {
		return result, fmt.Errorf("parse error: %s", err.Error())
	}

	// evaluate the expression
	amount, err := expr.Evaluate()
	if err != nil {
		return result, fmt.Errorf("execution error: %s", err.Error())
	}

	// output the result
	return fmt.Sprintf("%.2f %s", amount.Value, amount.Units), nil
}
