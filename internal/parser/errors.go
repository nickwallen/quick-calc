package parser

import (
	"fmt"
	"github.com/nickwallen/quick-calc/internal/types"
	"strings"
)

// UnexpectedToken is an error that occurs when an unexpected token is found.
type UnexpectedToken struct {
	Expected []types.TokenType // the type of token(s) that are expected
	BadToken types.Token       // the token that we got
	Position int               // the position of the error
	Width    int               // the width of the error
}

// ErrorUnexpectedToken creates a new unexpected token error.
func ErrorUnexpectedToken(badToken types.Token, expected ...types.TokenType) *UnexpectedToken {
	return &UnexpectedToken{
		Expected: expected,
		BadToken: badToken,
		Position: badToken.Position,
		Width:    len(badToken.Value),
	}
}

func (u *UnexpectedToken) Error() string {
	var expected []string
	for _, e := range u.Expected {
		expected = append(expected, fmt.Sprintf("%s", e))
	}
	return fmt.Sprintf("at position %d, got '%s', but expected %s", u.Position, u.BadToken.Value, strings.Join(expected, ", "))
}

// UnexpectedEOF is an error that occurs when the end of input is reached prematurely.
type UnexpectedEOF struct {
	Expected  []types.TokenType // the type of token(s) that are expected
	LastToken types.Token       // the last token read
	Position  int               // the position of the error
}

// ErrorUnexpectedEOF creates a new unexpected EOF errors.
func ErrorUnexpectedEOF(lastToken types.Token, expected ...types.TokenType) *UnexpectedEOF {
	return &UnexpectedEOF{
		Expected:  expected,
		LastToken: lastToken,
		Position:  lastToken.Position,
	}
}

func (u *UnexpectedEOF) Error() string {
	var expected []string
	for _, e := range u.Expected {
		expected = append(expected, fmt.Sprintf("%s", e))
	}
	return fmt.Sprintf("at position %d, reached end of input, but expected %s", u.Position, strings.Join(expected, ", "))
}

// ReadFailed is an error that occurs when tokens cannot be read.
type ReadFailed struct {
	Cause error
}

// ErrorReadFailed Creates a new read failed error.
func ErrorReadFailed(cause error) *ReadFailed {
	return &ReadFailed{cause}
}

// ErrorReadFailedNoCause Creates a new read failed error where no cause is known.
func ErrorReadFailedNoCause() *ReadFailed {
	return &ReadFailed{fmt.Errorf("unknown Cause")}
}

func (f *ReadFailed) Error() string {
	return fmt.Sprintf("unable to read tokens, %s", f.Cause)
}

// InvalidUnits is an error that occurs when an invalid unit name is used.
type InvalidUnits struct {
	Name     string // the name that is invalid
	Position int    // the position of the invalid unit name
	Width    int    // the width of the invalid unit name
}

// ErrorInvalidUnits Creates a new invalid units error.
func ErrorInvalidUnits(badToken types.Token) *InvalidUnits {
	return &InvalidUnits{
		Name:     badToken.Value,
		Position: badToken.Position,
		Width:    len(badToken.Value),
	}
}

func (u *InvalidUnits) Error() string {
	return fmt.Sprintf("at position %d, '%s' is not a known measurement unit", u.Position, u.Name)
}

// InvalidOperator occurs when an unsupported operator is used.
type InvalidOperator struct {
	operator types.Token
}

func (i *InvalidOperator) Error() string {
	return fmt.Sprintf("at position %d, found invalid operator %s", i.operator.Position, i.operator.Value)
}

// ErrorInvalidOperator Creates an invalid operator error.
func ErrorInvalidOperator(operator types.Token) *InvalidOperator {
	return &InvalidOperator{operator}
}
