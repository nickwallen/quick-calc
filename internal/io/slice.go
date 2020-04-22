package io

import (
	"fmt"
	"github.com/nickwallen/toks/internal/tokens"
)

// TokenSlice Enables tokens to be read from and written to a slice.
type TokenSlice struct {
	slice    []tokens.Token
	position int
}

// NewTokenSlice Creates a new token slice.
func NewTokenSlice(capacity int) TokenSlice {
	return TokenSlice{
		slice:    make([]tokens.Token, capacity),
		position: 0,
	}
}

// TokenSliceOf Create a new slice containing a set of tokens.
func TokenSliceOf(input ...tokens.Token) TokenSlice {
	return TokenSlice{slice: input, position: 0}
}

// Tokens Returns the slice of tokens.
func (t *TokenSlice) Tokens() []tokens.Token {
	return t.slice
}

// ReadToken Allows tokens to be read from a slice.
func (t *TokenSlice) ReadToken() (tokens.Token, error) {
	var token tokens.Token
	if t.position >= len(t.slice) {
		return token, fmt.Errorf("no tokens left")
	}
	token = t.slice[t.position]
	t.position++
	return token, nil
}

// WriteToken Writes tokens to a slice.
func (t *TokenSlice) WriteToken(token tokens.Token) error {
	if t.position >= len(t.slice) {
		return fmt.Errorf("no space left")
	}
	t.slice[t.position] = token
	t.position++
	return nil
}

// Close Close the token slice.
func (t *TokenSlice) Close() {
	// nothing to do
}
