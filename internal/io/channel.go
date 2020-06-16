package io

import (
	"fmt"
	"github.com/nickwallen/quick-calc/internal/types"
)

// TokenChannel Enables tokens to be read from and written to a channel.
type TokenChannel struct {
	channel chan types.Token // the channel on which tokens are published
	input   string           // the original input to be tokenized
}

// NewTokenChannel Creates a new token channel.
func NewTokenChannel(input string) TokenChannel {
	return TokenChannel{
		channel: make(chan types.Token, 2),
		input:   input,
	}
}

// ReadToken Reads tokens from a channel.
func (t TokenChannel) ReadToken() (types.Token, error) {
	var token types.Token
	token, ok := <-t.channel
	if !ok {
		return token, fmt.Errorf("no more tokens; channel is closed")
	}
	return token, nil
}

// WriteToken Writes tokens to a channel.
func (t TokenChannel) WriteToken(token types.Token) error {
	t.channel <- token
	return nil
}

// Close Closes the token channel.
func (t TokenChannel) Close() {
	t.channel <- types.EOF.TokenAt("", len(t.input))
	close(t.channel)
}

// Input returns the original input that was tokenized.
func (t TokenChannel) Input() string {
	return t.input
}
