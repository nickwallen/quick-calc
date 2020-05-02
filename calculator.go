package calc

import (
	"fmt"
	"github.com/nickwallen/quick-calc/internal/io"
	"github.com/nickwallen/quick-calc/internal/parser"
	"github.com/nickwallen/quick-calc/internal/tokenizer"
)

// Calculate Calculates the value of an input Expression.
func Calculate(input string) (string, error) {
	var result string

	// the tokenizer runs in the background populating the token channel
	tokens := io.NewTokenChannel()
	go tokenizer.Tokenize(input, tokens)

	// parse the tokens
	expression, err := parser.Parse(tokens)
	if err != nil {
		return result, fmt.Errorf("parse error: %s", err.Error())
	}

	// evaluate the expression
	amount, err := expression.Evaluate()
	if err != nil {
		return result, fmt.Errorf("execution error: %s", err.Error())
	}

	// output the result
	return fmt.Sprintf("%.2f %s", amount.Value, amount.Units), nil
}
