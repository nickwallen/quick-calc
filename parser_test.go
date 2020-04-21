package toks

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseAmount(t *testing.T) {
	actual, err := NewParser("23 pounds").Parse()
	if err != nil {
		assert.Fail(t, "unexpected error: %s", err)
	}
	expected := AmountOf(23, UnitsOf("pounds"))
	assert.Equal(t, expected, actual)
}

func TestParseAmountNoUnits(t *testing.T) {
	_, err := NewParser("23 ").Parse()
	assert.Equal(t, "expected units, like kilograms, but got end-of-line", err.Error())
}

func TestParseAmountNoNumber(t *testing.T) {
	_, err := NewParser(" pounds ").Parse()
	assert.Equal(t, "expected number, but got 'p'", err.Error())
}

func TestParseSum(t *testing.T) {
	actual, err := NewParser("23 kg + 23 pounds").Parse()
	if err != nil {
		assert.Fail(t, "unexpected error: %s", err)
	}
	expected := OperatorOf(
		AmountOf(23, UnitsOf("kg")),
		AmountOf(23, UnitsOf("pounds")),
		UnitsOf("kg"),
		Plus)
	assert.Equal(t, expected, actual)
}

func TestParseSubtract(t *testing.T) {
	actual, err := NewParser("23 kg - 23 pounds").Parse()
	if err != nil {
		assert.Fail(t, "unexpected error: %s", err)
	}
	expected := OperatorOf(
		AmountOf(23, UnitsOf("kg")),
		AmountOf(23, UnitsOf("pounds")),
		UnitsOf("kg"),
		Minus)
	assert.Equal(t, expected, actual)
}

func TestParseConversion(t *testing.T) {
	actual, err := NewParser("2 pounds in ounces").Parse()
	if err != nil {
		assert.Fail(t, "unexpected error: %s", err)
	}
	expected := UnitConverterOf(
		AmountOf(2, UnitsOf("pounds")),
		UnitsOf("ounces"))
	assert.Equal(t, expected, actual)
}

func TestParseSumAndConvert(t *testing.T) {
	actual, err := NewParser("2 ounces + 2 pounds in pounds").Parse()
	if err != nil {
		assert.Fail(t, "unexpected error: %s", err)
	}
	expected := OperatorOf(
		AmountOf(2, UnitsOf("ounces")),
		AmountOf(2, UnitsOf("pounds")),
		UnitsOf("pounds"),
		Plus)
	assert.Equal(t, expected, actual)
}

func TestParseSubtractAndConvert(t *testing.T) {
	actual, err := NewParser("2 pounds - 2 ounces in ounces").Parse()
	if err != nil {
		assert.Fail(t, "unexpected error: %s", err)
	}
	expected := OperatorOf(
		AmountOf(2, UnitsOf("pounds")),
		AmountOf(2, UnitsOf("ounces")),
		UnitsOf("ounces"),
		Minus)
	assert.Equal(t, expected, actual)
}
