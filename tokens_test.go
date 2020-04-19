package toks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecimals(t *testing.T) {
	expected := []Token{NumberToken.Of("22"), EOFToken.Of("")}
	inputs := []string{"22", "  22", "22    "}
	for _, input := range inputs {
		expect(t, expected, New(input).Tokens())
	}
}

func TestPositiveDecimals(t *testing.T) {
	expected := []Token{NumberToken.Of("+22"), EOFToken.Of("")}
	inputs := []string{"+22", "  +22", "+22    ", "+  22"}
	for _, input := range inputs {
		expect(t, expected, New(input).Tokens())
	}
}

func TestNegativeDecimals(t *testing.T) {
	expected := []Token{NumberToken.Of("-22"), EOFToken.Of("")}
	inputs := []string{"-22", "  -22", "-22    ", "  - 22"}
	for _, input := range inputs {
		expect(t, expected, New(input).Tokens())
	}
}

func TestBadDecimal(t *testing.T) {
	expected := []Token{ErrorToken.Of("expected number, but got '2A'")}
	inputs := []string{"2A", "   2A", "2A   "}
	for _, input := range inputs {
		expect(t, expected, New(input).Tokens())
	}
}

func TestHexaDecimal(t *testing.T) {
	expected := []Token{NumberToken.Of("0xAF"), EOFToken.Of("")}
	inputs := []string{"0xAF", "   0xAF", "0xAF   "}
	for _, input := range inputs {
		expect(t, expected, New(input).Tokens())
	}
}

func TestBadHexaDecimal(t *testing.T) {
	expected := []Token{ErrorToken.Of("expected number, but got '0xG'")}
	inputs := []string{"0xG2", "   0xG2", "0xG2   "}
	for _, input := range inputs {
		expect(t, expected, New(input).Tokens())
	}
}

func TestFloats(t *testing.T) {
	expected := []Token{NumberToken.Of("2.22"), EOFToken.Of("")}
	inputs := []string{"2.22", "   2.22", "2.22   "}
	for _, input := range inputs {
		expect(t, expected, New(input).Tokens())
	}
}

func TestExpNotation(t *testing.T) {
	expected := []Token{NumberToken.Of("2E10"), EOFToken.Of("")}
	inputs := []string{"2E10", "   2E10", "2E10   "}
	for _, input := range inputs {
		expect(t, expected, New(input).Tokens())
	}
}

func TestPlus(t *testing.T) {
	expected := []Token{NumberToken.Of("2"), PlusToken.Of("+"), NumberToken.Of("2"), EOFToken.Of("")}
	inputs := []string{"2 + 2", "   2+2", "   2 +   2   ", "2+2"}
	for _, input := range inputs {
		expect(t, expected, New(input).Tokens())
	}
}

func TestPlusNegatives(t *testing.T) {
	expected := []Token{NumberToken.Of("2"), PlusToken.Of("+"), NumberToken.Of("-2"), EOFToken.Of("")}
	inputs := []string{"2 + -2", "   2+-2", "   2 +   -2   ", "2+-2"}
	for _, input := range inputs {
		expect(t, expected, New(input).Tokens())
	}
}

func TestPlusPositives(t *testing.T) {
	expected := []Token{NumberToken.Of("2"), PlusToken.Of("+"), NumberToken.Of("+2"), EOFToken.Of("")}
	inputs := []string{"2 + +2", "   2++2", "   2 +   +2   ", "2++2"}
	for _, input := range inputs {
		expect(t, expected, New(input).Tokens())
	}
}

func TestTooManyPlus(t *testing.T) {
	expected := []Token{NumberToken.Of("2"), PlusToken.Of("+"), ErrorToken.Of("expected number, but got '++'")}
	inputs := []string{"2 +++ 2", "   2+++2", "   2+++   2   ", "2+++2"}
	for _, input := range inputs {
		expect(t, expected, New(input).Tokens())
	}
}

func TestMinus(t *testing.T) {
	expected := []Token{NumberToken.Of("2"), MinusToken.Of("-"), NumberToken.Of("2"), EOFToken.Of("")}
	inputs := []string{"2 - 2", "   2-2", "   2 -   2   ", "2-2"}
	for _, input := range inputs {
		expect(t, expected, New(input).Tokens())
	}
}

func TestMinusNegatives(t *testing.T) {
	expected := []Token{NumberToken.Of("2"), MinusToken.Of("-"), NumberToken.Of("-2"), EOFToken.Of("")}
	inputs := []string{"2 - -2", "   2--2", "   2 -   -2   ", "2--2"}
	for _, input := range inputs {
		expect(t, expected, New(input).Tokens())
	}
}

func TestMinusPositives(t *testing.T) {
	expected := []Token{NumberToken.Of("2"), MinusToken.Of("-"), NumberToken.Of("+2"), EOFToken.Of("")}
	inputs := []string{"2 - +2", "   2-+2", "   2 -   +2   ", "2-+2"}
	for _, input := range inputs {
		expect(t, expected, New(input).Tokens())
	}
}

func TestTooManyMinus(t *testing.T) {
	expected := []Token{NumberToken.Of("2"), MinusToken.Of("-"), ErrorToken.Of("expected number, but got '--'")}
	inputs := []string{"2 --- 2", "   2---2", "   2 ---   2   ", "2---2"}
	for _, input := range inputs {
		expect(t, expected, New(input).Tokens())
	}
}

func TestMultiply(t *testing.T) {
	expected := []Token{NumberToken.Of("2"), MultiplyToken.Of("*"), NumberToken.Of("2"), EOFToken.Of("")}
	inputs := []string{"2 * 2", "   2*2", "   2 *   2   ", "2*2"}
	for _, input := range inputs {
		expect(t, expected, New(input).Tokens())
	}
}

func TestManyMultiplies(t *testing.T) {
	expected := []Token{NumberToken.Of("2"), MultiplyToken.Of("*"), ErrorToken.Of("expected number, but got '*'")}
	inputs := []string{"2 ** 2", "   2**2", "   2 **   2   ", "2**2"}
	for _, input := range inputs {
		expect(t, expected, New(input).Tokens())
	}
}

func TestDivide(t *testing.T) {
	expected := []Token{NumberToken.Of("2"), DivisionToken.Of("/"), NumberToken.Of("2"), EOFToken.Of("")}
	inputs := []string{"2 / 2", "   2/2", "   2 /   2   ", "2/2"}
	for _, input := range inputs {
		expect(t, expected, New(input).Tokens())
	}
}

func TestManyDivides(t *testing.T) {
	expected := []Token{NumberToken.Of("2"), DivisionToken.Of("/"), ErrorToken.Of("expected number, but got '/'")}
	inputs := []string{"2 // 2", "   2//2", "   2 //   2   ", "2//2"}
	for _, input := range inputs {
		expect(t, expected, New(input).Tokens())
	}
}

func expect(t *testing.T, expected []Token, tokens chan Token) {
	// pul the actual tokens off the channel
	actuals := make([]Token, 0)
	for token := range tokens {
		actuals = append(actuals, token)
	}
	assert.ElementsMatch(t, expected, actuals)
}
