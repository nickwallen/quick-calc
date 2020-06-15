package types

import (
	"fmt"
)

// InvalidUnitConversion is an error that occurs when a unit conversion is invalid; like meters to pounds.
type InvalidUnitConversion struct {
	From string
	To   string
}

func (e *InvalidUnitConversion) Error() string {
	return fmt.Sprintf("cannot convert from %s to %s", e.From, e.To)
}

// ErrorInvalidUnitConversion Creates an invalid unit conversion error.
func ErrorInvalidUnitConversion(from, to string) *InvalidUnitConversion {
	return &InvalidUnitConversion{from, to}
}
