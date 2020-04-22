package io

import (
	"github.com/nickwallen/toks/internal/tokens"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTokenSlice_ReadToken(t *testing.T) {
	token1 := tokens.Number.Token("1")
	token2 := tokens.Number.Token("2")
	tokenSlice, err := TokenSliceOf(token1, token2)
	assert.Nil(t, err)

	// read the first token
	actual1, err := tokenSlice.ReadToken()
	assert.Equal(t, token1, actual1)
	assert.Nil(t, err)

	// read the second token
	actual2, err := tokenSlice.ReadToken()
	assert.Equal(t, token2, actual2)
	assert.Nil(t, err)

	// there should be no tokens left
	_, err = tokenSlice.ReadToken()
	assert.Equal(t, "no tokens left", err.Error())
}

func TestTokenSlice_ReadTokenNoWrite(t *testing.T) {
	// there are  no tokens tor ead
	tokenSlice := NewTokenSlice(2)
	_, err := tokenSlice.ReadToken()
	assert.Equal(t, "no tokens left", err.Error())
}

func TestTokenSlice_WriteToken(t *testing.T) {
	tokenSlice := NewTokenSlice(2)

	// write a token
	token1 := tokens.Number.Token("1")
	err := tokenSlice.WriteToken(token1)
	assert.Nil(t, err)

	// read the token
	actual1, err := tokenSlice.ReadToken()
	assert.Equal(t, token1, actual1)
	assert.Nil(t, err)

	// write another token, but should be at capacity
	token2 := tokens.Number.Token("2")
	err = tokenSlice.WriteToken(token2)
	assert.Equal(t, "no space left", err.Error())

	// read the token
	actual2, err := tokenSlice.ReadToken()
	assert.Equal(t, token2, actual2)
	assert.Nil(t, err)

	// write another token, but should be at capacity
	token3 := tokens.Number.Token("3")
	err = tokenSlice.WriteToken(token3)
	assert.Equal(t, "no space left", err.Error())
}

func TestTokenSlice_UnreadToken(t *testing.T) {
	token1 := tokens.Number.Token("1")
	tokenSlice, err := TokenSliceOf(token1)
	assert.Nil(t, err)

	// read a token
	actual1, err := tokenSlice.ReadToken()
	assert.Equal(t, token1, actual1)
	assert.Nil(t, err)

	// there should be no tokens left
	_, err = tokenSlice.ReadToken()
	assert.Equal(t, "no tokens left", err.Error())

	// unread the last token
	err = tokenSlice.UnreadToken(actual1)
	assert.Nil(t, err)

	// read the token again
	actual1, err = tokenSlice.ReadToken()
	assert.Equal(t, token1, actual1)
	assert.Nil(t, err)
}
