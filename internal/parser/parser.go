package parser

import (
	u "github.com/bcicen/go-units"
	"github.com/nickwallen/quick-calc/internal/types"
	"strconv"
)

// the interface used by the parser to read tokens
type tokenReader interface {
	// ReadToken reads the next token
	ReadToken() (types.Token, error)
}

// Parse a series of tokens and returns an expression.
func Parse(reader tokenReader) (types.Expression, error) {
	var expr types.Expression

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
	case types.Plus, types.Minus:
		return expectOperation(reader, value1, nextToken)
	case types.In:
		return expectConversion(reader, value1)
	case types.EOF:
		return value1, nil
	default:
		return expr, errorUnexpectedToken(nextToken, types.Plus, types.Minus, types.In)
	}
}

func expectConversion(reader tokenReader, from types.Expression) (expr types.Expression, err error) {
	// expect the units to convert to
	units, err := expectUnits(reader)
	if err != nil {
		return expr, err
	}
	// expect EOF
	_, err = nextToken(reader, types.EOF)
	if err != nil {
		return expr, err
	}
	return types.UnitConversionExpr(from, units), nil
}

func expectOperation(reader tokenReader, prevValue types.Expression, operator types.Token) (expr types.Expression, err error) {
	// expect the next operand in the expression
	nextValue, err := expectValue(reader)
	if err != nil {
		return expr, err
	}
	// we have an expression like prevValue + nextValue
	expr, err = operationExpr(operator, prevValue, nextValue)
	if err != nil {
		return expr, err
	}
	token, err := reader.ReadToken()
	if err != nil {
		return expr, err
	}
	switch token.TokenType {
	case types.Plus, types.Minus:
		// the operation has more operands; prevValue + nextValue + ...
		return expectOperation(reader, expr, token)
	case types.In:
		return expectConversion(reader, expr)
	case types.EOF:
		return expr, nil
	default:
		return expr, errorUnexpectedToken(token, types.Plus, types.Minus, types.In, types.EOF)
	}
}

func expectValue(reader tokenReader) (expr types.Expression, err error) {
	token, err := nextToken(reader, types.Number)
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
	expr = types.NewValue(number, units)
	return expr, nil
}

func expectUnits(reader tokenReader) (units string, err error) {
	token, err := nextToken(reader, types.Units)
	if err != nil {
		return units, err
	}
	// ensure that the units are valid
	_, err = u.Find(token.Value)
	if err != nil {
		return units, errorInvalidUnits(token)
	}
	return token.Value, nil
}

func nextToken(reader tokenReader, expected types.TokenType) (nextToken types.Token, err error) {
	nextToken, err = reader.ReadToken()
	if err != nil {
		return nextToken, errorReadFailed(err)
	}
	if nextToken.TokenType == types.Error {
		return nextToken, errorTokenizerError(nextToken)
	}
	if expected != nextToken.TokenType {
		if nextToken.TokenType == types.EOF {
			return nextToken, errorUnexpectedEOF(nextToken, expected)
		}
		return nextToken, errorUnexpectedToken(nextToken, expected)
	}
	return nextToken, nil
}

// operationExpr Create an expression where two values are acted on by an operator.
func operationExpr(operator types.Token, left types.Expression, right types.Expression) (expr types.Expression, err error) {
	switch operator.TokenType {
	case types.Plus:
		return types.AdditionExpr(left, right), nil
	case types.Minus:
		return types.SubtractionExpr(left, right), nil
	default:
		return expr, errorInvalidOperator(operator)
	}
}
