package toks

import (
	"fmt"
	"strconv"
)

// parser Parses input text and outputs an Expression.
type parser struct {
	// the channel from which tokens can be read
	tokens chan Token
}

// Parse Parse a channel containing a series of tokens.
func Parse(tokens chan Token) (Expression, error) {
	var expression Expression

	// an expression should start with an amount like '23 pounds'
	parser := &parser{tokens: tokens}
	amount1, err := parser.expectAmount()
	if err != nil {
		return expression, err
	}

	// the next token defines if this is an operation or a conversion
	token := <-parser.tokens
	switch token.TokenType {
	case Plus, Minus, Multiply, Divide:
		return parser.expectOperation(amount1, token.TokenType)
	case In:
		return parser.expectConversion(amount1)
	case EOF:
		// all tokens have been consumed
		return amount1, nil
	default:
		// something bad happened because tokens remain that we not parsed
		return expression, fmt.Errorf("parsing error on input '%s'", token.Value)
	}
}

func (parser *parser) expectConversion(amount1 amount) (Expression, error) {
	var expression Expression

	// expect the units to convert to
	units, err := parser.expectUnits()
	if err != nil {
		return expression, err
	}

	// expect EOF
	_, err = parser.nextToken(EOF)
	if err != nil {
		return expression, err
	}

	// success
	return unitConverterOf(amount1, units), nil
}

func (parser *parser) expectOperation(amount1 amount, operator TokenType) (Expression, error) {
	// to this point, we've already seen... operand1 +
	var expression Expression

	// expect the second operand
	amount2, err := parser.expectAmount()
	if err != nil {
		return expression, err
	}

	// what units should the result be in?
	token := <-parser.tokens
	switch token.TokenType {
	case EOF:
		// success; default to units of the first operand for expressions like '2 kg + 2 g'
		return operatorOf(amount1, amount2, amount1.units, operator), nil

	case In:
		// the units have been specified for expressions like '2 kg + 2 g in grams'
		units, err := parser.expectUnits()
		if err != nil {
			return expression, err
		}

		// expect EOF
		_, err = parser.nextToken(EOF)
		if err != nil {
			return expression, err
		}

		// success
		return operatorOf(amount1, amount2, units, operator), nil

	default:
		return expression, fmt.Errorf("parsing error on unexpected input '%s'", token.Value)
	}
}

func (parser *parser) expectAmount() (amount, error) {
	var amount amount
	token, err := parser.nextToken(Number)
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
	expression := amountOf(number, units)
	return expression, nil
}

func (parser *parser) expectUnits() (units amountUnits, err error) {
	token, err := parser.nextToken(Units)
	if err != nil {
		return units, err
	}
	units, err = unitsOf(token.Value)
	if err != nil {
		return units, err
	}
	return units, nil
}

func (parser *parser) nextToken(expected TokenType) (Token, error) {
	token := <-parser.tokens
	if token.TokenType == Error {
		return token, fmt.Errorf(token.Value)
	}
	if expected != token.TokenType {
		value := fmt.Sprintf("got '%s'", token.Value)
		if token.TokenType == EOF {
			value = "reached end of input"
		}
		return token, fmt.Errorf("expected %s, but %s", expected, value)
	}
	return token, nil
}
