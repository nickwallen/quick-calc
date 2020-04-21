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
	switch token.TokenType {
	case Plus, Minus, Multiply, Divide:
		return parser.expectOperation(amount1, token.TokenType)
	case In:
		return parser.expectConversion(amount1)
	case EOF:
		// all tokens have been consumed
		return amount1, nil
	default:
		// tokens remain, something bad happened
		return expression, fmt.Errorf("parsing error on unexpected input '%s'", token.Value)
	}
}

func (parser *Parser) expectConversion(amount1 Amount) (Expression, error) {
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
	return UnitConverterOf(amount1, units), nil
}

func (parser *Parser) expectOperation(amount1 Amount, operator TokenType) (Expression, error) {
	// to this point, we've already seen... operand1 +
	var expression Expression

	// expect the second operand
	amount2, err := parser.expectAmount()
	if err != nil {
		return expression, err
	}

	// what units should the result be in?
	token := <-parser.tokenizer.Tokens()
	switch token.TokenType {
	case EOF:
		// success; default to units of the first operand for expressions like '2 kg + 2 g'
		return OperatorOf(amount1, amount2, amount1.Units, operator), nil

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
		return OperatorOf(amount1, amount2, units, operator), nil

	default:
		return expression, fmt.Errorf("parsing error on unexpected input '%s'", token.Value)
	}
}

func (parser *Parser) expectAmount() (Amount, error) {
	var amount Amount
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
	expression := AmountOf(number, units)
	return expression, nil
}

func (parser *Parser) expectUnits() (units AmountUnits, err error) {
	token, err := parser.nextToken(Units)
	if err != nil {
		return units, err
	}
	units, err = UnitsOf(token.Value)
	if err != nil {
		return units, err
	}
	return units, nil
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
