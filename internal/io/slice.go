package io

import (
	"fmt"
	"github.com/nickwallen/toks/internal/tokens"
)

// TokenSlice Enables tokens to be read from and written to a slice.
type TokenSlice struct {
	slice   []tokens.Token
	writeAt int
	readAt  int
}

// NewTokenSlice Creates a new token slice.
func NewTokenSlice(size int) TokenSlice {
	slice := make([]tokens.Token, size)
	return TokenSlice{slice, 0, 0}
}

// TokenSliceOf Create a new slice containing a set of tokens.
func TokenSliceOf(input ...tokens.Token) (TokenSlice, error) {
	slice := NewTokenSlice(len(input))
	for _, token := range input {
		err := slice.WriteToken(token)
		if err != nil {
			return slice, err
		}
	}
	return slice, nil
}

// Tokens Returns a slice containing the tokens.
func (t *TokenSlice) Tokens() []tokens.Token {
	return t.slice
}

// ReadToken Allows tokens to be read from a slice.
func (t *TokenSlice) ReadToken() (tokens.Token, error) {
	var token tokens.Token
	if len(t.slice) == 0 {
		return token, fmt.Errorf("no tokens left")
	}
	token = t.slice[0]
	t.slice = t.slice[1:]
	return token, nil
}

// UnreadToken Puts a token back so that it can be read again.
func (t *TokenSlice) UnreadToken(token tokens.Token) error {
	t.slice = append([]tokens.Token{token}, t.slice...)
	return nil
}

// WriteToken Writes tokens to a slice.
func (t *TokenSlice) WriteToken(token tokens.Token) error {
	t.slice = append(t.slice, token)
	return nil
}

// Close Close the token slice.
func (t *TokenSlice) Close() {
	// nothing to do
}
