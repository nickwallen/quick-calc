package types

import (
	"fmt"
	u "github.com/bcicen/go-units"
)

// Amount is the result of evaluating an expression.
type Amount struct {
	Value float64
	Units string
}

// Expression is something that can be evaluated.
type Expression interface {
	Eval() (Amount, error)
}

// Value represents a fixed Value like "2 pounds".
type Value struct {
	number float64
	unit   string
}

// NewValue creates a new Value.
func NewValue(number float64, unit string) Value {
	return Value{number, unit}
}

// Eval evaluates a simple Value expression.
func (v Value) Eval() (Amount, error) {
	return Amount{
		Value: v.number,
		Units: v.unit,
	}, nil
}

func (v Value) String() string {
	return fmt.Sprintf("%.2f %s", v.number, v.unit)
}

// Addition is an expression that performs addition.
type Addition struct {
	left  Expression
	right Expression
}

// DoAddition creates a new expression that performs Addition.
func DoAddition(left, right Expression) Addition {
	return Addition{left, right}
}

// Eval evaluates an Addition expression.
func (s Addition) Eval() (sum Amount, err error) {
	add := func(l float64, r float64) float64 { return l + r }
	return eval(s.left, s.right, add)
}

func (s Addition) String() string {
	return fmt.Sprintf("%s + %s", s.left, s.right)
}

// Subtraction an expression that performs subtraction.
type Subtraction struct {
	left  Expression
	right Expression
}

// DoSubtraction creates a new expression that performs Subtraction.
func DoSubtraction(left, right Expression) Subtraction {
	return Subtraction{left, right}
}

// Eval evaluates a Subtraction expression.
func (s Subtraction) Eval() (diff Amount, err error) {
	subtract := func(l float64, r float64) float64 { return l - r }
	return eval(s.left, s.right, subtract)
}

func (s Subtraction) String() string {
	return fmt.Sprintf("%s - %s", s.left, s.right)
}

// UnitConversion converts between units of measure
type UnitConversion struct {
	expr        Expression
	targetUnits string
}

// DoUnitConversion creates a new unit conversion expression.
func DoUnitConversion(expr Expression, targetUnits string) UnitConversion {
	return UnitConversion{expr, targetUnits}
}

// Eval evaluates a unit conversion expression.s
func (c UnitConversion) Eval() (amount Amount, err error) {
	// evaluate the expression
	amount, err = c.expr.Eval()
	if err != nil {
		return amount, err
	}
	// is unit conversion needed?
	if amount.Units != c.targetUnits {
		fromUnits, err := u.Find(amount.Units)
		if err != nil {
			return amount, err
		}

		toUnits, err := u.Find(c.targetUnits)
		if err != nil {
			return amount, err
		}

		if fromUnits.Name == toUnits.Name {
			// no conversion is necessary; for example 2 kilograms in kg requires no conversion
			amount = Amount{amount.Value, c.targetUnits}
			return amount, nil
		}

		// unit conversion
		value, err := u.ConvertFloat(amount.Value, fromUnits, toUnits)
		if err != nil {
			return amount, ErrorInvalidUnitConversion(amount.Units, c.targetUnits)
		}
		amount = Amount{value.Float(), c.targetUnits}
	}
	return amount, nil
}

func (c UnitConversion) String() string {
	return fmt.Sprintf("%s in %s", c.expr, c.targetUnits)
}

type opFunction func(float64, float64) float64

func eval(leftExpr Expression, rightExpr Expression, opFunc opFunction) (Amount, error) {
	var result Amount

	// evaluate the left side
	left, err := leftExpr.Eval()
	if err != nil {
		return result, err
	}

	// prefer the units of the left side
	targetUnit := left.Units

	// unit conversion, if needed
	if left.Units != targetUnit {
		left, err = DoUnitConversion(leftExpr, targetUnit).Eval()
		if err != nil {
			return result, err
		}
	}

	// evaluate the right side
	right, err := rightExpr.Eval()
	if err != nil {
		return result, err
	}

	// unit conversion, if needed
	if right.Units != targetUnit {
		right, err = DoUnitConversion(rightExpr, targetUnit).Eval()
		if err != nil {
			return result, err
		}
	}

	// do the summation
	result = Amount{
		Value: opFunc(left.Value, right.Value),
		Units: targetUnit,
	}
	return result, nil
}
