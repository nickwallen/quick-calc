package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCalculate(t *testing.T) {
	writer := bytes.NewBufferString("")
	calculate("23 kg + 23 kg", writer)
	assert.Equal(t, "46.00 kilograms \n", writer.String())
}

func TestTokenize(t *testing.T) {
	writer := bytes.NewBufferString("")
	tokenize("2 + 2", writer)
	assert.Equal(t, "NUM[2]  SYM[+]  NUM[2]  EOF  ", writer.String())
}
