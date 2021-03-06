package calc_test

import (
	"fmt"
	"github.com/nickwallen/quick-calc"
	"github.com/stretchr/testify/assert"
	"testing"
)

var expressions = map[string]string{
	"2 oz":                                        "2.00 oz",
	"45 lbs":                                      "45.00 lbs",
	"2 kilograms in kg":                           "2.00 kg",
	"2 kg + 2000g":                                "4.00 kg",
	"2 kilograms + 2 kilograms":                   "4.00 kilograms",
	"2 pounds + 2 kilograms":                      "6.41 pounds",
	"2 feet - 2 feet":                             "0.00 feet",
	"2 meters - 2 feet":                           "1.39 meters",
	"2 pounds in ounces":                          "32.00 ounces",
	"2 pounds + 2 kilograms in kilograms":         "2.91 kilograms",
	"2 meters - 2 feet in feet":                   "4.56 feet",
	"2kg + 34g in grams":                          "2034.00 grams",
	"2 miles + 2 meters in feet":                  "10566.56 feet",
	"12 years in days":                            "4383.00 days",
	"12 mmmH2O + 12 mmmH2O":                       "24.00 mmmH2O",
	"2 oz + 3 oz + 4 oz + 5 oz":                   "14.00 oz",
	"2 oz - 3 oz + 4 oz + 5 oz":                   "8.00 oz",
	"2 oz + 3 oz - 4 oz + 5 oz":                   "6.00 oz",
	"2 oz + 3 oz + 4 oz - 5 oz":                   "4.00 oz",
	"2 oz + 3 oz + 4 oz - 5 oz in pounds":         "0.25 pounds",
	"12 oz + 1.2 lbs + 24 oz - 0.8 lbs in pounds": "2.65 pounds",
}

var badExpressions = map[string]string{
	"32 googles":         "'googles' is not a known measurement unit",
	"2 miles / 500 feet": "got '/', but expected '+', '-', 'in'",
	"2 miles + 3 pounds": "cannot convert from pounds to miles",
	"pounds":             "got 'pounds', but expected a number",
}

func TestCalculate(t *testing.T) {
	for input, expected := range expressions {
		t.Run(input, func(t *testing.T) {
			actual, err := calc.Calculate(input)
			assert.Nil(t, err, input)
			assert.Equal(t, expected, actual, input)
		})
	}
}

func TestCalculateAmount(t *testing.T) {
	for input, expected := range expressions {
		t.Run(input, func(t *testing.T) {
			actual, err := calc.CalculateAmount(input)
			actualStr := fmt.Sprintf("%.2f %s", actual.Value, actual.Units.Value)
			assert.Nil(t, err)
			assert.Equal(t, expected, actualStr)
		})
	}
}

func TestCalculateBadExpr(t *testing.T) {
	for input, expectedErr := range badExpressions {
		t.Run(input, func(t *testing.T) {
			_, err := calc.Calculate(input)
			assert.NotNil(t, err)
			assert.Equal(t, expectedErr, fmt.Sprintf("%s", err))
		})
	}
}
