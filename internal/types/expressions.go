package types

import (
	"fmt"
	u "github.com/bcicen/go-units"
)

// Amount is the result of evaluating an expression.
type Amount struct {
	Value float64
	Units Token
}

// Expression is something that can be evaluated.
type Expression interface {
	Eval(input string) (Amount, InputError)
}

// Value represents a fixed Value like "2 pounds".
type Value struct {
	number float64
	unit   Token
}

// NewValue creates a new Value.
func NewValue(number float64, unit Token) Value {
	return Value{number, unit}
}

// Eval evaluates a simple Value expression.
func (v Value) Eval(input string) (Amount, InputError) {
	// validate the units
	var amount Amount
	_, err := u.Find(v.unit.Value)
	if err != nil {
		return amount, ErrorInvalidUnits(input, v.unit)
	}
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

// AdditionExpr creates a new expression that performs Addition.
func AdditionExpr(left, right Expression) Addition {
	return Addition{left, right}
}

// Eval evaluates an Addition expression.
func (s Addition) Eval(input string) (sum Amount, err InputError) {
	add := func(l float64, r float64) float64 { return l + r }
	return eval(s.left, s.right, add, input)
}

func (s Addition) String() string {
	return fmt.Sprintf("%s + %s", s.left, s.right)
}

// Subtraction an expression that performs subtraction.
type Subtraction struct {
	left  Expression
	right Expression
}

// SubtractionExpr creates a new expression that performs Subtraction.
func SubtractionExpr(left, right Expression) Subtraction {
	return Subtraction{left, right}
}

// Eval evaluates a Subtraction expression.
func (s Subtraction) Eval(input string) (diff Amount, err InputError) {
	subtract := func(l float64, r float64) float64 { return l - r }
	return eval(s.left, s.right, subtract, input)
}

func (s Subtraction) String() string {
	return fmt.Sprintf("%s - %s", s.left, s.right)
}

// UnitConversion converts between units of measure
type UnitConversion struct {
	expr        Expression
	targetUnits Token
}

// UnitConversionExpr creates a new unit conversion expression.
func UnitConversionExpr(expr Expression, targetUnits Token) UnitConversion {
	return UnitConversion{expr, targetUnits}
}

// Eval evaluates a unit conversion expression.s
func (c UnitConversion) Eval(input string) (amount Amount, err InputError) {
	// evaluate the expression
	amount, err = c.expr.Eval(input)
	if err != nil {
		return amount, err
	}
	// is unit conversion needed?
	if amount.Units.String() != c.targetUnits.Value {
		fromUnits, unitErr := u.Find(amount.Units.Value)
		if unitErr != nil {
			return amount, ErrorInvalidUnits(input, amount.Units)
		}

		toUnits, unitErr := u.Find(c.targetUnits.Value)
		if unitErr != nil {
			return amount, ErrorInvalidUnits(input, c.targetUnits)
		}

		if fromUnits.Name == toUnits.Name {
			// no conversion is necessary; for example 2 kilograms in kg requires no conversion
			amount = Amount{amount.Value, c.targetUnits}
			return amount, nil
		}

		// unit conversion
		value, err := u.ConvertFloat(amount.Value, fromUnits, toUnits)
		if err != nil {
			return amount, ErrorInvalidUnitConversion(input, amount.Units, c.targetUnits)
		}
		amount = Amount{value.Float(), c.targetUnits}
	}
	return amount, nil
}

func (c UnitConversion) String() string {
	return fmt.Sprintf("%s in %s", c.expr, c.targetUnits)
}

type opFunction func(float64, float64) float64

func eval(leftExpr Expression, rightExpr Expression, opFunc opFunction, input string) (Amount, InputError) {
	var result Amount

	// evaluate the left side
	left, err := leftExpr.Eval(input)
	if err != nil {
		return result, err
	}

	// prefer the units of the left side
	targetUnit := left.Units

	// unit conversion, if needed
	if left.Units != targetUnit {
		left, err = UnitConversionExpr(leftExpr, targetUnit).Eval(input)
		if err != nil {
			return result, err
		}
	}

	// evaluate the right side
	right, err := rightExpr.Eval(input)
	if err != nil {
		return result, err
	}

	// unit conversion, if needed
	if right.Units != targetUnit {
		right, err = UnitConversionExpr(rightExpr, targetUnit).Eval(input)
		if err != nil {
			return result, err
		}
	}

	result = Amount{
		Value: opFunc(left.Value, right.Value),
		Units: targetUnit,
	}
	return result, nil
}
