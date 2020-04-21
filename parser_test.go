package toks

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseAmount(t *testing.T) {
	actual, err := Parse("23 pounds")
	if err != nil {
		assert.Fail(t, "unexpected error: %s", err)
	}
	expected := amountOf(23, pounds())
	assert.Equal(t, expected, actual)
}

func TestParseAmountNoUnits(t *testing.T) {
	_, err := Parse("23")
	assert.Equal(t, "expected units, like kilograms, but got end-of-line", err.Error())
}

func TestParseAmountNoNumber(t *testing.T) {
	_, err := Parse(" pounds ")
	assert.Equal(t, "expected number, but got 'p'", err.Error())
}

func TestParseSum(t *testing.T) {
	actual, err := Parse("23 kg + 23 pounds")
	if err != nil {
		assert.Fail(t, "unexpected error: %s", err)
	}
	expected := operatorOf(amountOf(23, kilos()), amountOf(23, pounds()), kilos(), Plus)
	assert.Equal(t, expected, actual)
}

func TestParseSubtract(t *testing.T) {
	actual, err := Parse("23 kg - 23 pounds")
	if err != nil {
		assert.Fail(t, "unexpected error: %s", err)
	}
	expected := operatorOf(amountOf(23, kilos()), amountOf(23, pounds()), kilos(), Minus)
	assert.Equal(t, expected, actual)
}

func TestParseConversion(t *testing.T) {
	actual, err := Parse("2 pounds in ounces")
	if err != nil {
		assert.Fail(t, "unexpected error: %s", err)
	}
	expected := unitConverterOf(amountOf(2, pounds()), ounces())
	assert.Equal(t, expected, actual)
}

func TestParseSumAndConvert(t *testing.T) {
	actual, err := Parse("2 ounces + 2 pounds in pounds")
	if err != nil {
		assert.Fail(t, "unexpected error: %s", err)
	}
	expected := operatorOf(amountOf(2, ounces()), amountOf(2, pounds()), pounds(), Plus)
	assert.Equal(t, expected, actual)
}

func TestParseSubtractAndConvert(t *testing.T) {
	actual, err := Parse("2 pounds - 2 ounces in ounces")
	if err != nil {
		assert.Fail(t, "unexpected error: %s", err)
	}
	expected := operatorOf(amountOf(2, pounds()), amountOf(2, ounces()), ounces(), Minus)
	assert.Equal(t, expected, actual)
}
