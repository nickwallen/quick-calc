package toks

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var testCases = map[string]string{
	"2 kg + 2000g":              "4.00 kg",
	"2 kilograms + 2 kilograms": "4.00 kilograms",
	"2 pounds + 2 kilograms":    "6.41 pounds",
	"2 feet - 2 feet":           "0.00 feet",
	"2 meters - 2 feet":         "1.39 meters",
}

func TestCalculator(t *testing.T) {
	for expression, expected := range testCases {
		actual, err := Calculate(expression)
		assert.Equal(t, expected, actual)
		assert.Nil(t, err)
	}
}
