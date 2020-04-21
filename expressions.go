package toks

import (
	units "github.com/bcicen/go-units"
)

// Expression An Expression can be evaluated resulting in an amount.
type Expression interface {
	Evaluate() (amount, error)
}

// amount Expressions should be evaluated to an amount.
type amount struct {
	value float64
	units amountUnits
}

// amountOf Create a new amount.
func amountOf(value float64, units amountUnits) amount {
	return amount{value: value, units: units}
}

// Evaluate An amount is an Expression and can be evaluated.
func (amount amount) Evaluate() (amount, error) {
	return amount, nil
}

// amountUnits The units of an amount, like kilograms.
type amountUnits struct {
	units string
}

// unitsOf Create a new amountUnits.
func unitsOf(input string) (amountUnits, error) {
	var result amountUnits

	// ensure that the units are valid
	unit, err := units.Find(input)
	if err != nil {
		return result, err
	}

	// use a 'standard' name which may differ from what the user input
	result.units = unit.PluralName()
	return result, nil
}

func (units amountUnits) String() string {
	return units.units
}

// operator An Expression that operates on multiple operands like +, -, *, and /.
type operator struct {
	left     Expression
	right    Expression
	units    amountUnits
	operator TokenType
}

// operatorOf Create a new operator Expression.
func operatorOf(left Expression, right Expression, units amountUnits, op TokenType) operator {
	return operator{left: left, right: right, units: units, operator: op}
}

// allows multiple operations like +, -, etc to be handled by operator
type opFunction func(float64, float64) float64

// Evaluate operator is an Expression that can be evaluated.
func (op operator) Evaluate() (amount, error) {
	var amount amount
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

	// convert the units of the first result
	if result1.units != op.units {
		result1, err = unitConverterOf(result1, op.units).Evaluate()
		if err != nil {
			return amount, err
		}
	}

	// convert the units of the second result
	if result2.units != op.units {
		result2, err = unitConverterOf(result2, op.units).Evaluate()
		if err != nil {
			return amount, err
		}
	}

	// execute the operation
	value := operationFn(result1.value, result2.value)
	return amountOf(value, op.units), nil
}

// unitConverter Converts from one unit to another.
type unitConverter struct {
	from    amount
	toUnits amountUnits
}

// unitConverterOf Converts from one unit to another.
func unitConverterOf(from amount, toUnits amountUnits) unitConverter {
	return unitConverter{from: from, toUnits: toUnits}
}

// Evaluate Converts from one unit to another.
func (c unitConverter) Evaluate() (amount, error) {
	var amount amount
	fromUnits, err := units.Find(c.from.units.units)
	if err != nil {
		return amount, err
	}
	toUnits, err := units.Find(c.toUnits.units)
	if err != nil {
		return amount, err
	}
	value, err := units.ConvertFloat(c.from.value, fromUnits, toUnits)
	if err != nil {
		return amount, err
	}
	return amountOf(value.Float(), c.toUnits), nil
}
