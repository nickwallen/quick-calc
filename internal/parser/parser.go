package parser

import (
	"fmt"
	"github.com/nickwallen/toks/internal/tokenizer"
	"strconv"
)

// TODO no longer need to have this type; just pass around a tokenReader
// parser Parses input text and outputs an Expression.
type parser struct {
	// the channel from which tokens can be read
	reader tokenReader
}

// the interface used by the parser to read tokens
type tokenReader interface {
	ReadToken() (tokenizer.Token, error)
}

// Parse Parses a series of tokens and returns an expression.
func Parse(reader tokenReader) (Expression, error) {
	var expression Expression

	// an expression should start with an Amount like '23 pounds'
	parser := &parser{reader: reader}
	amount1, err := parser.expectAmount()
	if err != nil {
		return expression, err
	}

	token, err := reader.ReadToken()
	if err != nil {
		return expression, err
	}

	// the next token defines if this is an operation or a conversion
	switch token.TokenType {
	case tokenizer.Plus, tokenizer.Minus, tokenizer.Multiply, tokenizer.Divide:
		return parser.expectOperation(amount1, token.TokenType)

	case tokenizer.In:
		return parser.expectConversion(amount1)

	case tokenizer.EOF:
		// all tokens have been consumed, all we have is the first amount
		return amount1, nil

	default:
		// something bad happened because tokens remain that were not parsed
		return expression, fmt.Errorf("parsing error on input '%s'", token.Value)
	}
}

func (parser *parser) expectConversion(amount1 Amount) (Expression, error) {
	var expression Expression

	// expect the Units to convert to
	units, err := parser.expectUnits()
	if err != nil {
		return expression, err
	}

	// expect EOF
	_, err = parser.nextToken(tokenizer.EOF)
	if err != nil {
		return expression, err
	}

	// success
	return unitConverterOf(amount1, units), nil
}

func (parser *parser) expectOperation(amount1 Amount, operator tokenizer.TokenType) (Expression, error) {
	// to this point, we've already seen... operand1 +
	var expression Expression

	// expect the second operand
	amount2, err := parser.expectAmount()
	if err != nil {
		return expression, err
	}

	token, err := parser.reader.ReadToken()
	if err != nil {
		return expression, err
	}

	// what units should the result be in?
	switch token.TokenType {
	case tokenizer.EOF:
		// success; default to Units of the first operand for expressions like '2 kg + 2 g'
		return operatorOf(amount1, amount2, amount1.Units, operator), nil

	case tokenizer.In:
		// the Units have been specified for expressions like '2 kg + 2 g in grams'
		units, err := parser.expectUnits()
		if err != nil {
			return expression, err
		}

		// expect EOF
		_, err = parser.nextToken(tokenizer.EOF)
		if err != nil {
			return expression, err
		}

		// success
		return operatorOf(amount1, amount2, units, operator), nil

	default:
		return expression, fmt.Errorf("parsing error on unexpected input '%s'", token.Value)
	}
}

func (parser *parser) expectAmount() (Amount, error) {
	var amount Amount
	token, err := parser.nextToken(tokenizer.Number)
	if err != nil {
		return amount, err
	}
	// TODO where to handle hexadecimal vs decimal?
	number, err := strconv.ParseFloat(token.Value, 64)
	if err != nil {
		return amount, err
	}
	units, err := parser.expectUnits()
	if err != nil {
		return amount, err
	}
	expression := AmountOf(number, units)
	return expression, nil
}

func (parser *parser) expectUnits() (units Units, err error) {
	token, err := parser.nextToken(tokenizer.Units)
	if err != nil {
		return units, err
	}
	units, err = UnitsOf(token.Value)
	if err != nil {
		return units, err
	}
	return units, nil
}

func (parser *parser) nextToken(expected tokenizer.TokenType) (token tokenizer.Token, err error) {
	// get the next token
	token, err = parser.reader.ReadToken()
	if err != nil {
		return token, err
	}

	// if the tokenizer raises an error, its an error
	if token.TokenType == tokenizer.Error {
		return token, fmt.Errorf(token.Value)
	}

	// if the token is not the right type, its an error
	if expected != token.TokenType {
		value := fmt.Sprintf("got '%s'", token.Value)
		if token.TokenType == tokenizer.EOF {
			value = "reached end of input"
		}
		return token, fmt.Errorf("expected %s, but %s", expected, value)
	}
	return token, nil
}
