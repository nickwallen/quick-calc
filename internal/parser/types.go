package parser

import (
	"fmt"
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

func (e Expression) String() string {
	switch e.Op {
	case ValueOp:
		return fmt.Sprintf("[%.2f %s]", e.Value, e.ValueUnits)
	case ConvertOp:
		return fmt.Sprintf("[%s] -> [%s]", e.Left, e.TargetUnits)
	default:
		return fmt.Sprintf("%s (%v, %v) in %s", e.Op, e.Left, e.Right, e.TargetUnits)
	}
}

// Operator The operator of an expression, like + for addition.
type Operator string

// Operator types
const (
	Noop      Operator = ""
	PlusOp    Operator = "+"
	MinusOp   Operator = "-"
	ConvertOp Operator = "->"
	ValueOp   Operator = "="
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

	// use a 'standard' badName which may differ from what the user input
	result.units = input
	return result, nil
}

func (units Units) String() string {
	return units.units
}

// binaryExpr Create an expression where two values are acted on by an operator.
func binaryExpr(opType tokens.TokenType, left Expression, right Expression, units Units) Expression {
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
