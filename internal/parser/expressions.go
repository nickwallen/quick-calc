package parser

import (
	"github.com/bcicen/go-units"
	"github.com/nickwallen/qcalc/internal/tokens"
)

// Expression An Expression can be evaluated resulting in an Amount.
type Expression interface {
	Evaluate() (Amount, error)
}

// Amount Expressions should be evaluated to an Amount.
type Amount struct {
	Value float64
	Units Units
}

// AmountOf Create a new Amount.
func AmountOf(value float64, units Units) Amount {
	return Amount{Value: value, Units: units}
}

// Evaluate An Amount is an Expression and can be evaluated.
func (amount Amount) Evaluate() (Amount, error) {
	return amount, nil
}

// Units The units of an Amount, like kilograms.
type Units struct {
	units string
}

// UnitsOf Create a new Units.
func UnitsOf(input string) (Units, error) {
	var result Units

	// ensure that the Units are valid
	_, err := units.Find(input)
	if err != nil {
		return result, err
	}

	// use a 'standard' name which may differ from what the user input
	result.units = input
	return result, nil
}

func (units Units) String() string {
	return units.units
}

// operator An Expression that operates on multiple operands like +, -, *, and /.
type operator struct {
	left     Expression
	right    Expression
	units    Units
	operator tokens.TokenType
}

// operatorOf Create a new operator Expression.
func operatorOf(left Expression, right Expression, units Units, op tokens.TokenType) operator {
	return operator{left: left, right: right, units: units, operator: op}
}

// allows multiple operations like +, -, etc to be handled by operator
type opFunction func(float64, float64) float64

// Evaluate operator is an Expression that can be evaluated.
func (op operator) Evaluate() (Amount, error) {
	var amount Amount
	var operationFn opFunction

	// what operation will be performed?
	switch op.operator {
	case tokens.Plus:
		operationFn = func(x float64, y float64) float64 { return x + y }
	case tokens.Minus:
		operationFn = func(x float64, y float64) float64 { return x - y }
	case tokens.Multiply:
		operationFn = func(x float64, y float64) float64 { return x * y }
	case tokens.Divide:
		operationFn = func(x float64, y float64) float64 { return x / y }
	}

	// Evaluate the first Expression
	result1, err := op.left.Evaluate()
	if err != nil {
		return amount, err
	}

	// Evaluate the second Expression
	result2, err := op.right.Evaluate()
	if err != nil {
		return amount, err
	}

	// convert the Units of the first result
	if result1.Units != op.units {
		result1, err = unitConverterOf(result1, op.units).Evaluate()
		if err != nil {
			return amount, err
		}
	}

	// convert the Units of the second result
	if result2.Units != op.units {
		result2, err = unitConverterOf(result2, op.units).Evaluate()
		if err != nil {
			return amount, err
		}
	}

	// execute the operation
	value := operationFn(result1.Value, result2.Value)
	return AmountOf(value, op.units), nil
}

// unitConverter Converts from one unit to another.
type unitConverter struct {
	from    Amount
	toUnits Units
}

// unitConverterOf Converts from one unit to another.
func unitConverterOf(from Amount, toUnits Units) unitConverter {
	return unitConverter{from: from, toUnits: toUnits}
}

// Evaluate Converts from one unit to another.
func (c unitConverter) Evaluate() (Amount, error) {
	var amount Amount
	fromUnits, err := units.Find(c.from.Units.units)
	if err != nil {
		return amount, err
	}
	toUnits, err := units.Find(c.toUnits.units)
	if err != nil {
		return amount, err
	}

	// no conversion may be necessary; for example 2 kilograms in kg
	if fromUnits.Name == toUnits.Name {
		return AmountOf(c.from.Value, c.from.Units), nil
	}

	value, err := units.ConvertFloat(c.from.Value, fromUnits, toUnits)
	if err != nil {
		return amount, err
	}
	return AmountOf(value.Float(), c.toUnits), nil
}
