package calc

import (
	"fmt"
	"github.com/nickwallen/quick-calc/internal/io"
	"github.com/nickwallen/quick-calc/internal/parser"
	"github.com/nickwallen/quick-calc/internal/tokenizer"
	"github.com/nickwallen/quick-calc/internal/types"
)

// Calculate evaluates an input expression and returns the value as a string.
func Calculate(input string) (string, error) {
	var result string
	amt, err := CalculateAmount(input)
	if err != nil {
		return result, err
	}
	result = fmt.Sprintf("%.2f %s", amt.Value, amt.Units)
	return result, nil
}

// CalculateAmount evaluates an input expression and returns an Amount object.
func CalculateAmount(input string) (amt types.Amount, err error) {
	tokens := io.NewTokenChannel()
	go tokenizer.Tokenize(input, tokens)
	expr, err := parser.Parse(tokens)
	if err != nil {
		return amt, err
	}
	return expr.Eval()
}
