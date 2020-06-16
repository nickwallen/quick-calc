package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestCalculate(t *testing.T) {
	writer := bytes.NewBufferString("")
	calculate("23 kg + 23 kg", writer)
	assert.Equal(t, "46.00 kg \n", writer.String())
}

func TestTokenize(t *testing.T) {
	writer := bytes.NewBufferString("")
	tokenize("2 + 2", writer)
	assert.Equal(t, "NUM[2]  SYM[+]  NUM[2]  EOF  ", writer.String())
}

var badExpressions = map[string]string{
	"32 googles": `
error: 'googles' is not a known measurement unit at position 4
  |
  | 32 googles
  |    ^^^^^^^
	`,
	"2 miles / 500 feet": `
error: got '/', but expected '+', '-', 'in' at position 9
  |
  | 2 miles / 500 feet
  |         ^
	`,
	"2 miles + 3 pounds": `
error: cannot convert from pounds to miles at position 13
  |
  | 2 miles + 3 pounds
  |             ^^^^^^
`,
	"pounds": `
error: expected number, but got 'p' at position 1
  |
  | pounds
  | ^
	`,
}

func TestPrintError(t *testing.T) {
	for expr, expectedErr := range badExpressions {
		writer := bytes.NewBufferString("")
		calculate(expr, writer)
		assert.Equal(t, strings.TrimSpace(expectedErr), strings.TrimSpace(writer.String()))
	}
}
