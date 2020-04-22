package parser

import (
	"github.com/nickwallen/toks/internal/io"
	"github.com/nickwallen/toks/internal/tokens"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseAmount(t *testing.T) {
	input := io.TokenSliceOf(tokens.Number.Token("23"), tokens.Units.Token("pounds"), tokens.EOF.Token(""))
	actual, err := Parse(&input)
	expected := AmountOf(23, pounds())
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestParseAmountNoUnits(t *testing.T) {
	input := io.TokenSliceOf(tokens.Number.Token("23"), tokens.EOF.Token(""))
	_, err := Parse(&input)
	assert.Equal(t, "expected units, but reached end of input", err.Error())
}

func TestParseAmountNoNumber(t *testing.T) {
	input := io.TokenSliceOf(tokens.Units.Token("pounds"), tokens.EOF.Token(""))
	_, err := Parse(&input)
	assert.Equal(t, "expected number, but got 'pounds'", err.Error())
}

func TestParseSum(t *testing.T) {
	input := io.TokenSliceOf(
		tokens.Number.Token("23"),
		tokens.Units.Token("kg"),
		tokens.Plus.Token("+"),
		tokens.Number.Token("23"),
		tokens.Units.Token("pounds"),
		tokens.EOF.Token(""))
	actual, err := Parse(&input)
	expected := operatorOf(AmountOf(23, kilos()), AmountOf(23, pounds()), kilos(), tokens.Plus)
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestParseSubtract(t *testing.T) {
	input := io.TokenSliceOf(
		tokens.Number.Token("23"),
		tokens.Units.Token("kg"),
		tokens.Minus.Token("-"),
		tokens.Number.Token("23"),
		tokens.Units.Token("pounds"),
		tokens.EOF.Token(""))
	actual, err := Parse(&input)
	expected := operatorOf(AmountOf(23, kilos()), AmountOf(23, pounds()), kilos(), tokens.Minus)
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestParseConversion(t *testing.T) {
	input := io.TokenSliceOf(
		tokens.Number.Token("2"),
		tokens.Units.Token("pounds"),
		tokens.In.Token("in"),
		tokens.Units.Token("ounces"),
		tokens.EOF.Token(""))
	actual, err := Parse(&input)
	expected := unitConverterOf(AmountOf(2, pounds()), ounces())
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestParseSumAndConvert(t *testing.T) {
	input := io.TokenSliceOf(
		tokens.Number.Token("2"),
		tokens.Units.Token("ounces"),
		tokens.Plus.Token("+"),
		tokens.Number.Token("2"),
		tokens.Units.Token("pounds"),
		tokens.In.Token("in"),
		tokens.Units.Token("pounds"),
		tokens.EOF.Token(""))
	actual, err := Parse(&input)
	expected := operatorOf(AmountOf(2, ounces()), AmountOf(2, pounds()), pounds(), tokens.Plus)
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestParseSubtractAndConvert(t *testing.T) {
	input := io.TokenSliceOf(
		tokens.Number.Token("2"),
		tokens.Units.Token("pounds"),
		tokens.Minus.Token("-"),
		tokens.Number.Token("2"),
		tokens.Units.Token("ounces"),
		tokens.In.Token("in"),
		tokens.Units.Token("ounces"),
		tokens.EOF.Token(""))
	actual, err := Parse(&input)
	expected := operatorOf(AmountOf(2, pounds()), AmountOf(2, ounces()), ounces(), tokens.Minus)
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}
