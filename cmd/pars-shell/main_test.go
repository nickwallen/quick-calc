package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBasic(t *testing.T) {
	reader := bytes.NewBufferString("23 kg + 23 kg")
	writer := bytes.NewBufferString("")
	parse(reader, writer)
	assert.Equal(t, "\n > 46 kg \n", writer.String())
}
