package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasic(t *testing.T) {
	reader := bytes.NewBufferString("2 + 2")
	writer := bytes.NewBufferString("")
	tokenize(reader, writer)
	assert.Equal(t, "\n > NUM[2]  SYM[+]  NUM[2]  EOF  ", writer.String())
}
