package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSum(t *testing.T) {
	reader := bytes.NewBufferString("23 kg + 23 kg")
	writer := bytes.NewBufferString("")
	calculate(reader, writer)
	assert.Equal(t, "\n > 46.00 kg \n", writer.String())
}
