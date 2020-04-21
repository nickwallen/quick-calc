package toks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmpty(t *testing.T) {
	expected := []Token{Error.token("expected number, but got ''")}
	inputs := []string{"", "  "}
	for _, input := range inputs {
		expect(t, expected, NewTokenizer(input).Tokens())
	}
}

func TestDecimals(t *testing.T) {
	expected := []Token{Number.token("22"), EOF.token("")}
	inputs := []string{"22", "  22", "22    "}
	for _, input := range inputs {
		expect(t, expected, NewTokenizer(input).Tokens())
	}
}

func TestZeros(t *testing.T) {
	expected := []Token{Number.token("0"), EOF.token("")}
	inputs := []string{"0", "  0", "0    "}
	for _, input := range inputs {
		expect(t, expected, NewTokenizer(input).Tokens())
	}
}

func TestCommas(t *testing.T) {
	expected := []Token{Number.token("2,200,123"), EOF.token("")}
	inputs := []string{"2,200,123", "  2,200,123", "2,200,123    "}
	for _, input := range inputs {
		expect(t, expected, NewTokenizer(input).Tokens())
	}
}

func TestLeadComma(t *testing.T) {
	input := ",200,200"
	expect(t, []Token{Error.token("expected number, but got ',2'")}, NewTokenizer(input).Tokens())
}

func TestPositiveDecimals(t *testing.T) {
	expected := []Token{Number.token("+22"), EOF.token("")}
	inputs := []string{"+22", "  +22", "+22    ", "+  22"}
	for _, input := range inputs {
		expect(t, expected, NewTokenizer(input).Tokens())
	}
}

func TestNegativeDecimals(t *testing.T) {
	expected := []Token{Number.token("-22"), EOF.token("")}
	inputs := []string{"-22", "  -22", "-22    ", "  - 22"}
	for _, input := range inputs {
		expect(t, expected, NewTokenizer(input).Tokens())
	}
}

func TestBadDecimal(t *testing.T) {
	expected := []Token{Number.token("2"), Error.token("expected symbol, but got '?'")}
	inputs := []string{"2?", "   2?", "2?   "}
	for _, input := range inputs {
		expect(t, expected, NewTokenizer(input).Tokens())
	}
}

func TestHexaDecimal(t *testing.T) {
	expected := []Token{Number.token("0xAF"), EOF.token("")}
	inputs := []string{"0xAF", "   0xAF", "0xAF   "}
	for _, input := range inputs {
		expect(t, expected, NewTokenizer(input).Tokens())
	}
}

func TestBadHexaDecimal(t *testing.T) {
	expected := []Token{Error.token("expected number, but got '0xG'")}
	inputs := []string{"0xG2", "   0xG2", "0xG2   "}
	for _, input := range inputs {
		expect(t, expected, NewTokenizer(input).Tokens())
	}
}

func TestNoHexaDecimal(t *testing.T) {
	expected := []Token{Error.token("expected number, but got '0x'")}
	inputs := []string{"0x", "   0x", "0x   "}
	for _, input := range inputs {
		expect(t, expected, NewTokenizer(input).Tokens())
	}
}

func TestFloats(t *testing.T) {
	expected := []Token{Number.token("2.22"), EOF.token("")}
	inputs := []string{"2.22", "   2.22", "2.22   "}
	for _, input := range inputs {
		expect(t, expected, NewTokenizer(input).Tokens())
	}
}

func TestExpNotation(t *testing.T) {
	expected := []Token{Number.token("2E10"), EOF.token("")}
	inputs := []string{"2E10", "   2E10", "2E10   "}
	for _, input := range inputs {
		expect(t, expected, NewTokenizer(input).Tokens())
	}
}

func TestPlus(t *testing.T) {
	expected := []Token{Number.token("2"), Plus.token("+"), Number.token("2"), EOF.token("")}
	inputs := []string{"2 + 2", "   2+2", "   2 +   2   ", "2+2"}
	for _, input := range inputs {
		expect(t, expected, NewTokenizer(input).Tokens())
	}
}

func TestPlusNegatives(t *testing.T) {
	expected := []Token{Number.token("2"), Plus.token("+"), Number.token("-2"), EOF.token("")}
	inputs := []string{"2 + -2", "   2+-2", "   2 +   -2   ", "2+-2"}
	for _, input := range inputs {
		expect(t, expected, NewTokenizer(input).Tokens())
	}
}

func TestPlusPositives(t *testing.T) {
	expected := []Token{Number.token("2"), Plus.token("+"), Number.token("+2"), EOF.token("")}
	inputs := []string{"2 + +2", "   2++2", "   2 +   +2   ", "2++2"}
	for _, input := range inputs {
		expect(t, expected, NewTokenizer(input).Tokens())
	}
}

func TestTooManyPlus(t *testing.T) {
	expected := []Token{Number.token("2"), Plus.token("+"), Error.token("expected number, but got '++'")}
	inputs := []string{"2 +++ 2", "   2+++2", "   2+++   2   ", "2+++2"}
	for _, input := range inputs {
		expect(t, expected, NewTokenizer(input).Tokens())
	}
}

func TestMinus(t *testing.T) {
	expected := []Token{Number.token("2"), Minus.token("-"), Number.token("2"), EOF.token("")}
	inputs := []string{"2 - 2", "   2-2", "   2 -   2   ", "2-2"}
	for _, input := range inputs {
		expect(t, expected, NewTokenizer(input).Tokens())
	}
}

func TestMinusNegatives(t *testing.T) {
	expected := []Token{Number.token("2"), Minus.token("-"), Number.token("-2"), EOF.token("")}
	inputs := []string{"2 - -2", "   2--2", "   2 -   -2   ", "2--2"}
	for _, input := range inputs {
		expect(t, expected, NewTokenizer(input).Tokens())
	}
}

func TestMinusPositives(t *testing.T) {
	expected := []Token{Number.token("2"), Minus.token("-"), Number.token("+2"), EOF.token("")}
	inputs := []string{"2 - +2", "   2-+2", "   2 -   +2   ", "2-+2"}
	for _, input := range inputs {
		expect(t, expected, NewTokenizer(input).Tokens())
	}
}

func TestTooManyMinus(t *testing.T) {
	expected := []Token{Number.token("2"), Minus.token("-"), Error.token("expected number, but got '--'")}
	inputs := []string{"2 --- 2", "   2---2", "   2 ---   2   ", "2---2"}
	for _, input := range inputs {
		expect(t, expected, NewTokenizer(input).Tokens())
	}
}

func TestMultiply(t *testing.T) {
	expected := []Token{Number.token("2"), Multiply.token("*"), Number.token("2"), EOF.token("")}
	inputs := []string{"2 * 2", "   2*2", "   2 *   2   ", "2*2"}
	for _, input := range inputs {
		expect(t, expected, NewTokenizer(input).Tokens())
	}
}

func TestManyMultiplies(t *testing.T) {
	expected := []Token{Number.token("2"), Multiply.token("*"), Error.token("expected number, but got '*'")}
	inputs := []string{"2 ** 2", "   2**2", "   2 **   2   ", "2**2"}
	for _, input := range inputs {
		expect(t, expected, NewTokenizer(input).Tokens())
	}
}

func TestDivide(t *testing.T) {
	expected := []Token{Number.token("2"), Divide.token("/"), Number.token("2"), EOF.token("")}
	inputs := []string{"2 / 2", "   2/2", "   2 /   2   ", "2/2"}
	for _, input := range inputs {
		expect(t, expected, NewTokenizer(input).Tokens())
	}
}

func TestManyDivides(t *testing.T) {
	expected := []Token{Number.token("2"), Divide.token("/"), Error.token("expected number, but got '/'")}
	inputs := []string{"2 // 2", "   2//2", "   2 //   2   ", "2//2"}
	for _, input := range inputs {
		expect(t, expected, NewTokenizer(input).Tokens())
	}
}

func TestUnits(t *testing.T) {
	expected := []Token{Number.token("245"), Units.token("pounds"), EOF.token("")}
	inputs := []string{"245 pounds", "    245 pounds", "245pounds"}
	for _, input := range inputs {
		expect(t, expected, NewTokenizer(input).Tokens())
	}
}

func TestAddUnits(t *testing.T) {
	expected := []Token{
		Number.token("245"),
		Units.token("pounds"),
		Plus.token("+"),
		Number.token("37.50"),
		Units.token("kg"),
		EOF.token("")}
	inputs := []string{"245 pounds + 37.50kg", "245   pounds   + 37.50   kg"}
	for _, input := range inputs {
		expect(t, expected, NewTokenizer(input).Tokens())
	}
}

func TestConversion(t *testing.T) {
	expected := []Token{Number.token("20"), Units.token("lbs"), In.token("in"),
		Units.token("kg"), EOF.token("")}
	inputs := []string{"20 lbs in kg", "   20lbs in   kg   "}
	for _, input := range inputs {
		expect(t, expected, NewTokenizer(input).Tokens())
	}
}

func TestUnitsStartWithIn(t *testing.T) {
	expected := []Token{Number.token("20"), Units.token("ints"), EOF.token("")}
	inputs := []string{"20 ints", "   20ints   "}
	for _, input := range inputs {
		expect(t, expected, NewTokenizer(input).Tokens())
	}
}

func TestAddThenConvert(t *testing.T) {
	expected := []Token{
		Number.token("245"),
		Units.token("pounds"),
		Plus.token("+"),
		Number.token("37.50"),
		Units.token("kg"),
		In.token("in"),
		Units.token("kilos"),
		EOF.token("")}
	inputs := []string{"245 pounds + 37.50 kg in kilos"}
	for _, input := range inputs {
		expect(t, expected, NewTokenizer(input).Tokens())
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
