package parser

import (
	"fmt"
	"github.com/nickwallen/toks/internal/tokenizer"
	"github.com/stretchr/testify/assert"
	"testing"
)

type tokenSlice struct {
	slice    []tokenizer.Token
	position int
}

// allows the tests to read tokens from a slice
func (t *tokenSlice) ReadToken() (tokenizer.Token, error) {
	var token tokenizer.Token
	if t.position >= len(t.slice) {
		return token, fmt.Errorf("no tokens left")
	}
	token = t.slice[t.position]
	t.position++
	return token, nil
}

func tokens(input ...tokenizer.Token) tokenSlice {
	return tokenSlice{slice: input, position: 0}
}

func TestParseAmount(t *testing.T) {
	tokens := tokens(tokenizer.Number.Token("23"), tokenizer.Units.Token("pounds"), tokenizer.EOF.Token(""))
	actual, err := Parse(&tokens)
	expected := AmountOf(23, pounds())
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestParseAmountNoUnits(t *testing.T) {
	tokens := tokens(tokenizer.Number.Token("23"), tokenizer.EOF.Token(""))
	_, err := Parse(&tokens)
	assert.Equal(t, "expected units, but reached end of input", err.Error())
}

func TestParseAmountNoNumber(t *testing.T) {
	tokens := tokens(tokenizer.Units.Token("pounds"), tokenizer.EOF.Token(""))
	_, err := Parse(&tokens)
	assert.Equal(t, "expected number, but got 'pounds'", err.Error())
}

func TestParseSum(t *testing.T) {
	tokens := tokens(
		tokenizer.Number.Token("23"),
		tokenizer.Units.Token("kg"),
		tokenizer.Plus.Token("+"),
		tokenizer.Number.Token("23"),
		tokenizer.Units.Token("pounds"),
		tokenizer.EOF.Token(""))
	actual, err := Parse(&tokens)
	expected := operatorOf(AmountOf(23, kilos()), AmountOf(23, pounds()), kilos(), tokenizer.Plus)
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestParseSubtract(t *testing.T) {
	tokens := tokens(
		tokenizer.Number.Token("23"),
		tokenizer.Units.Token("kg"),
		tokenizer.Minus.Token("-"),
		tokenizer.Number.Token("23"),
		tokenizer.Units.Token("pounds"),
		tokenizer.EOF.Token(""))
	actual, err := Parse(&tokens)
	expected := operatorOf(AmountOf(23, kilos()), AmountOf(23, pounds()), kilos(), tokenizer.Minus)
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestParseConversion(t *testing.T) {
	tokens := tokens(
		tokenizer.Number.Token("2"),
		tokenizer.Units.Token("pounds"),
		tokenizer.In.Token("in"),
		tokenizer.Units.Token("ounces"),
		tokenizer.EOF.Token(""))
	actual, err := Parse(&tokens)
	expected := unitConverterOf(AmountOf(2, pounds()), ounces())
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestParseSumAndConvert(t *testing.T) {
	tokens := tokens(
		tokenizer.Number.Token("2"),
		tokenizer.Units.Token("ounces"),
		tokenizer.Plus.Token("+"),
		tokenizer.Number.Token("2"),
		tokenizer.Units.Token("pounds"),
		tokenizer.In.Token("in"),
		tokenizer.Units.Token("pounds"),
		tokenizer.EOF.Token(""))
	actual, err := Parse(&tokens)
	expected := operatorOf(AmountOf(2, ounces()), AmountOf(2, pounds()), pounds(), tokenizer.Plus)
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestParseSubtractAndConvert(t *testing.T) {
	tokens := tokens(
		tokenizer.Number.Token("2"),
		tokenizer.Units.Token("pounds"),
		tokenizer.Minus.Token("-"),
		tokenizer.Number.Token("2"),
		tokenizer.Units.Token("ounces"),
		tokenizer.In.Token("in"),
		tokenizer.Units.Token("ounces"),
		tokenizer.EOF.Token(""))
	actual, err := Parse(&tokens)
	expected := operatorOf(AmountOf(2, pounds()), AmountOf(2, ounces()), ounces(), tokenizer.Minus)
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}
