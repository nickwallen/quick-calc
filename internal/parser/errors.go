package parser

import (
	"fmt"
	"github.com/nickwallen/quick-calc/internal/types"
	"strings"
)

// unexpectedToken is an error that occurs when an unexpected token is found.
type unexpectedToken struct {
	Expected []types.TokenType // the type of token(s) that are expected
	BadToken types.Token       // the token that we got
	Position int               // the position of the error
	Width    int               // the width of the error
}

// errorUnexpectedToken creates a new unexpected token error.
func errorUnexpectedToken(badToken types.Token, expected ...types.TokenType) *unexpectedToken {
	return &unexpectedToken{
		Expected: expected,
		BadToken: badToken,
		Position: badToken.Position,
		Width:    len(badToken.Value),
	}
}

func (u *unexpectedToken) Error() string {
	var expected []string
	for _, e := range u.Expected {
		expected = append(expected, fmt.Sprintf("%s", e))
	}
	return fmt.Sprintf("at position %d, got '%s', but expected %s", u.Position, u.BadToken.Value, strings.Join(expected, ", "))
}

// enexpectedEOF is an error that occurs when the end of input is reached prematurely.
type enexpectedEOF struct {
	Expected  []types.TokenType // the type of token(s) that are expected
	LastToken types.Token       // the last token read
	Position  int               // the position of the error
}

// errorUnexpectedEOF creates a new unexpected EOF errors.
func errorUnexpectedEOF(lastToken types.Token, expected ...types.TokenType) *enexpectedEOF {
	return &enexpectedEOF{
		Expected:  expected,
		LastToken: lastToken,
		Position:  lastToken.Position,
	}
}

func (u *enexpectedEOF) Error() string {
	var expected []string
	for _, e := range u.Expected {
		expected = append(expected, fmt.Sprintf("%s", e))
	}
	return fmt.Sprintf("at position %d, reached end of input, but expected %s", u.Position, strings.Join(expected, ", "))
}

// readFailed is an error that occurs when tokens cannot be read.
type readFailed struct {
	Cause error
}

// errorReadFailed Creates a new read failed error.
func errorReadFailed(cause error) *readFailed {
	return &readFailed{cause}
}

func (f *readFailed) Error() string {
	return fmt.Sprintf("%s", f.Cause)
}

type tokenizerError struct {
	errorToken types.Token
}

func errorTokenizerError(errorToken types.Token) *tokenizerError {
	return &tokenizerError{errorToken}
}

func (t *tokenizerError) Error() string {
	return fmt.Sprintf("at position %d, %s", t.errorToken.Position, t.errorToken.Value)
}

// invalidUnits is an error that occurs when an invalid unit name is used.
type invalidUnits struct {
	Name     string // the name that is invalid
	Position int    // the position of the invalid unit name
	Width    int    // the width of the invalid unit name
}

// errorInvalidUnits Creates a new invalid units error.
func errorInvalidUnits(badToken types.Token) *invalidUnits {
	return &invalidUnits{
		Name:     badToken.Value,
		Position: badToken.Position,
		Width:    len(badToken.Value),
	}
}

func (u *invalidUnits) Error() string {
	return fmt.Sprintf("at position %d, '%s' is not a known measurement unit", u.Position, u.Name)
}

// invalidOperator occurs when an unsupported operator is used.
type invalidOperator struct {
	operator types.Token
}

func (i *invalidOperator) Error() string {
	return fmt.Sprintf("at position %d, found invalid operator %s", i.operator.Position, i.operator.Value)
}

// errorInvalidOperator Creates an invalid operator error.
func errorInvalidOperator(operator types.Token) *invalidOperator {
	return &invalidOperator{operator}
}
