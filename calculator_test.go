package calc_test

import (
	"fmt"
	"github.com/nickwallen/quick-calc"
	"github.com/stretchr/testify/assert"
	"testing"
)

var expressions = map[string]string{
	"2 oz":                                "2.00 oz",
	"45 lbs":                              "45.00 lbs",
	"2 kilograms in kg":                   "2.00 kilograms",
	"2 kg + 2000g":                        "4.00 kg",
	"2 kilograms + 2 kilograms":           "4.00 kilograms",
	"2 pounds + 2 kilograms":              "6.41 pounds",
	"2 feet - 2 feet":                     "0.00 feet",
	"2 meters - 2 feet":                   "1.39 meters",
	"2 pounds in ounces":                  "32.00 ounces",
	"2 pounds + 2 kilograms in kilograms": "2.91 kilograms",
	"2 meters - 2 feet in feet":           "4.56 feet",
	"2kg + 34g in grams":                  "2034.00 grams",
	"2 miles + 2 meters in feet":          "10566.56 feet",
	"12 years in days":                    "4383.00 days",
	"12 mmmH2O + 12 mmmH2O":               "24.00 mmmH2O",
}

func TestCalculate(t *testing.T) {
	for expression, expected := range expressions {
		actual, err := calc.Calculate(expression)
		assert.Equal(t, expected, actual)
		assert.Nil(t, err)
	}
}

func TestCalculateAmount(t *testing.T) {
	for expression, expected := range expressions {
		actual, err := calc.CalculateAmount(expression)
		actualStr := fmt.Sprintf("%.2f %s", actual.Value, actual.Units)
		assert.Equal(t, expected, actualStr)
		assert.Nil(t, err)
	}
}

var badExpressions = []string{
	"2 miles + 3 pounds",
	"32 googles",
}

func TestCalculateBadExpr(t *testing.T) {
	for _, expression := range badExpressions {
		_, err := calc.Calculate(expression)
		assert.NotNil(t, err)
	}
}
