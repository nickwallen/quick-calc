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
	amt, err := CalculateAmount(input)
	if err != nil {
		return result, fmt.Errorf("execution error: %s", err.Error())
	}
	result = fmt.Sprintf("%.2f %s", amt.Value, amt.Units)
	return result, nil
}

// CalculateAmount evaluates an input expression and returns an Amount object.
func CalculateAmount(input string) (amt Amount, err error) {
	tokens := io.NewTokenChannel()
	go tokenizer.Tokenize(input, tokens)
	expr, err := parser.Parse(tokens)
	if err != nil {
		return amt, fmt.Errorf("parse error: %s", err.Error())
	}
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
		var amt Amount
		return amt, fmt.Errorf("unsupported op: '%s'", expr.Op)
	}
}

func evalBinaryExpr(opFunc opFunction, expr *parser.Expression) (amt Amount, err error) {
	left, err := eval(expr.Left)
	if err != nil {
		return amt, err
	}
	right, err := eval(expr.Right)
	if err != nil {
		return amt, err
	}
	if left.Units != expr.TargetUnits.String() {
		left, err = convertUnits(left, expr.TargetUnits)
		if err != nil {
			return amt, err
		}
	}
	if right.Units != expr.TargetUnits.String() {
		right, err = convertUnits(right, expr.TargetUnits)
		if err != nil {
			return amt, err
		}
	}
	result := Amount{opFunc(left.Value, right.Value), left.Units}
	return result, nil
}

func evalValue(expr *parser.Expression) (Amount, error) {
	amt := Amount{expr.Value, expr.ValueUnits.String()}
	if expr.ValueUnits != expr.TargetUnits {
		amt, err := convertUnits(amt, expr.TargetUnits)
		if err != nil {
			return amt, err
		}
	}
	return amt, nil
}

func evalConversion(expr *parser.Expression) (amt Amount, err error) {
	amt, err = eval(expr.Left)
	if err != nil {
		return amt, err
	}
	if amt.Units != expr.TargetUnits.String() {
		amt, err = convertUnits(amt, expr.TargetUnits)
		if err != nil {
			return amt, err
		}
	}
	return amt, nil
}

func convertUnits(from Amount, units parser.Units) (Amount, error) {
	var amt Amount
	fromUnits, err := u.Find(from.Units)
	if err != nil {
		return amt, err
	}
	toUnits, err := u.Find(units.String())
	if err != nil {
		return amt, err
	}
	if fromUnits.Name == toUnits.Name {
		// no conversion is necessary; for example 2 kilograms in kg requires no conversion
		amt = Amount{from.Value, units.String()}
		return amt, nil
	}
	value, err := u.ConvertFloat(from.Value, fromUnits, toUnits)
	if err != nil {
		return amt, err
	}
	amt = Amount{value.Float(), units.String()}
	return amt, nil
}
