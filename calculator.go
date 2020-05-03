package calc

import (
	"fmt"
	"github.com/nickwallen/quick-calc/internal/io"
	"github.com/nickwallen/quick-calc/internal/parser"
	"github.com/nickwallen/quick-calc/internal/tokenizer"
)

// Calculate evaluates an input expression and returns the value as a string.
func Calculate(input string) (string, error) {
	var result string

	// evaluate the expression
	amount, err := CalculateAmount(input)
	if err != nil {
		return result, fmt.Errorf("execution error: %s", err.Error())
	}

	// output the result
	result = fmt.Sprintf("%.2f %s", amount.Value, amount.Units)
	return result, nil
}

// CalculateAmount evaluates an input expression and returns an Amount object.
func CalculateAmount(input string) (parser.Amount, error) {
	var amount parser.Amount

	// the tokenizer runs in the background populating the token channel
	tokens := io.NewTokenChannel()
	go tokenizer.Tokenize(input, tokens)

	// parse the tokens
	expression, err := parser.Parse(tokens)
	if err != nil {
		return amount, fmt.Errorf("parse error: %s", err.Error())
	}

	// evaluate the expression
	return expression.Evaluate()
}
