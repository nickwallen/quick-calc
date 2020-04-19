package toks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecimals(t *testing.T) {
	expected := []Token{NumberToken.Of("22"), EOFToken.Of("")}
	expect(t, expected, New("22").Tokens())
	expect(t, expected, New("  22").Tokens())
	expect(t, expected, New("22  ").Tokens())
}

func TestBadDecimal(t *testing.T) {
	expected := []Token{ErrorToken.Of("expected number, but got \"2A\"")}
	expect(t, expected, New("2A").Tokens())
	expect(t, expected, New("  2A").Tokens())
	expect(t, expected, New("2A  ").Tokens())
}

func TestHexaDecimal(t *testing.T) {
	expected := []Token{NumberToken.Of("0xAF"), EOFToken.Of("")}
	expect(t, expected, New("0xAF").Tokens())
	expect(t, expected, New("  0xAF").Tokens())
	expect(t, expected, New("0xAF  ").Tokens())
}

func TestBadHexaDecimal(t *testing.T) {
	expected := []Token{ErrorToken.Of("expected number, but got \"0xG\"")}
	expect(t, expected, New("0xG2").Tokens())
	expect(t, expected, New("  0xG2").Tokens())
	expect(t, expected, New("0xG2  ").Tokens())
}

func TestFloats(t *testing.T) {
	expected := []Token{NumberToken.Of("2.22"), EOFToken.Of("")}
	expect(t, expected, New("2.22").Tokens())
	expect(t, expected, New("  2.22").Tokens())
	expect(t, expected, New("2.22  ").Tokens())
}

func TestExpNotation(t *testing.T) {
	expected := []Token{NumberToken.Of("2E10"), EOFToken.Of("")}
	expect(t, expected, New("2E10").Tokens())
	expect(t, expected, New("  2E10").Tokens())
	expect(t, expected, New("2E10  ").Tokens())
}

func TestPlus(t *testing.T) {
	expected := []Token{NumberToken.Of("2"), PlusToken.Of("+"), NumberToken.Of("2"), EOFToken.Of("")}
	expect(t, expected, New("2 + 2").Tokens())
	expect(t, expected, New("2+2").Tokens())
	expect(t, expected, New("   2 + 2").Tokens())
	expect(t, expected, New("2+2   ").Tokens())
}

func TestManyPlus(t *testing.T) {
	expected := []Token{NumberToken.Of("2"), PlusToken.Of("+"), ErrorToken.Of("expected number, but got \"+\"")}
	expect(t, expected, New("2 ++ 2").Tokens())
}

func TestMinus(t *testing.T) {
	expected := []Token{NumberToken.Of("2"), MinusToken.Of("-"), NumberToken.Of("2"), EOFToken.Of("")}
	expect(t, expected, New("2 - 2").Tokens())
	expect(t, expected, New("2-2").Tokens())
	expect(t, expected, New("   2 - 2").Tokens())
	expect(t, expected, New("2-2   ").Tokens())
}

func TestManyMinus(t *testing.T) {
	expected := []Token{NumberToken.Of("2"), MinusToken.Of("-"), ErrorToken.Of("expected number, but got \"-\"")}
	expect(t, expected, New("2 --  2").Tokens())
}

func TestNextToken(t *testing.T) {
	tok := New("2 + 2")
	assert.Equal(t, NumberToken.Of("2"), tok.NextToken())
	assert.Equal(t, PlusToken.Of("+"), tok.NextToken())
	assert.Equal(t, NumberToken.Of("2"), tok.NextToken())
	assert.Equal(t, EOFToken.Of(""), tok.NextToken())
}

func expect(t *testing.T, expected []Token, tokens chan Token) {
	// pul the actual tokens off the channel
	actuals := make([]Token, 0)
	for token := range tokens {
		actuals = append(actuals, token)
	}
	assert.ElementsMatch(t, expected, actuals)
}