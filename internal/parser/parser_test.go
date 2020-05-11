package parser

import (
	"github.com/nickwallen/quick-calc/internal/io"
	"github.com/nickwallen/quick-calc/internal/tokens"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseAmount(t *testing.T) {
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

func TestParseAmountNoUnits(t *testing.T) {
	input := io.NewTokenChannel()
	go func() {
		input.WriteToken(tokens.Number.Token("23"))
		input.WriteToken(tokens.EOF.Token(""))
	}()
	_, err := Parse(&input)
	assert.Equal(t, "expected units, but reached end of input", err.Error())
}

func TestParseAmountNoNumber(t *testing.T) {
	input := io.NewTokenChannel()
	go func() {
		input.WriteToken(tokens.Units.Token("pounds"))
		input.WriteToken(tokens.EOF.Token(""))
	}()
	_, err := Parse(&input)
	assert.Equal(t, "expected number, but got 'pounds'", err.Error())
}

func TestParseSum(t *testing.T) {
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
	expected := binaryExpr(valueExpr(23, kilos()), valueExpr(23, pounds()), kilos(), tokens.Plus)
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestParseSubtract(t *testing.T) {
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
	expected := binaryExpr(valueExpr(23, kilos()), valueExpr(23, pounds()), kilos(), tokens.Minus)
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

func TestParseSumAndConvert(t *testing.T) {
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
	expected := binaryExpr(valueExpr(2, ounces()), valueExpr(2, pounds()), pounds(), tokens.Plus)
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestParseSubtractAndConvert(t *testing.T) {
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
	expected := binaryExpr(valueExpr(2, pounds()), valueExpr(2, ounces()), ounces(), tokens.Minus)
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
