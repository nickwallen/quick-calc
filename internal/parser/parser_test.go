package parser

import (
	"github.com/nickwallen/quick-calc/internal/io"
	"github.com/nickwallen/quick-calc/internal/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseValue(t *testing.T) {
	input := io.NewTokenChannel()
	go func() {
		input.WriteToken(types.Number.Token("23"))
		input.WriteToken(types.Units.Token("pounds"))
		input.WriteToken(types.EOF.Token(""))
	}()
	actual, err := Parse(&input)
	expected := types.NewValue(23, "pounds")
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestParseValueNoUnits(t *testing.T) {
	input := io.NewTokenChannel()
	go func() {
		input.WriteToken(types.Number.TokenAt("23", 1))
		input.WriteToken(types.EOF.TokenAt("", 2))
	}()
	_, err := Parse(&input)
	assert.Equal(t, "at position 2, reached end of input, but expected a unit", err.Error())
}

func TestParseValueNoNumber(t *testing.T) {
	input := io.NewTokenChannel()
	go func() {
		input.WriteToken(types.Units.TokenAt("pounds", 1))
		input.WriteToken(types.EOF.TokenAt("", 2))
	}()
	_, err := Parse(&input)
	assert.Equal(t, "at position 1, got 'pounds', but expected a number", err.Error())
}

func TestParseBinaryAdd(t *testing.T) {
	input := io.NewTokenChannel()
	go func() {
		input.WriteToken(types.Number.Token("23"))
		input.WriteToken(types.Units.Token("kg"))
		input.WriteToken(types.Plus.Token("+"))
		input.WriteToken(types.Number.Token("23"))
		input.WriteToken(types.Units.Token("pounds"))
		input.WriteToken(types.EOF.Token(""))
	}()
	actual, err := Parse(&input)
	expected := types.DoAddition(
		types.NewValue(23, "kg"),
		types.NewValue(23, "pounds"))
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestParseBinarySubtract(t *testing.T) {
	input := io.NewTokenChannel()
	go func() {
		input.WriteToken(types.Number.Token("23"))
		input.WriteToken(types.Units.Token("kg"))
		input.WriteToken(types.Minus.Token("-"))
		input.WriteToken(types.Number.Token("23"))
		input.WriteToken(types.Units.Token("pounds"))
		input.WriteToken(types.EOF.Token(""))
	}()
	actual, err := Parse(&input)
	expected := types.DoSubtraction(
		types.NewValue(23, "kg"),
		types.NewValue(23, "pounds"))
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestParseConversion(t *testing.T) {
	input := io.NewTokenChannel()
	go func() {
		input.WriteToken(types.Number.Token("2"))
		input.WriteToken(types.Units.Token("pounds"))
		input.WriteToken(types.In.Token("in"))
		input.WriteToken(types.Units.Token("ounces"))
		input.WriteToken(types.EOF.Token(""))
	}()
	actual, err := Parse(&input)
	expected := types.DoUnitConversion(types.NewValue(2, "pounds"), "ounces")
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestParseBinarySumAndConvert(t *testing.T) {
	input := io.NewTokenChannel()
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
	expected := types.DoUnitConversion(
		types.DoAddition(
			types.NewValue(2, "ounces"),
			types.NewValue(2, "pounds")),
		"pounds")
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestParseBinarySubtractAndConvert(t *testing.T) {
	input := io.NewTokenChannel()
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
	expected := types.DoUnitConversion(
		types.DoSubtraction(
			types.NewValue(2, "pounds"),
			types.NewValue(2, "ounces")),
		"ounces")
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestParseSum(t *testing.T) {
	input := io.NewTokenChannel()
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
	expected := types.DoAddition(
		types.DoAddition(
			types.NewValue(2, "ounces"),
			types.NewValue(3, "ounces")),
		types.NewValue(4, "ounces"))
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestParseSum2(t *testing.T) {
	input := io.NewTokenChannel()
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
	expected := types.DoAddition(
		types.DoSubtraction(
			types.DoAddition(
				types.NewValue(2, "ounces"),
				types.NewValue(3, "ounces")),
			types.NewValue(4, "ounces")),
		types.NewValue(5, "ounces"))
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}
