package util

import (
	"fmt"
	"github.com/nickwallen/toks/internal/tokenizer"
)

// TokenChannel allows the parser and tokenizer to read/write to channels.
type TokenChannel chan tokenizer.Token

// ReadToken Allows the parser to read from a channel of tokens.
func (ch TokenChannel) ReadToken() (tokenizer.Token, error) {
	var token tokenizer.Token
	token, ok := <-ch
	if !ok {
		return token, fmt.Errorf("no more tokens; channel is closed")
	}
	return token, nil
}

// WriteToken Allows the tokenizer to write to a channel of tokens.
func (ch TokenChannel) WriteToken(token tokenizer.Token) error {
	ch <- token
	return nil
}

// Close Closes the channel
func (ch TokenChannel) Close() {
	close(ch)
}
