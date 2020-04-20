package toks

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnitsOf(t *testing.T) {
	expected := "kilograms"
	units := UnitsOf(expected)
	assert.Equal(t, expected, units.String())
}

func TestAmountOf(t *testing.T) {
	var expected float64 = 20
	units := UnitsOf("kilograms")
	amount := AmountOf(expected, units)
	assert.Equal(t, expected, amount.Value)
	assert.Equal(t, units, amount.Units)
}

func TestSumOf(t *testing.T) {
	units := UnitsOf("kilograms")
	amount := AmountOf(20, units)
	sum := SumOf(amount, amount, UnitsOf("kilograms"))
	assert.Equal(t, amount, sum.left)
	assert.Equal(t, amount, sum.right)
	assert.Equal(t, units, sum.units)
}

func TestSum(t *testing.T) {
	// 20 kg + 20 kg = ? kg
	sum := SumOf(
		AmountOf(20, UnitsOf("kg")),
		AmountOf(20, UnitsOf("kg")),
		UnitsOf("kg"))
	actual, err := sum.Evaluate()
	expected := AmountOf(40, UnitsOf("kg"))
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestSumDifferentUnits(t *testing.T) {
	// 2 kg + 2000 g = ? kg
	sum := SumOf(
		AmountOf(2, UnitsOf("kg")),
		AmountOf(2000, UnitsOf("g")),
		UnitsOf("kg"))
	actual, err := sum.Evaluate()
	expected := AmountOf(4, UnitsOf("kg"))
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestSumThenConvert(t *testing.T) {
	// 2000 g + 2000 g = ? kg
	sum := SumOf(
		AmountOf(2000, UnitsOf("g")),
		AmountOf(2000, UnitsOf("g")),
		UnitsOf("kg"))
	actual, err := sum.Evaluate()
	expected := AmountOf(4, UnitsOf("kg"))
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestConvert(t *testing.T) {
	convert := convOf(AmountOf(2, UnitsOf("kg")), UnitsOf("g"))
	actual, err := convert.Evaluate()
	expected := AmountOf(2000, UnitsOf("g"))
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}
