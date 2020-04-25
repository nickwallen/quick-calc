package parser

import (
	"github.com/nickwallen/qcalc/internal/tokens"
	"github.com/stretchr/testify/assert"
	"testing"
)

func pounds() (pounds Units) {
	pounds, _ = UnitsOf("pounds")
	return pounds
}

func kilos() (kilos Units) {
	kilos, _ = UnitsOf("kg")
	return kilos
}

func grams() (grams Units) {
	grams, _ = UnitsOf("g")
	return grams
}

func ounces() (ounces Units) {
	ounces, _ = UnitsOf("ounces")
	return ounces
}

func TestUnitsOf(t *testing.T) {
	expected := "kilograms"
	units, err := UnitsOf(expected)
	assert.Equal(t, expected, units.String())
	assert.Nil(t, err)
}

func TestAmountOf(t *testing.T) {
	var expected float64 = 20
	amount := AmountOf(expected, kilos())
	assert.Equal(t, expected, amount.Value)
	assert.Equal(t, kilos(), amount.Units)
}

func TestSumOf(t *testing.T) {
	amount := AmountOf(20, kilos())
	sum := operatorOf(amount, amount, kilos(), tokens.Plus)
	assert.Equal(t, amount, sum.left)
	assert.Equal(t, amount, sum.right)
	assert.Equal(t, kilos(), sum.units)
}

func TestSum(t *testing.T) {
	// 20 kg + 20 kg = ? kg
	sum := operatorOf(AmountOf(20, kilos()), AmountOf(20, kilos()), kilos(), tokens.Plus)
	actual, err := sum.Evaluate()
	expected := AmountOf(40, kilos())
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestSumDifferentUnits(t *testing.T) {
	// 2 kg + 2000 g = ? kg
	sum := operatorOf(AmountOf(2, kilos()), AmountOf(2000, grams()), kilos(), tokens.Plus)
	actual, err := sum.Evaluate()
	expected := AmountOf(4, kilos())
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestSumThenConvert(t *testing.T) {
	// 2000 g + 2000 g = ? kg
	sum := operatorOf(AmountOf(2000, grams()), AmountOf(2000, grams()), kilos(), tokens.Plus)
	actual, err := sum.Evaluate()
	expected := AmountOf(4, kilos())
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestConvert(t *testing.T) {
	convert := unitConverterOf(AmountOf(2, kilos()), grams())
	actual, err := convert.Evaluate()
	expected := AmountOf(2000, grams())
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}
