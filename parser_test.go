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
	expected := SumOf(
		AmountOf(23, UnitsOf("kg")),
		AmountOf(23, UnitsOf("pounds")),
		UnitsOf("kg"))
	assert.Equal(t, expected, actual)
}
