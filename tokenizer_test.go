package toks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenize(t *testing.T) {
	output := make(chan Token, 2)
	go Tokenize("2 grams + 2 pounds", output)
	assert.Equal(t, Number.token("2"), <-output)
	assert.Equal(t, Units.token("grams"), <-output)
	assert.Equal(t, Plus.token("+"), <-output)
	assert.Equal(t, Number.token("2"), <-output)
	assert.Equal(t, Units.token("pounds"), <-output)
	assert.Equal(t, EOF.token(""), <-output)
}
