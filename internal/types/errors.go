package types

import (
	"fmt"
	"strings"
)

// InputError represents an error that occurs when tokenizing, parsing or evaluating the user's input.
type InputError interface {
	// Error returns an error message.
	Error() string
	// Input returns the input string.
	Input() string
	// Position returns the position of the error.
	Position() (start, width int)
}

// ErrorUnexpectedToken creates a new unexpected token error.
func ErrorUnexpectedToken(input string, badToken Token, expected ...TokenType) *UnexpectedToken {
	return &UnexpectedToken{
		expected: expected,
		badToken: badToken,
		position: badToken.Position,
		width:    len(badToken.Value),
		input:    input,
	}
}

// ErrorUnexpectedEOF creates a new unexpected EOF errors.
func ErrorUnexpectedEOF(input string, lastToken Token, expected ...TokenType) *UnexpectedEOF {
	return &UnexpectedEOF{
		expected:  expected,
		lastToken: lastToken,
		position:  lastToken.Position,
		input:     input,
	}
}

// ErrorReadFailed Creates a new read failed error.
func ErrorReadFailed(input string, cause error) *ReadFailed {
	return &ReadFailed{cause, input}
}

// ErrorTokenizerError Creates a new tokenizer error.
func ErrorTokenizerError(input string, errorToken Token) *TokenizerError {
	return &TokenizerError{
		errorToken: errorToken,
		position:   errorToken.Position,
		width:      1,
		input:      input,
	}
}

// ErrorInvalidUnits Creates a new invalid units error.
func ErrorInvalidUnits(input string, badToken Token) *InvalidUnits {
	return &InvalidUnits{
		invalidName: badToken.Value,
		position:    badToken.Position,
		width:       len(badToken.Value),
		input:       input,
	}
}

// ErrorInvalidOperator Creates an invalid operator error.
func ErrorInvalidOperator(input string, invalid Token) *InvalidOperator {
	return &InvalidOperator{
		invalid:  invalid,
		input:    input,
		position: invalid.Position,
		width:    len(invalid.Value),
	}
}

// ErrorInvalidNumber Creates an invalid number error.
func ErrorInvalidNumber(input string, invalid Token) *InvalidNumber {
	return &InvalidNumber{
		invalid:  invalid,
		input:    input,
		position: invalid.Position,
		width:    len(invalid.Value),
	}
}

// ErrorInvalidUnitConversion Creates an invalid unit conversion error.
func ErrorInvalidUnitConversion(input string, from Token, to Token) *InvalidUnitConversion {
	return &InvalidUnitConversion{
		from:     from.Value,
		to:       to.Value,
		position: from.Position,
		width:    len(from.Value),
		input:    input,
	}
}

// InvalidUnitConversion is an error that occurs when a unit conversion is invalid; like meters to pounds.
type InvalidUnitConversion struct {
	from     string
	to       string
	position int
	width    int
	input    string
}

// Error returns an error message.
func (e *InvalidUnitConversion) Error() string {
	return fmt.Sprintf("cannot convert from %s to %s", e.from, e.to)
}

// Input returns the input string.
func (e *InvalidUnitConversion) Input() string {
	return e.input
}

// Position returns the position of the error.
func (e *InvalidUnitConversion) Position() (start, width int) {
	return e.position, e.width
}

// UnexpectedToken is an error indicated that an unexpected token found.
type UnexpectedToken struct {
	expected []TokenType // the token(s) that were expected
	badToken Token       // the unexpected token
	position int         // the position of the error
	width    int         // the width of the error
	input    string      // the input string
}

// Input returns the input string.
func (u *UnexpectedToken) Input() string {
	return u.input
}

// Position returns the position of the error.
func (u *UnexpectedToken) Position() (start, width int) {
	return u.position, u.width
}

func (u *UnexpectedToken) Error() string {
	var expected []string
	for _, e := range u.expected {
		expected = append(expected, fmt.Sprintf("%s", e))
	}
	return fmt.Sprintf("got '%s', but expected %s", u.badToken.Value, strings.Join(expected, ", "))
}

// UnexpectedEOF is an error indicating that a premature EOF was encountered.
type UnexpectedEOF struct {
	expected  []TokenType // the token(s) that were expected
	lastToken Token       // the last token read
	position  int         // the position of the error
	input     string      // the input string
}

func (u *UnexpectedEOF) Error() string {
	var expected []string
	for _, e := range u.expected {
		expected = append(expected, fmt.Sprintf("%s", e))
	}
	return fmt.Sprintf("reached end of input, but expected %s", strings.Join(expected, ", "))
}

// Input returns the input string.
func (u *UnexpectedEOF) Input() string {
	return u.input
}

// Position returns the position of the error.
func (u *UnexpectedEOF) Position() (start, width int) {
	return u.position, 1
}

// ReadFailed is an error that occurs when tokens cannot be read.
type ReadFailed struct {
	cause error  // the cause of the read failure
	input string // the input string
}

func (r *ReadFailed) Error() string {
	return fmt.Sprintf("%s", r.cause)
}

// Input returns the input string.
func (r *ReadFailed) Input() string {
	return r.input
}

// Position returns the position of the error.
func (r *ReadFailed) Position() (start, width int) {
	return 1, 1
}

// TokenizerError is an error encountered during tokenization.
type TokenizerError struct {
	errorToken Token  // the error token from the tokenizer
	position   int    // the position of the error
	width      int    // the width of the error
	input      string // the input string
}

func (t *TokenizerError) Error() string {
	return fmt.Sprintf("%s", t.errorToken.Value)
}

// Input returns the input string.
func (t *TokenizerError) Input() string {
	return t.input
}

// Position returns the position of the error.
func (t *TokenizerError) Position() (start, width int) {
	return t.position, t.width
}

// InvalidUnits is an error indicating an invalid unit of measure was encountered.
type InvalidUnits struct {
	invalidName string // the name that is not a valid unit of measurement
	position    int    // the position of the invalid unit invalidName
	width       int    // the width of the invalid unit invalidName
	input       string // the input string
}

func (u *InvalidUnits) Error() string {
	return fmt.Sprintf("'%s' is not a known measurement unit", u.invalidName)
}

// Input returns the input string.
func (u *InvalidUnits) Input() string {
	return u.input
}

// Position returns the position of the error.
func (u *InvalidUnits) Position() (start, width int) {
	return u.position, u.width
}

// InvalidOperator is an error indicating an invalid operator was encountered.
type InvalidOperator struct {
	invalid  Token  // the operator that is not valid
	input    string // the input string
	position int    // the position of the error
	width    int    // the width of the error
}

func (i *InvalidOperator) Error() string {
	return fmt.Sprintf("found invalid operator %s", i.invalid)
}

// Input returns the input string.
func (i *InvalidOperator) Input() string {
	return i.input
}

// Position returns the position of the error.
func (i *InvalidOperator) Position() (start, width int) {
	return i.position, i.width
}

// InvalidNumber is an error indicating an invalid number was encountered.
type InvalidNumber struct {
	invalid  Token  // the number that is not valid
	input    string // the input string
	position int    // the position of the error
	width    int    // the width of the error
}

func (i InvalidNumber) Error() string {
	return i.Error()
}

// Input returns the input string.
func (i InvalidNumber) Input() string {
	return i.input
}

// Position returns the position of the error.
func (i InvalidNumber) Position() (start, width int) {
	return i.position, i.width
}
