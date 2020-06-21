package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValue_Eval(t *testing.T) {
	input := "22 lbs"
	number := float64(22)
	units := Units.TokenAt("lbs", 1)
	amount, err := NewValue(number, units).Eval(input)
	assert.Nil(t, err)
	assert.Equal(t, number, amount.Value)
	assert.Equal(t, units, amount.Units)
}

func TestValue_Eval_InvalidUnits(t *testing.T) {
	input := "22 googles"
	number := float64(22)
	units := Units.TokenAt("googles", 4)
	_, err := NewValue(number, units).Eval(input)
	assert.NotNil(t, err)
}

func TestAddition_Eval(t *testing.T) {
	input := "2 pounds + 1 stone"
	pounds := Units.TokenAt("pounds", 3)
	stones := Units.TokenAt("stone", 14)
	amount, err := AdditionExpr(NewValue(2, pounds), NewValue(1, stones)).Eval(input)
	assert.Nil(t, err)
	assert.InDelta(t, float64(16), amount.Value, 0.01)
	assert.Equal(t, pounds, amount.Units)
}

func TestAddition_Eval_InvalidUnitConversion(t *testing.T) {
	input := "2 miles + 1 hour"
	miles := Units.TokenAt("miles", 3)
	hours := Units.TokenAt("hour", 13)
	_, err := AdditionExpr(NewValue(2, miles), NewValue(1, hours)).Eval(input)
	if assert.NotNil(t, err) {
		_, ok := err.(*InvalidUnitConversion)
		assert.True(t, ok, "expected invalid unit conversion error, got %s", err)
	}
}

func TestAddition_Eval_InvalidUnits(t *testing.T) {
	input := "2 googles + 1 googles"
	googles := Units.TokenAt("googles", 3)
	_, err := AdditionExpr(NewValue(2, googles), NewValue(1, googles)).Eval(input)
	if assert.NotNil(t, err) {
		_, ok := err.(*InvalidUnits)
		assert.True(t, ok, "expected invalid units error, got ", err)
	}
}

func TestSubtraction_Eval(t *testing.T) {
	input := "2 stones - 1 pound"
	stones := Units.TokenAt("stone", 2)
	pounds := Units.TokenAt("pounds", 14)
	amount, err := SubtractionExpr(NewValue(2, stones), NewValue(1, pounds)).Eval(input)
	assert.Nil(t, err)
	assert.InDelta(t, 1.92, amount.Value, 0.01)
	assert.Equal(t, stones, amount.Units)
}

func TestSubtraction_Eval_InvalidUnitConversion(t *testing.T) {
	input := "2 miles - 1 hour"
	miles := Units.TokenAt("miles", 3)
	hours := Units.TokenAt("hour", 13)
	_, err := SubtractionExpr(NewValue(2, miles), NewValue(1, hours)).Eval(input)
	if assert.NotNil(t, err) {
		_, ok := err.(*InvalidUnitConversion)
		assert.True(t, ok, "expected invalid unit conversion error, got %s", err)
	}
}

func TestSubtraction_Eval_InvalidUnits(t *testing.T) {
	input := "2 googles - 1 googles"
	googles := Units.TokenAt("googles", 3)
	_, err := SubtractionExpr(NewValue(2, googles), NewValue(1, googles)).Eval(input)
	if assert.NotNil(t, err) {
		_, ok := err.(*InvalidUnits)
		assert.True(t, ok, "expected invalid units error, got ", err)
	}
}

func TestUnitConversion_Eval(t *testing.T) {
	input := "2 stones in pounds"
	stones := Units.TokenAt("stones", 3)
	pounds := Units.TokenAt("pounds", 13)
	amount, err := UnitConversionExpr(NewValue(2.0, stones), pounds).Eval(input)
	assert.Nil(t, err)
	assert.InDelta(t, 28, amount.Value, 0.01)
	assert.Equal(t, pounds, amount.Units)
}

func TestUnitConversion_Eval_InvalidUnitConversion(t *testing.T) {
	input := "2 hours in miles"
	hours := Units.TokenAt("hours", 3)
	miles := Units.TokenAt("miles", 12)
	_, err := UnitConversionExpr(NewValue(2.0, hours), miles).Eval(input)
	if assert.NotNil(t, err) {
		_, ok := err.(*InvalidUnitConversion)
		assert.True(t, ok, "expected invalid unit conversion error, got %s", err)
	}
}

func TestUnitConversion_Eval_InvalidUnits(t *testing.T) {
	input := "2 googles in miles"
	googles := Units.TokenAt("googles", 3)
	miles := Units.TokenAt("miles", 12)
	_, err := UnitConversionExpr(NewValue(2.0, googles), miles).Eval(input)
	if assert.NotNil(t, err) {
		_, ok := err.(*InvalidUnits)
		assert.True(t, ok, "expected invalid unit conversion error, got %s", err)
	}
}
