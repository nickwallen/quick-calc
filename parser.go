package toks

import (
	"fmt"
	"strconv"
)

// Parser Parses input text and outputs an expression.
type Parser struct {
	tokenizer   *Tokenizer
	expressions chan Expression
}

// NewParser Creates a new parser.
func NewParser(input string) *Parser {
	parser := &Parser{tokenizer: NewTokenizer(input)}
	return parser
}

// Parse Parse the input text and return an expression
func (parser *Parser) Parse() (Expression, error) {
	var expression Expression
	amount1, err := parser.expectAmount()
	if err != nil {
		return expression, err
	}

	token := <-parser.tokenizer.Tokens()
	if token.TokenType != EOF {

		// expect an operation
		if token.TokenType == Plus {
			amount2, err := parser.expectAmount()
			if err != nil {
				return expression, err
			}
			return SumOf(amount1, amount2, amount1.Units), nil
		}
	} else {
		expression = amount1
	}

	return expression, nil
}

func (parser *Parser) expectAmount() (Amount, error) {
	var amount Amount

	token, err := parser.nextToken(Number)
	if err != nil {
		return amount, err
	}

	// TODO where to handle hexadecimal vs decimal?
	number, err := strconv.ParseInt(token.Value, 10, 64)
	if err != nil {
		return amount, err
	}

	units, err := parser.nextToken(Units)
	if err != nil {
		return amount, err
	}

	expression := AmountOf(number, UnitsOf(units.Value))
	return expression, nil
}

func (parser *Parser) nextToken(expected TokenType) (Token, error) {
	token := <-parser.tokenizer.Tokens()
	if token.TokenType == Error {
		return token, fmt.Errorf(token.Value)
	}
	if expected != token.TokenType {
		return token, fmt.Errorf("expected %s, but got %s", expected, token.TokenType)
	}
	return token, nil
}
