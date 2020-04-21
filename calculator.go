package toks

import (
	"fmt"
	"github.com/nickwallen/toks/internal/parser"
	"github.com/nickwallen/toks/internal/tokenizer"
	"github.com/nickwallen/toks/internal/util"
)

// Calculate Calculates the value of an input Expression.
func Calculate(input string) (string, error) {
	var result string
	var tokens util.TokenChannel

	// the tokenizer runs in the background populating the token channel
	tokens = make(chan tokenizer.Token, 2)
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
