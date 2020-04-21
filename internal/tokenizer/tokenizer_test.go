package tokenizer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTokenize(t *testing.T) {
	output := make(chan Token, 2)
	go Tokenize("2 grams + 2 pounds", output)
	assert.Equal(t, Number.Token("2"), <-output)
	assert.Equal(t, Units.Token("grams"), <-output)
	assert.Equal(t, Plus.Token("+"), <-output)
	assert.Equal(t, Number.Token("2"), <-output)
	assert.Equal(t, Units.Token("pounds"), <-output)
	assert.Equal(t, EOF.Token(""), <-output)
}

func TestTokenizeToSlice(t *testing.T) {
	actual := TokenizeToSlice("2 grams + 2 pounds")
	expected := []Token{
		Number.Token("2"),
		Units.Token("grams"),
		Plus.Token("+"),
		Number.Token("2"),
		Units.Token("pounds"),
		EOF.Token("")}
	assert.ElementsMatch(t, expected, actual)
}
