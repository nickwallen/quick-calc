package toks

import "fmt"

// Expression An expression can be evaluated resulting in an amount.
type Expression interface {
	Evaluate() (Amount, error)
}

// Amount Expressions should be evaluated to an amount.
type Amount struct {
	Value int64
	Units AmountUnits
}

// AmountOf Create a new Amount.
func AmountOf(value int64, units AmountUnits) Amount {
	return Amount{Value: value, Units: units}
}

// Evaluate An Amount is an expression and can be evaluated.
func (amount Amount) Evaluate() (Amount, error) {
	return amount, nil
}

// AmountUnits The units of an Amount, like kilograms.
type AmountUnits struct {
	units string
}

// UnitsOf Create a new AmountUnits.
func UnitsOf(units string) AmountUnits {
	return AmountUnits{units: units}
}

func (units AmountUnits) String() string {
	return units.units
}

// Sum An expression that can sum two expressions.
type Sum struct {
	left  Expression
	right Expression
	units AmountUnits
}

// SumOf Create a new Sum expression.
func SumOf(left Expression, right Expression, units AmountUnits) Sum {
	return Sum{left: left, right: right, units: units}
}

// Evaluate Sum is an expression that can be evaluated.
func (sum Sum) Evaluate() (Amount, error) {
	var amount Amount

	// evaluate the left expression
	leftResult, err := sum.left.Evaluate()
	if err != nil {
		return amount, err
	}

	// evaluate the right expression
	rightResult, err := sum.right.Evaluate()
	if err != nil {
		return amount, err
	}

	if leftResult.Units == rightResult.Units {
		// same units, so just sum the two
		amount = Amount{
			Value: leftResult.Value + rightResult.Value,
			Units: sum.units,
		}
		return amount, nil
	}

	// TODO support unit conversion
	return amount, fmt.Errorf("unit conversion from '%s' to '%s' not yet supported", leftResult.Units, rightResult.Units)
}
