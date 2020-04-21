package toks

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var testCases = map[string]string{
	"2 kg + 2000g":                        "4.00 kg",
	"2 kilograms + 2 kilograms":           "4.00 kilograms",
	"2 pounds + 2 kilograms":              "6.41 pounds",
	"2 feet - 2 feet":                     "0.00 feet",
	"2 meters - 2 feet":                   "1.39 meters",
	"2 pounds in ounces":                  "32.00 ounces",
	"2 pounds + 2 kilograms in kilograms": "2.91 kilograms",
	"2 meters - 2 feet in feet":           "4.56 feet",
}

func TestCalculator(t *testing.T) {
	for expression, expected := range testCases {
		actual, err := Calculate(expression)
		assert.Equal(t, expected, actual)
		assert.Nil(t, err)
	}
}
