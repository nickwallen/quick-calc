package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSum(t *testing.T) {
	reader := bytes.NewBufferString("23 kg + 23 kg")
	writer := bytes.NewBufferString("")
	parse(reader, writer)
	assert.Equal(t, "\n > 46.00 kg \n", writer.String())
}

func TestSumDiffUnits(t *testing.T) {
	reader := bytes.NewBufferString("2 kilograms + 2000 g")
	writer := bytes.NewBufferString("")
	parse(reader, writer)
	assert.Equal(t, "\n > 4.00 kilograms \n", writer.String())
}

//func TestConversion(t *testing.T) {
//	reader := bytes.NewBufferString("2.0 fahrenheit in celsius")
//	writer := bytes.NewBufferString("")
//	parse(reader, writer)
//	assert.Equal(t, "\n > 4.00 kilograms \n", writer.String())
//}
