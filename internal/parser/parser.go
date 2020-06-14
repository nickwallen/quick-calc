package parser

import (
	"github.com/nickwallen/quick-calc/internal/tokens"
	"strconv"
)

// the interface used by the parser to read tokens
type tokenReader interface {
	// ReadToken reads the next token
	ReadToken() (tokens.Token, error)
}

// Parse a series of tokens and returns an expression.
func Parse(reader tokenReader) (Expression, error) {
	var expr Expression

	// an expression should start with a value like '23 pounds'
	value1, err := expectValue(reader)
	if err != nil {
		return expr, err
	}
	nextToken, err := reader.ReadToken()
	if err != nil {
		return expr, err
	}
	switch nextToken.TokenType {
	case tokens.Plus, tokens.Minus:
		return expectOperation(reader, value1, nextToken.TokenType)
	case tokens.In:
		return expectConversion(reader, value1)
	case tokens.EOF:
		return value1, nil
	default:
		return expr, ErrorUnexpectedToken(nextToken, tokens.Plus, tokens.Minus, tokens.In)
	}
}

func expectConversion(reader tokenReader, from Expression) (expr Expression, err error) {
	// expect the units to convert to
	units, err := expectUnits(reader)
	if err != nil {
		return expr, err
	}
	// expect EOF
	_, err = nextToken(reader, tokens.EOF)
	if err != nil {
		return expr, err
	}
	return conversion(from, units), nil
}

func expectOperation(reader tokenReader, value1 Expression, operator tokens.TokenType) (expr Expression, err error) {
	// to this point, we've already seen... operand1 +
	// now expect the second operand
	value2, err := expectValue(reader)
	if err != nil {
		return expr, err
	}
	token, err := reader.ReadToken()
	if err != nil {
		return expr, err
	}
	// what units should the result be in?
	switch token.TokenType {
	case tokens.Plus, tokens.Minus:
		// the operation has more operands
		right, err := expectOperation(reader, value2, token.TokenType)
		if err != nil {
			return expr, err
		}
		return binaryExpr(operator, value1, right, value1.TargetUnits), nil
	case tokens.In:
		// the units have been specified, for example '2 kg + 2 g in grams'
		units, err := expectUnits(reader)
		if err != nil {
			return expr, err
		}
		// expect EOF
		_, err = nextToken(reader, tokens.EOF)
		if err != nil {
			return expr, err
		}
		return binaryExpr(operator, value1, value2, units), nil
	case tokens.EOF:
		// operation complete; default to the units of the first operand, for example '2 kg + 2 g'
		return binaryExpr(operator, value1, value2, value1.TargetUnits), nil
	default:
		return expr, ErrorUnexpectedToken(token, tokens.Plus, tokens.Minus, tokens.In, tokens.EOF)
	}
}

func expectValue(reader tokenReader) (Expression, error) {
	var expr Expression
	token, err := nextToken(reader, tokens.Number)
	if err != nil {
		return expr, err
	}
	// TODO where to handle hexadecimal vs decimal?
	number, err := strconv.ParseFloat(token.Value, 64)
	if err != nil {
		return expr, err
	}
	units, err := expectUnits(reader)
	if err != nil {
		return expr, err
	}
	expr = valueExpr(number, units)
	return expr, nil
}

func expectUnits(reader tokenReader) (units Units, err error) {
	token, err := nextToken(reader, tokens.Units)
	if err != nil {
		return units, err
	}
	units, err = UnitsOf(token.Value)
	if err != nil {
		return units, ErrorInvalidUnits(token)
	}
	return units, nil
}

func nextToken(reader tokenReader, expected tokens.TokenType) (tokens.Token, error) {
	var nextToken tokens.Token
	nextToken, err := reader.ReadToken()
	if err != nil {
		return nextToken, ErrorReadFailed(err)
	}
	if nextToken.TokenType == tokens.Error {
		return nextToken, ErrorReadFailedNoCause()
	}
	if expected != nextToken.TokenType {
		if nextToken.TokenType == tokens.EOF {
			return nextToken, ErrorUnexpectedEOF(nextToken, expected)
		}
		return nextToken, ErrorUnexpectedToken(nextToken, expected)
	}
	return nextToken, nil
}
