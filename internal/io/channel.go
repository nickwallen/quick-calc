package io

import (
	"fmt"
	"github.com/nickwallen/toks/internal/tokens"
)

// TokenChannel Enables tokens to be read from and written to a channel.
type TokenChannel chan tokens.Token

// NewTokenChannel Creates a new token channel.
func NewTokenChannel() TokenChannel {
	return make(chan tokens.Token, 2)
}

// ReadToken Reads tokens from a channel.
func (ch TokenChannel) ReadToken() (tokens.Token, error) {
	var token tokens.Token
	token, ok := <-ch
	if !ok {
		return token, fmt.Errorf("no more tokens; channel is closed")
	}
	return token, nil
}

// WriteToken Writes tokens to a channel.
func (ch TokenChannel) WriteToken(token tokens.Token) error {
	ch <- token
	return nil
}

// Close Closes the token channel.
func (ch TokenChannel) Close() {
	close(ch)
}
