package parser

import (
	"github.com/nickwallen/quick-calc/internal/io"
	"github.com/nickwallen/quick-calc/internal/tokens"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseValue(t *testing.T) {
	input := io.NewTokenChannel()
	go func() {
		input.WriteToken(tokens.Number.Token("23"))
		input.WriteToken(tokens.Units.Token("pounds"))
		input.WriteToken(tokens.EOF.Token(""))
	}()
	actual, err := Parse(&input)
	expected := valueExpr(23, pounds())
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestParseValueNoUnits(t *testing.T) {
	input := io.NewTokenChannel()
	go func() {
		input.WriteToken(tokens.Number.Token("23"))
		input.WriteToken(tokens.EOF.Token(""))
	}()
	_, err := Parse(&input)
	assert.Equal(t, "expected units, but reached end of input", err.Error())
}

func TestParseValueNoNumber(t *testing.T) {
	input := io.NewTokenChannel()
	go func() {
		input.WriteToken(tokens.Units.Token("pounds"))
		input.WriteToken(tokens.EOF.Token(""))
	}()
	_, err := Parse(&input)
	assert.Equal(t, "expected number, but got 'pounds'", err.Error())
}

func TestParseBinaryAdd(t *testing.T) {
	input := io.NewTokenChannel()
	go func() {
		input.WriteToken(tokens.Number.Token("23"))
		input.WriteToken(tokens.Units.Token("kg"))
		input.WriteToken(tokens.Plus.Token("+"))
		input.WriteToken(tokens.Number.Token("23"))
		input.WriteToken(tokens.Units.Token("pounds"))
		input.WriteToken(tokens.EOF.Token(""))
	}()
	actual, err := Parse(&input)
	expected := binaryExpr(
		tokens.Plus,
		valueExpr(23, kilos()),
		valueExpr(23, pounds()),
		kilos())
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestParseBinarySubtract(t *testing.T) {
	input := io.NewTokenChannel()
	go func() {
		input.WriteToken(tokens.Number.Token("23"))
		input.WriteToken(tokens.Units.Token("kg"))
		input.WriteToken(tokens.Minus.Token("-"))
		input.WriteToken(tokens.Number.Token("23"))
		input.WriteToken(tokens.Units.Token("pounds"))
		input.WriteToken(tokens.EOF.Token(""))
	}()
	actual, err := Parse(&input)
	expected := binaryExpr(
		tokens.Minus,
		valueExpr(23, kilos()),
		valueExpr(23, pounds()),
		kilos())
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestParseConversion(t *testing.T) {
	input := io.NewTokenChannel()
	go func() {
		input.WriteToken(tokens.Number.Token("2"))
		input.WriteToken(tokens.Units.Token("pounds"))
		input.WriteToken(tokens.In.Token("in"))
		input.WriteToken(tokens.Units.Token("ounces"))
		input.WriteToken(tokens.EOF.Token(""))
	}()
	actual, err := Parse(&input)
	expected := conversion(valueExpr(2, pounds()), ounces())
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestParseBinarySumAndConvert(t *testing.T) {
	input := io.NewTokenChannel()
	go func() {
		input.WriteToken(tokens.Number.Token("2"))
		input.WriteToken(tokens.Units.Token("ounces"))
		input.WriteToken(tokens.Plus.Token("+"))
		input.WriteToken(tokens.Number.Token("2"))
		input.WriteToken(tokens.Units.Token("pounds"))
		input.WriteToken(tokens.In.Token("in"))
		input.WriteToken(tokens.Units.Token("pounds"))
		input.WriteToken(tokens.EOF.Token(""))
	}()
	actual, err := Parse(&input)
	expected := binaryExpr(
		tokens.Plus,
		valueExpr(2, ounces()),
		valueExpr(2, pounds()),
		pounds())
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestParseBinarySubtractAndConvert(t *testing.T) {
	input := io.NewTokenChannel()
	go func() {
		input.WriteToken(tokens.Number.Token("2"))
		input.WriteToken(tokens.Units.Token("pounds"))
		input.WriteToken(tokens.Minus.Token("-"))
		input.WriteToken(tokens.Number.Token("2"))
		input.WriteToken(tokens.Units.Token("ounces"))
		input.WriteToken(tokens.In.Token("in"))
		input.WriteToken(tokens.Units.Token("ounces"))
		input.WriteToken(tokens.EOF.Token(""))
	}()
	actual, err := Parse(&input)
	expected := binaryExpr(
		tokens.Minus,
		valueExpr(2, pounds()),
		valueExpr(2, ounces()),
		ounces())
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestParseSum(t *testing.T) {
	input := io.NewTokenChannel()
	go func() {
		input.WriteToken(tokens.Number.Token("2"))
		input.WriteToken(tokens.Units.Token("ounces"))
		input.WriteToken(tokens.Plus.Token("+"))
		input.WriteToken(tokens.Number.Token("3"))
		input.WriteToken(tokens.Units.Token("ounces"))
		input.WriteToken(tokens.Plus.Token("+"))
		input.WriteToken(tokens.Number.Token("4"))
		input.WriteToken(tokens.Units.Token("ounces"))
		input.WriteToken(tokens.EOF.Token(""))
	}()
	actual, err := Parse(&input)
	expected := binaryExpr(
		tokens.Plus,
		valueExpr(2, ounces()),
		binaryExpr(
			tokens.Plus,
			valueExpr(3, ounces()),
			valueExpr(4, ounces()),
			ounces()),
		ounces())
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestParseSum2(t *testing.T) {
	input := io.NewTokenChannel()
	go func() {
		input.WriteToken(tokens.Number.Token("2"))
		input.WriteToken(tokens.Units.Token("ounces"))
		input.WriteToken(tokens.Plus.Token("+"))
		input.WriteToken(tokens.Number.Token("3"))
		input.WriteToken(tokens.Units.Token("ounces"))
		input.WriteToken(tokens.Plus.Token("-"))
		input.WriteToken(tokens.Number.Token("4"))
		input.WriteToken(tokens.Units.Token("ounces"))
		input.WriteToken(tokens.Plus.Token("+"))
		input.WriteToken(tokens.Number.Token("5"))
		input.WriteToken(tokens.Units.Token("ounces"))
		input.WriteToken(tokens.EOF.Token(""))
	}()
	actual, err := Parse(&input)
	expected := binaryExpr(
		tokens.Plus,
		valueExpr(2, ounces()),
		binaryExpr(
			tokens.Plus,
			valueExpr(3, ounces()),
			binaryExpr(
				tokens.Plus,
				valueExpr(4, ounces()),
				valueExpr(5, ounces()),
				ounces()),
			ounces()),
		ounces())
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func pounds() (pounds Units) {
	pounds, _ = UnitsOf("pounds")
	return pounds
}

func kilos() (kilos Units) {
	kilos, _ = UnitsOf("kg")
	return kilos
}

func ounces() (ounces Units) {
	ounces, _ = UnitsOf("ounces")
	return ounces
}
