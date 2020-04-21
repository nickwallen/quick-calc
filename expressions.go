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

// Operator An expression that operates on multiple operands like +, -, *, and /.
type Operator struct {
	left     Expression
	right    Expression
	units    AmountUnits
	operator TokenType
}

// OperatorOf Create a new Operator expression.
func OperatorOf(left Expression, right Expression, units AmountUnits, operator TokenType) Operator {
	return Operator{left: left, right: right, units: units, operator: operator}
}

// allows multiple operations like +, -, etc to be handled by Operator
type opFunction func(float64, float64) float64

// Evaluate Operator is an expression that can be evaluated.
func (op Operator) Evaluate() (Amount, error) {
	var amount Amount
	var operationFn opFunction

	// what operation will be performed?
	switch op.operator {
	case Plus:
		operationFn = func(x float64, y float64) float64 { return x + y }
	case Minus:
		operationFn = func(x float64, y float64) float64 { return x - y }
	case Multiply:
		operationFn = func(x float64, y float64) float64 { return x * y }
	case Divide:
		operationFn = func(x float64, y float64) float64 { return x / y }
	}

	// evaluate the first expression
	result1, err := op.left.Evaluate()
	if err != nil {
		return amount, err
	}

	// evaluate the second expression
	result2, err := op.right.Evaluate()
	if err != nil {
		return amount, err
	}

	// convert the units of the first result
	if result1.Units != op.units {
		result1, err = UnitConverterOf(result1, op.units).Evaluate()
		if err != nil {
			return amount, err
		}
	}

	// convert the units of the second result
	if result2.Units != op.units {
		result2, err = UnitConverterOf(result2, op.units).Evaluate()
		if err != nil {
			return amount, err
		}
	}

	// execute the operation
	value := operationFn(result1.Value, result2.Value)
	return AmountOf(value, op.units), nil
}

// UnitConverter Converts from one unit to another.
type UnitConverter struct {
	from    Amount
	toUnits AmountUnits
}

// UnitConverterOf Converts from one unit to another.
func UnitConverterOf(from Amount, toUnits AmountUnits) UnitConverter {
	return UnitConverter{from: from, toUnits: toUnits}
}

// Evaluate Converts from one unit to another.
func (c UnitConverter) Evaluate() (Amount, error) {
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
