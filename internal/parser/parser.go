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

	// Returns the original input that was tokenized.
	Input() string
}

// Parse a series of tokens and returns an expression.
func Parse(reader tokenReader) (types.Expression, types.InputError) {
	var expr types.Expression

	// an expression should start with a value like '23 pounds'
	value1, err := expectValue(reader)
	if err != nil {
		return expr, err
	}
	nextToken, readErr := reader.ReadToken()
	if readErr != nil {
		return expr, types.ErrorReadFailed(reader.Input(), readErr)
	}
	switch nextToken.TokenType {
	case types.Plus, types.Minus:
		return expectOperation(reader, value1, nextToken)
	case types.In:
		return expectConversion(reader, value1)
	case types.EOF:
		return value1, nil
	default:
		return expr, types.ErrorUnexpectedToken(reader.Input(), nextToken, types.Plus, types.Minus, types.In)
	}
}

func expectConversion(reader tokenReader, from types.Expression) (expr types.Expression, err types.InputError) {
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
	return types.UnitConversionExpr(from, units, reader.Input()), nil
}

func expectOperation(reader tokenReader, prevValue types.Expression, operator types.Token) (expr types.Expression, err types.InputError) {
	// expect the next operand in the expression
	nextValue, err := expectValue(reader)
	if err != nil {
		return expr, err
	}
	// we have an expression like prevValue + nextValue
	expr, err = operationExpr(operator, prevValue, nextValue, reader.Input())
	if err != nil {
		return expr, err
	}
	token, readErr := reader.ReadToken()
	if readErr != nil {
		return expr, types.ErrorReadFailed(reader.Input(), readErr)
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
		return expr, types.ErrorUnexpectedToken(reader.Input(), token, types.Plus, types.Minus, types.In, types.EOF)
	}
}

func expectValue(reader tokenReader) (expr types.Expression, err types.InputError) {
	token, err := nextToken(reader, types.Number)
	if err != nil {
		return expr, err
	}
	// TODO where to handle hexadecimal vs decimal?
	number, parseErr := strconv.ParseFloat(token.Value, 64)
	if parseErr != nil {
		return expr, types.ErrorInvalidNumber(reader.Input(), token)
	}
	units, err := expectUnits(reader)
	if err != nil {
		return expr, err
	}
	expr = types.NewValue(number, units)
	return expr, nil
}

func expectUnits(reader tokenReader) (units types.Token, err types.InputError) {
	token, err := nextToken(reader, types.Units)
	if err != nil {
		return units, err
	}
	// ensure that the units are valid
	_, unitErr := u.Find(token.Value)
	if unitErr != nil {
		return units, types.ErrorInvalidUnits(reader.Input(), token)
	}
	return token, nil
}

func nextToken(reader tokenReader, expected types.TokenType) (nextToken types.Token, err types.InputError) {
	nextToken, readErr := reader.ReadToken()
	if readErr != nil {
		return nextToken, types.ErrorReadFailed(reader.Input(), readErr)
	}
	if nextToken.TokenType == types.Error {
		return nextToken, types.ErrorTokenizerError(reader.Input(), nextToken)
	}
	if expected != nextToken.TokenType {
		if nextToken.TokenType == types.EOF {
			return nextToken, types.ErrorUnexpectedEOF(reader.Input(), nextToken, expected)
		}
		return nextToken, types.ErrorUnexpectedToken(reader.Input(), nextToken, expected)
	}
	return nextToken, nil
}

// operationExpr Create an expression where two values are acted on by an operator.
func operationExpr(operator types.Token, left types.Expression, right types.Expression, input string) (expr types.Expression, err types.InputError) {
	switch operator.TokenType {
	case types.Plus:
		return types.AdditionExpr(left, right, input), nil
	case types.Minus:
		return types.SubtractionExpr(left, right, input), nil
	default:
		return expr, types.ErrorInvalidOperator(input, operator)
	}
}
