package types

import (
	"fmt"
)

// InvalidUnitConversion is an error that occurs when a unit conversion is invalid; like meters to pounds.
type InvalidUnitConversion struct {
	From     string
	To       string
	Position int
	Width    int
}

func (e *InvalidUnitConversion) Error() string {
	return fmt.Sprintf("at position %d, cannot convert from %s to %s", e.Position, e.From, e.To)
}

// ErrorInvalidUnitConversion Creates an invalid unit conversion error.
func ErrorInvalidUnitConversion(from Token, to Token) *InvalidUnitConversion {
	return &InvalidUnitConversion{
		From:     from.Value,
		To:       to.Value,
		Position: from.Position,
		Width:    len(from.Value),
	}
}
