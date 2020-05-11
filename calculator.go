package calc

import (
	"fmt"
	u "github.com/bcicen/go-units"
	"github.com/nickwallen/quick-calc/internal/io"
	"github.com/nickwallen/quick-calc/internal/parser"
	"github.com/nickwallen/quick-calc/internal/tokenizer"
)

// Amount is the result of evaluating an expression.
type Amount struct {
	Value float64
	Units string
}

// Calculate evaluates an input expression and returns the value as a string.
func Calculate(input string) (string, error) {
	var result string
	amount, err := CalculateAmount(input)
	if err != nil {
		return result, fmt.Errorf("execution error: %s", err.Error())
	}
	result = fmt.Sprintf("%.2f %s", amount.Value, amount.Units)
	return result, nil
}

// CalculateAmount evaluates an input expression and returns an Amount object.
func CalculateAmount(input string) (Amount, error) {
	var amount Amount

	// the tokenizer runs in the background populating the token channel
	tokens := io.NewTokenChannel()
	go tokenizer.Tokenize(input, tokens)

	// parse the tokens
	expr, err := parser.Parse(tokens)
	if err != nil {
		return amount, fmt.Errorf("parse error: %s", err.Error())
	}

	// eval the expression
	return eval(&expr)
}

type opFunction func(float64, float64) float64

func eval(expr *parser.Expression) (Amount, error) {
	switch expr.Op {
	case parser.PlusOp:
		opFunc := func(l float64, r float64) float64 { return l + r }
		return evalBinaryExpr(opFunc, expr)

	case parser.MinusOp:
		opFunc := func(l float64, r float64) float64 { return l - r }
		return evalBinaryExpr(opFunc, expr)

	case parser.ValueOp:
		return evalValue(expr)

	case parser.ConvertOp:
		return evalConversion(expr)

	default:
		var amount Amount
		return amount, fmt.Errorf("unsupported op: '%s'", expr.Op)
	}
}

func evalBinaryExpr(opFunc opFunction, expr *parser.Expression) (amount Amount, err error) {
	left, err := eval(expr.Left)
	if err != nil {
		return amount, err
	}

	right, err := eval(expr.Right)
	if err != nil {
		return amount, err
	}

	// unit conversion, if necessary
	if left.Units != expr.TargetUnits.String() {
		left, err = convert(left, expr.TargetUnits)
		if err != nil {
			return amount, err
		}
	}

	// unit conversion, if necessary
	if right.Units != expr.TargetUnits.String() {
		right, err = convert(right, expr.TargetUnits)
		if err != nil {
			return amount, err
		}
	}

	result := Amount{opFunc(left.Value, right.Value), left.Units}
	return result, nil
}

func evalValue(expr *parser.Expression) (Amount, error) {
	amount := Amount{expr.Value, expr.ValueUnits.String()}

	// unit conversion, if necessary
	if expr.ValueUnits != expr.TargetUnits {
		amount, err := convert(amount, expr.TargetUnits)
		if err != nil {
			return amount, err
		}
	}
	return amount, nil
}

func evalConversion(expr *parser.Expression) (amount Amount, err error) {
	//amount := Amount{expr.Value, expr.ValueUnits.String()}
	amount, err = eval(expr.Left)
	if err != nil {
		return amount, err
	}

	// unit conversion, if necessary
	if amount.Units != expr.TargetUnits.String() {
		amount, err = convert(amount, expr.TargetUnits)
		if err != nil {
			return amount, err
		}
	}
	return amount, nil
}

func convert(from Amount, units parser.Units) (Amount, error) {
	var amount Amount
	fromUnits, err := u.Find(from.Units)
	if err != nil {
		return amount, err
	}

	toUnits, err := u.Find(units.String())
	if err != nil {
		return amount, err
	}

	// no conversion may be necessary; for example 2 kilograms in kg
	if fromUnits.Name == toUnits.Name {
		amount = Amount{from.Value, units.String()}
		return amount, nil
	}

	value, err := u.ConvertFloat(from.Value, fromUnits, toUnits)
	if err != nil {
		return amount, err
	}

	amount = Amount{value.Float(), units.String()}
	return amount, nil
}
