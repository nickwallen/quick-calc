package toks

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseAmount(t *testing.T) {
	actual, err := parse(Number.token("23"), Units.token("pounds"), EOF.token(""))
	expected := amountOf(23, pounds())
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestParseAmountNoUnits(t *testing.T) {
	_, err := parse(Number.token("23"), EOF.token(""))
	assert.Equal(t, "expected units, but reached end of input", err.Error())
}

func TestParseAmountNoNumber(t *testing.T) {
	_, err := parse(Units.token("pounds"), EOF.token(""))
	assert.Equal(t, "expected number, but got 'pounds'", err.Error())
}

func TestParseSum(t *testing.T) {
	actual, err := parse(
		Number.token("23"),
		Units.token("kg"),
		Plus.token("+"),
		Number.token("23"),
		Units.token("pounds"),
		EOF.token(""))
	expected := operatorOf(amountOf(23, kilos()), amountOf(23, pounds()), kilos(), Plus)
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestParseSubtract(t *testing.T) {
	actual, err := parse(
		Number.token("23"),
		Units.token("kg"),
		Minus.token("-"),
		Number.token("23"),
		Units.token("pounds"),
		EOF.token(""))
	expected := operatorOf(amountOf(23, kilos()), amountOf(23, pounds()), kilos(), Minus)
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestParseConversion(t *testing.T) {
	actual, err := parse(
		Number.token("2"),
		Units.token("pounds"),
		In.token("in"),
		Units.token("ounces"),
		EOF.token(""))
	expected := unitConverterOf(amountOf(2, pounds()), ounces())
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestParseSumAndConvert(t *testing.T) {
	actual, err := parse(
		Number.token("2"),
		Units.token("ounces"),
		Plus.token("+"),
		Number.token("2"),
		Units.token("pounds"),
		In.token("in"),
		Units.token("pounds"),
		EOF.token(""))
	expected := operatorOf(amountOf(2, ounces()), amountOf(2, pounds()), pounds(), Plus)
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestParseSubtractAndConvert(t *testing.T) {
	actual, err := parse(
		Number.token("2"),
		Units.token("pounds"),
		Minus.token("-"),
		Number.token("2"),
		Units.token("ounces"),
		In.token("in"),
		Units.token("ounces"),
		EOF.token(""))
	expected := operatorOf(amountOf(2, pounds()), amountOf(2, ounces()), ounces(), Minus)
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func parse(input ...Token) (Expression, error) {
	// create channel from which the parser will read tokens
	tokenChan := make(chan Token, len(input))
	for _, token := range input {
		tokenChan <- token
	}
	// parse the tokens
	return Parse(tokenChan)
}
