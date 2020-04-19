package toks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNextToken(t *testing.T) {
	tok := New("2 + 2")
	assert.Equal(t, Number.Token("2"), tok.NextToken())
	assert.Equal(t, Plus.Token("+"), tok.NextToken())
	assert.Equal(t, Number.Token("2"), tok.NextToken())
	assert.Equal(t, EOF.Token(""), tok.NextToken())
}

func TestTokens(t *testing.T) {
	tok := New("2 + 2")
	assert.Equal(t, Number.Token("2"), <-tok.Tokens())
	assert.Equal(t, Plus.Token("+"), <-tok.Tokens())
	assert.Equal(t, Number.Token("2"), <-tok.Tokens())
	assert.Equal(t, EOF.Token(""), <-tok.Tokens())
}
