package parser

import (
	"github.com/nickwallen/toks/internal/tokenizer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseAmount(t *testing.T) {
	actual, err := parse(tokenizer.Number.Token("23"), tokenizer.Units.Token("pounds"), tokenizer.EOF.Token(""))
	expected := AmountOf(23, pounds())
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestParseAmountNoUnits(t *testing.T) {
	_, err := parse(tokenizer.Number.Token("23"), tokenizer.EOF.Token(""))
	assert.Equal(t, "expected units, but reached end of input", err.Error())
}

func TestParseAmountNoNumber(t *testing.T) {
	_, err := parse(tokenizer.Units.Token("pounds"), tokenizer.EOF.Token(""))
	assert.Equal(t, "expected number, but got 'pounds'", err.Error())
}

func TestParseSum(t *testing.T) {
	actual, err := parse(
		tokenizer.Number.Token("23"),
		tokenizer.Units.Token("kg"),
		tokenizer.Plus.Token("+"),
		tokenizer.Number.Token("23"),
		tokenizer.Units.Token("pounds"),
		tokenizer.EOF.Token(""))
	expected := operatorOf(AmountOf(23, kilos()), AmountOf(23, pounds()), kilos(), tokenizer.Plus)
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestParseSubtract(t *testing.T) {
	actual, err := parse(
		tokenizer.Number.Token("23"),
		tokenizer.Units.Token("kg"),
		tokenizer.Minus.Token("-"),
		tokenizer.Number.Token("23"),
		tokenizer.Units.Token("pounds"),
		tokenizer.EOF.Token(""))
	expected := operatorOf(AmountOf(23, kilos()), AmountOf(23, pounds()), kilos(), tokenizer.Minus)
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestParseConversion(t *testing.T) {
	actual, err := parse(
		tokenizer.Number.Token("2"),
		tokenizer.Units.Token("pounds"),
		tokenizer.In.Token("in"),
		tokenizer.Units.Token("ounces"),
		tokenizer.EOF.Token(""))
	expected := unitConverterOf(AmountOf(2, pounds()), ounces())
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestParseSumAndConvert(t *testing.T) {
	actual, err := parse(
		tokenizer.Number.Token("2"),
		tokenizer.Units.Token("ounces"),
		tokenizer.Plus.Token("+"),
		tokenizer.Number.Token("2"),
		tokenizer.Units.Token("pounds"),
		tokenizer.In.Token("in"),
		tokenizer.Units.Token("pounds"),
		tokenizer.EOF.Token(""))
	expected := operatorOf(AmountOf(2, ounces()), AmountOf(2, pounds()), pounds(), tokenizer.Plus)
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestParseSubtractAndConvert(t *testing.T) {
	actual, err := parse(
		tokenizer.Number.Token("2"),
		tokenizer.Units.Token("pounds"),
		tokenizer.Minus.Token("-"),
		tokenizer.Number.Token("2"),
		tokenizer.Units.Token("ounces"),
		tokenizer.In.Token("in"),
		tokenizer.Units.Token("ounces"),
		tokenizer.EOF.Token(""))
	expected := operatorOf(AmountOf(2, pounds()), AmountOf(2, ounces()), ounces(), tokenizer.Minus)
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func parse(input ...tokenizer.Token) (Expression, error) {
	// create channel from which the parser will read tokens
	tokenChan := make(chan tokenizer.Token, len(input))
	for _, token := range input {
		tokenChan <- token
	}
	// parse the tokens
	return Parse(tokenChan)
}
