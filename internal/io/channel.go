package io

import (
	"fmt"
	"github.com/nickwallen/toks/internal/tokens"
)

// TokenChannel Enables tokens to be read from and written to a channel.
type TokenChannel struct {
	channel    chan tokens.Token
	undoBuffer TokenSlice
}

// NewTokenChannel Creates a new token channel.
func NewTokenChannel() TokenChannel {
	return TokenChannel{
		channel:    make(chan tokens.Token, 2),
		undoBuffer: NewTokenSlice(2),
	}
}

// Tokens Returns the token channel.
func (ch *TokenChannel) Tokens() chan tokens.Token {
	return ch.channel
}

// ReadToken Reads tokens from a channel.
func (ch *TokenChannel) ReadToken() (tokens.Token, error) {
	var token tokens.Token

	// is there anything in the undo buffer?
	token, err := ch.undoBuffer.ReadToken()
	if err == nil {
		return token, nil
	}

	token, ok := <-ch.channel
	if !ok {
		return token, fmt.Errorf("no more tokens; channel is closed")
	}
	return token, nil
}

// UnreadToken Puts a token back so that it can be read again.
func (ch *TokenChannel) UnreadToken(token tokens.Token) error {
	return ch.undoBuffer.WriteToken(token)
}

// WriteToken Writes tokens to a channel.
func (ch *TokenChannel) WriteToken(token tokens.Token) error {
	ch.channel <- token
	return nil
}

// Close Closes the token channel.
func (ch *TokenChannel) Close() {
	close(ch.channel)
}
