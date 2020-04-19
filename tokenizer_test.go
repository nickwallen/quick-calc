package toks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNextToken(t *testing.T) {
	tok := New("2 + 2")
	assert.Equal(t, NumberToken.Of("2"), tok.NextToken())
	assert.Equal(t, PlusToken.Of("+"), tok.NextToken())
	assert.Equal(t, NumberToken.Of("2"), tok.NextToken())
	assert.Equal(t, EOFToken.Of(""), tok.NextToken())
}

func TestTokens(t *testing.T) {
	tok := New("2 + 2")
	assert.Equal(t, NumberToken.Of("2"), <-tok.Tokens())
	assert.Equal(t, PlusToken.Of("+"), <-tok.Tokens())
	assert.Equal(t, NumberToken.Of("2"), <-tok.Tokens())
	assert.Equal(t, EOFToken.Of(""), <-tok.Tokens())
}
