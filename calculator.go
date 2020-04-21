package toks

import (
	"fmt"
)

// Calculate Calculates the value of an input Expression.
func Calculate(input string) (string, error) {
	var result string

	// Parse the input
	expr, err := NewParser(input).Parse()
	if err != nil {
		return result, fmt.Errorf("Parse error: %s", err.Error())
	}

	// Evaluate the Expression
	amount, err := expr.Evaluate()
	if err != nil {
		return result, fmt.Errorf("execution error: %s", err.Error())
	}

	// output the result
	return fmt.Sprintf("%.2f %s", amount.value, amount.units), nil
}
