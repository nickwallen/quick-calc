package toks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecimals(t *testing.T) {
	expected := []Token{Number.Token("22"), EOF.Token("")}
	inputs := []string{"22", "  22", "22    "}
	for _, input := range inputs {
		expect(t, expected, New(input).Tokens())
	}
}

func TestPositiveDecimals(t *testing.T) {
	expected := []Token{Number.Token("+22"), EOF.Token("")}
	inputs := []string{"+22", "  +22", "+22    ", "+  22"}
	for _, input := range inputs {
		expect(t, expected, New(input).Tokens())
	}
}

func TestNegativeDecimals(t *testing.T) {
	expected := []Token{Number.Token("-22"), EOF.Token("")}
	inputs := []string{"-22", "  -22", "-22    ", "  - 22"}
	for _, input := range inputs {
		expect(t, expected, New(input).Tokens())
	}
}

func TestBadDecimal(t *testing.T) {
	expected := []Token{Error.Token("expected number, but got '2A'")}
	inputs := []string{"2A", "   2A", "2A   "}
	for _, input := range inputs {
		expect(t, expected, New(input).Tokens())
	}
}

func TestHexaDecimal(t *testing.T) {
	expected := []Token{Number.Token("0xAF"), EOF.Token("")}
	inputs := []string{"0xAF", "   0xAF", "0xAF   "}
	for _, input := range inputs {
		expect(t, expected, New(input).Tokens())
	}
}

func TestBadHexaDecimal(t *testing.T) {
	expected := []Token{Error.Token("expected number, but got '0xG'")}
	inputs := []string{"0xG2", "   0xG2", "0xG2   "}
	for _, input := range inputs {
		expect(t, expected, New(input).Tokens())
	}
}

func TestFloats(t *testing.T) {
	expected := []Token{Number.Token("2.22"), EOF.Token("")}
	inputs := []string{"2.22", "   2.22", "2.22   "}
	for _, input := range inputs {
		expect(t, expected, New(input).Tokens())
	}
}

func TestExpNotation(t *testing.T) {
	expected := []Token{Number.Token("2E10"), EOF.Token("")}
	inputs := []string{"2E10", "   2E10", "2E10   "}
	for _, input := range inputs {
		expect(t, expected, New(input).Tokens())
	}
}

func TestPlus(t *testing.T) {
	expected := []Token{Number.Token("2"), Plus.Token("+"), Number.Token("2"), EOF.Token("")}
	inputs := []string{"2 + 2", "   2+2", "   2 +   2   ", "2+2"}
	for _, input := range inputs {
		expect(t, expected, New(input).Tokens())
	}
}

func TestPlusNegatives(t *testing.T) {
	expected := []Token{Number.Token("2"), Plus.Token("+"), Number.Token("-2"), EOF.Token("")}
	inputs := []string{"2 + -2", "   2+-2", "   2 +   -2   ", "2+-2"}
	for _, input := range inputs {
		expect(t, expected, New(input).Tokens())
	}
}

func TestPlusPositives(t *testing.T) {
	expected := []Token{Number.Token("2"), Plus.Token("+"), Number.Token("+2"), EOF.Token("")}
	inputs := []string{"2 + +2", "   2++2", "   2 +   +2   ", "2++2"}
	for _, input := range inputs {
		expect(t, expected, New(input).Tokens())
	}
}

func TestTooManyPlus(t *testing.T) {
	expected := []Token{Number.Token("2"), Plus.Token("+"), Error.Token("expected number, but got '++'")}
	inputs := []string{"2 +++ 2", "   2+++2", "   2+++   2   ", "2+++2"}
	for _, input := range inputs {
		expect(t, expected, New(input).Tokens())
	}
}

func TestMinus(t *testing.T) {
	expected := []Token{Number.Token("2"), Minus.Token("-"), Number.Token("2"), EOF.Token("")}
	inputs := []string{"2 - 2", "   2-2", "   2 -   2   ", "2-2"}
	for _, input := range inputs {
		expect(t, expected, New(input).Tokens())
	}
}

func TestMinusNegatives(t *testing.T) {
	expected := []Token{Number.Token("2"), Minus.Token("-"), Number.Token("-2"), EOF.Token("")}
	inputs := []string{"2 - -2", "   2--2", "   2 -   -2   ", "2--2"}
	for _, input := range inputs {
		expect(t, expected, New(input).Tokens())
	}
}

func TestMinusPositives(t *testing.T) {
	expected := []Token{Number.Token("2"), Minus.Token("-"), Number.Token("+2"), EOF.Token("")}
	inputs := []string{"2 - +2", "   2-+2", "   2 -   +2   ", "2-+2"}
	for _, input := range inputs {
		expect(t, expected, New(input).Tokens())
	}
}

func TestTooManyMinus(t *testing.T) {
	expected := []Token{Number.Token("2"), Minus.Token("-"), Error.Token("expected number, but got '--'")}
	inputs := []string{"2 --- 2", "   2---2", "   2 ---   2   ", "2---2"}
	for _, input := range inputs {
		expect(t, expected, New(input).Tokens())
	}
}

func TestMultiply(t *testing.T) {
	expected := []Token{Number.Token("2"), Multiply.Token("*"), Number.Token("2"), EOF.Token("")}
	inputs := []string{"2 * 2", "   2*2", "   2 *   2   ", "2*2"}
	for _, input := range inputs {
		expect(t, expected, New(input).Tokens())
	}
}

func TestManyMultiplies(t *testing.T) {
	expected := []Token{Number.Token("2"), Multiply.Token("*"), Error.Token("expected number, but got '*'")}
	inputs := []string{"2 ** 2", "   2**2", "   2 **   2   ", "2**2"}
	for _, input := range inputs {
		expect(t, expected, New(input).Tokens())
	}
}

func TestDivide(t *testing.T) {
	expected := []Token{Number.Token("2"), Divide.Token("/"), Number.Token("2"), EOF.Token("")}
	inputs := []string{"2 / 2", "   2/2", "   2 /   2   ", "2/2"}
	for _, input := range inputs {
		expect(t, expected, New(input).Tokens())
	}
}

func TestManyDivides(t *testing.T) {
	expected := []Token{Number.Token("2"), Divide.Token("/"), Error.Token("expected number, but got '/'")}
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
