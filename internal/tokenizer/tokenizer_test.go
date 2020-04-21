package tokenizer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TokenChannel allows the parser and tokenizer to read/write to channels.
type TokenChannel chan Token

// WriteToken Allows the tokenizer to write to a channel of tokens.
func (ch TokenChannel) WriteToken(token Token) error {
	ch <- token
	return nil
}

// Close Closes the channel
func (ch TokenChannel) Close() {
	close(ch)
}

func TestTokenize(t *testing.T) {
	var output TokenChannel
	output = make(chan Token, 2)
	go Tokenize("2 grams + 2 pounds", output)
	assert.Equal(t, Number.Token("2"), <-output)
	assert.Equal(t, Units.Token("grams"), <-output)
	assert.Equal(t, Plus.Token("+"), <-output)
	assert.Equal(t, Number.Token("2"), <-output)
	assert.Equal(t, Units.Token("pounds"), <-output)
	assert.Equal(t, EOF.Token(""), <-output)
}
