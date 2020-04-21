package toks

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func pounds() (pounds amountUnits) {
	pounds, _ = unitsOf("pounds")
	return pounds
}

func kilos() (kilos amountUnits) {
	kilos, _ = unitsOf("kilograms")
	return kilos
}

func grams() (grams amountUnits) {
	grams, _ = unitsOf("grams")
	return grams
}

func ounces() (ounces amountUnits) {
	ounces, _ = unitsOf("ounces")
	return ounces
}

func TestUnitsOf(t *testing.T) {
	expected := "kilograms"
	units, err := unitsOf(expected)
	assert.Equal(t, expected, units.String())
	assert.Nil(t, err)
}

func TestAmountOf(t *testing.T) {
	var expected float64 = 20
	amount := amountOf(expected, kilos())
	assert.Equal(t, expected, amount.value)
	assert.Equal(t, kilos(), amount.units)
}

func TestSumOf(t *testing.T) {
	amount := amountOf(20, kilos())
	sum := operatorOf(amount, amount, kilos(), Plus)
	assert.Equal(t, amount, sum.left)
	assert.Equal(t, amount, sum.right)
	assert.Equal(t, kilos(), sum.units)
}

func TestSum(t *testing.T) {
	// 20 kg + 20 kg = ? kg
	sum := operatorOf(amountOf(20, kilos()), amountOf(20, kilos()), kilos(), Plus)
	actual, err := sum.Evaluate()
	expected := amountOf(40, kilos())
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestSumDifferentUnits(t *testing.T) {
	// 2 kg + 2000 g = ? kg
	sum := operatorOf(amountOf(2, kilos()), amountOf(2000, grams()), kilos(), Plus)
	actual, err := sum.Evaluate()
	expected := amountOf(4, kilos())
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestSumThenConvert(t *testing.T) {
	// 2000 g + 2000 g = ? kg
	sum := operatorOf(amountOf(2000, grams()), amountOf(2000, grams()), kilos(), Plus)
	actual, err := sum.Evaluate()
	expected := amountOf(4, kilos())
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestConvert(t *testing.T) {
	convert := unitConverterOf(amountOf(2, kilos()), grams())
	actual, err := convert.Evaluate()
	expected := amountOf(2000, grams())
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}
