package io

import (
	"github.com/nickwallen/toks/internal/tokens"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTokenChannel_ReadWriteToken(t *testing.T) {
	tokenChannel := NewTokenChannel()
	token1 := tokens.Number.Token("1")
	token2 := tokens.Number.Token("2")

	// write some tokens
	go func() {
		err := tokenChannel.WriteToken(token1)
		assert.Nil(t, err)
		err = tokenChannel.WriteToken(token2)
		assert.Nil(t, err)
	}()

	// read the first token
	actual1, err := tokenChannel.ReadToken()
	assert.Equal(t, token1, actual1)
	assert.Nil(t, err)

	// read the second token
	actual2, err := tokenChannel.ReadToken()
	assert.Equal(t, token2, actual2)
	assert.Nil(t, err)

	//// there should be no tokens left
	//_, err = tokenChannel.ReadToken()
	//assert.Equal(t, "no tokens left", err.Error())
}

func TestTokenChannel_UnreadToken(t *testing.T) {
	tokenChannel := NewTokenChannel()

	// write a token
	token1 := tokens.Number.Token("1")
	go tokenChannel.WriteToken(token1)

	// read a token
	actual1, err := tokenChannel.ReadToken()
	assert.Equal(t, token1, actual1)
	assert.Nil(t, err)

	//// there should be no tokens left
	//_, err = tokenChannel.ReadToken()
	//assert.Equal(t, "no tokens left", err.Error())

	// unread the last token
	err = tokenChannel.UnreadToken(actual1)
	assert.Nil(t, err)

	// read the token again
	actual1, err = tokenChannel.ReadToken()
	assert.Equal(t, token1, actual1)
	assert.Nil(t, err)
}
