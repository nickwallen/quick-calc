package parser

import (
	"github.com/bcicen/go-units"
	"github.com/nickwallen/quick-calc/internal/tokens"
)

// Expression is something that can be evaluated. The parser is responsible
// for constructing an expression.
type Expression struct {
	Op          Operator
	Left        *Expression
	Right       *Expression
	TargetUnits Units
	Value       float64
	ValueUnits  Units
}

// Operator The operator of an expression, like + for addition.
type Operator string

const (
	// PlusOp signifies binary addition.
	PlusOp Operator = "+"

	// MinusOp signifies binary subtraction.
	MinusOp Operator = "-"

	// ConvertOp signifies conversion between different units of measure.
	ConvertOp Operator = "->"

	// ValueOp signifies a known, fixed valueExpr.
	ValueOp Operator = "="
)

// Units are a measure of a physical property.
type Units struct {
	units string
}

// UnitsOf Create a new TargetUnits valueExpr.
func UnitsOf(input string) (Units, error) {
	var result Units

	// ensure that the TargetUnits are valid
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

// binaryExpr Create an expression where two values are acted on by an operator.
func binaryExpr(left Expression, right Expression, units Units, opType tokens.TokenType) Expression {
	var op Operator
	switch opType {
	case tokens.Plus:
		op = PlusOp
	case tokens.Minus:
		op = MinusOp
	}
	return Expression{
		Op:          op,
		Left:        &left,
		Right:       &right,
		TargetUnits: units,
	}
}

// valueExpr Creates an expression consisting of a known, fixed value.
func valueExpr(value float64, units Units) Expression {
	return Expression{
		Op:          ValueOp,
		Value:       value,
		ValueUnits:  units,
		TargetUnits: units,
	}
}

// conversion Converts from one unit of measure to another.
func conversion(from Expression, toUnits Units) Expression {
	return Expression{
		Op:          ConvertOp,
		Left:        &from,
		TargetUnits: toUnits,
	}
}
