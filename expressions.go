package toks

import (
	units "github.com/bcicen/go-units"
)

// Expression An expression can be evaluated resulting in an amount.
type Expression interface {
	Evaluate() (Amount, error)
}

// Amount Expressions should be evaluated to an amount.
type Amount struct {
	Value float64
	Units AmountUnits
}

// AmountOf Create a new Amount.
func AmountOf(value float64, units AmountUnits) Amount {
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
	add := func(x float64, y float64) float64 { return x + y }
	return evaluateOp(sum.left, sum.right, sum.units, add)
}

// Diff An expression that can find the difference of two expressions.
type Diff struct {
	left  Expression
	right Expression
	units AmountUnits
}

// DiffOf Create a new Sum expression.
func DiffOf(left Expression, right Expression, units AmountUnits) Diff {
	return Diff{left: left, right: right, units: units}
}

// Evaluate Sum is an expression that can be evaluated.
func (diff Diff) Evaluate() (Amount, error) {
	add := func(x float64, y float64) float64 { return x - y }
	return evaluateOp(diff.left, diff.right, diff.units, add)
}

type operandFn func(float64, float64) float64

func evaluateOp(left Expression, right Expression, units AmountUnits, op operandFn) (Amount, error) {
	var amount Amount

	// evaluate the left expression
	leftResult, err := left.Evaluate()
	if err != nil {
		return amount, err
	}

	// evaluate the right expression
	rightResult, err := right.Evaluate()
	if err != nil {
		return amount, err
	}

	if leftResult.Units != units {
		// convert the left side
		leftResult, err = convOf(leftResult, units).Evaluate()
		if err != nil {
			return amount, err
		}
	}

	if rightResult.Units != units {
		// convert the right side
		rightResult, err = convOf(rightResult, units).Evaluate()
		if err != nil {
			return amount, err
		}
	}

	return AmountOf(op(leftResult.Value, rightResult.Value), units), nil
}

type convertUnits struct {
	from    Amount
	toUnits AmountUnits
}

// ConfOf Create an expression that can convert units from one to another.
func convOf(from Amount, toUnits AmountUnits) convertUnits {
	return convertUnits{from: from, toUnits: toUnits}
}

// Evaluate Converts units from one to another
func (c convertUnits) Evaluate() (Amount, error) {
	var amount Amount
	fromUnits, err := units.Find(c.from.Units.units)
	if err != nil {
		return amount, err
	}
	toUnits, err := units.Find(c.toUnits.units)
	if err != nil {
		return amount, err
	}
	value, err := units.ConvertFloat(c.from.Value, fromUnits, toUnits)
	if err != nil {
		return amount, err
	}
	return AmountOf(value.Float(), c.toUnits), nil
}
