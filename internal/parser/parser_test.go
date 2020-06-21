package parser

import (
	"github.com/nickwallen/quick-calc/internal/io"
	"github.com/nickwallen/quick-calc/internal/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseValue(t *testing.T) {
	expr := "23 pounds"
	input := io.NewTokenChannel(expr)
	go func() {
		input.WriteToken(types.Number.Token("23"))
		input.WriteToken(types.Units.Token("pounds"))
		input.WriteToken(types.EOF.Token(""))
	}()
	actual, err := Parse(&input)
	expected := types.NewValue(23, types.Units.Token("pounds"))
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestParseValueNoUnits(t *testing.T) {
	expr := "23"
	input := io.NewTokenChannel(expr)
	go func() {
		input.WriteToken(types.Number.TokenAt("23", 1))
		input.WriteToken(types.EOF.TokenAt("", 2))
	}()
	_, err := Parse(&input)
	assert.NotNil(t, err)
	assert.Equal(t, "reached end of input, but expected a unit", err.Error())
}

func TestParseValueNoNumber(t *testing.T) {
	input := io.NewTokenChannel("pounds")
	go func() {
		input.WriteToken(types.Units.TokenAt("pounds", 1))
		input.WriteToken(types.EOF.TokenAt("", 2))
	}()
	_, err := Parse(&input)
	assert.NotNil(t, err)
	assert.Equal(t, "got 'pounds', but expected a number", err.Error())
}

func TestParseBinaryAdd(t *testing.T) {
	expr := "23 kg + 23 pounds"
	input := io.NewTokenChannel(expr)
	go func() {
		input.WriteToken(types.Number.Token("23"))
		input.WriteToken(types.Units.Token("kg"))
		input.WriteToken(types.Plus.Token("+"))
		input.WriteToken(types.Number.Token("23"))
		input.WriteToken(types.Units.Token("pounds"))
		input.WriteToken(types.EOF.Token(""))
	}()
	actual, err := Parse(&input)
	expected := types.AdditionExpr(
		types.NewValue(23, types.Units.Token("kg")),
		types.NewValue(23, types.Units.Token("pounds")))
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestParseBinarySubtract(t *testing.T) {
	expr := "23 kg - 23 pounds"
	input := io.NewTokenChannel(expr)
	go func() {
		input.WriteToken(types.Number.Token("23"))
		input.WriteToken(types.Units.Token("kg"))
		input.WriteToken(types.Minus.Token("-"))
		input.WriteToken(types.Number.Token("23"))
		input.WriteToken(types.Units.Token("pounds"))
		input.WriteToken(types.EOF.Token(""))
	}()
	actual, err := Parse(&input)
	expected := types.SubtractionExpr(
		types.NewValue(23, types.Units.Token("kg")),
		types.NewValue(23, types.Units.Token("pounds")))
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestParseConversion(t *testing.T) {
	expr := "2 pounds in ounces"
	input := io.NewTokenChannel(expr)
	go func() {
		input.WriteToken(types.Number.Token("2"))
		input.WriteToken(types.Units.Token("pounds"))
		input.WriteToken(types.In.Token("in"))
		input.WriteToken(types.Units.Token("ounces"))
		input.WriteToken(types.EOF.Token(""))
	}()
	actual, err := Parse(&input)
	expected := types.UnitConversionExpr(
		types.NewValue(2, types.Units.Token("pounds")),
		types.Units.Token("ounces"))
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestParseBinarySumAndConvert(t *testing.T) {
	expr := "2 ounces + 2 pounds in pounds"
	input := io.NewTokenChannel(expr)
	go func() {
		input.WriteToken(types.Number.Token("2"))
		input.WriteToken(types.Units.Token("ounces"))
		input.WriteToken(types.Plus.Token("+"))
		input.WriteToken(types.Number.Token("2"))
		input.WriteToken(types.Units.Token("pounds"))
		input.WriteToken(types.In.Token("in"))
		input.WriteToken(types.Units.Token("pounds"))
		input.WriteToken(types.EOF.Token(""))
	}()
	actual, err := Parse(&input)
	expected := types.UnitConversionExpr(
		types.AdditionExpr(
			types.NewValue(2, types.Units.Token("ounces")),
			types.NewValue(2, types.Units.Token("pounds"))),
		types.Units.Token("pounds"))
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestParseBinarySubtractAndConvert(t *testing.T) {
	expr := "2 pounds - 2 ounces in ounces"
	input := io.NewTokenChannel(expr)
	go func() {
		input.WriteToken(types.Number.Token("2"))
		input.WriteToken(types.Units.Token("pounds"))
		input.WriteToken(types.Minus.Token("-"))
		input.WriteToken(types.Number.Token("2"))
		input.WriteToken(types.Units.Token("ounces"))
		input.WriteToken(types.In.Token("in"))
		input.WriteToken(types.Units.Token("ounces"))
		input.WriteToken(types.EOF.Token(""))
	}()
	actual, err := Parse(&input)
	expected := types.UnitConversionExpr(
		types.SubtractionExpr(
			types.NewValue(2, types.Units.Token("pounds")),
			types.NewValue(2, types.Units.Token("ounces"))),
		types.Units.Token("ounces"))
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestParseSum(t *testing.T) {
	expr := "2 ounces + 3 ounces + 4 ounces"
	input := io.NewTokenChannel(expr)
	go func() {
		input.WriteToken(types.Number.Token("2"))
		input.WriteToken(types.Units.Token("ounces"))
		input.WriteToken(types.Plus.Token("+"))
		input.WriteToken(types.Number.Token("3"))
		input.WriteToken(types.Units.Token("ounces"))
		input.WriteToken(types.Plus.Token("+"))
		input.WriteToken(types.Number.Token("4"))
		input.WriteToken(types.Units.Token("ounces"))
		input.WriteToken(types.EOF.Token(""))
	}()
	actual, err := Parse(&input)
	expected := types.AdditionExpr(
		types.AdditionExpr(
			types.NewValue(2, types.Units.Token("ounces")),
			types.NewValue(3, types.Units.Token("ounces"))),
		types.NewValue(4, types.Units.Token("ounces")))
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestParseSum2(t *testing.T) {
	expr := "2 ounces + 3 ounces - 4 ounces + 5 ounces"
	input := io.NewTokenChannel(expr)
	go func() {
		input.WriteToken(types.Number.Token("2"))
		input.WriteToken(types.Units.Token("ounces"))
		input.WriteToken(types.Plus.Token("+"))
		input.WriteToken(types.Number.Token("3"))
		input.WriteToken(types.Units.Token("ounces"))
		input.WriteToken(types.Minus.Token("-"))
		input.WriteToken(types.Number.Token("4"))
		input.WriteToken(types.Units.Token("ounces"))
		input.WriteToken(types.Plus.Token("+"))
		input.WriteToken(types.Number.Token("5"))
		input.WriteToken(types.Units.Token("ounces"))
		input.WriteToken(types.EOF.Token(""))
	}()
	actual, err := Parse(&input)
	expected := types.AdditionExpr(
		types.SubtractionExpr(
			types.AdditionExpr(
				types.NewValue(2, types.Units.Token("ounces")),
				types.NewValue(3, types.Units.Token("ounces"))),
			types.NewValue(4, types.Units.Token("ounces"))),
		types.NewValue(5, types.Units.Token("ounces")))
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}
