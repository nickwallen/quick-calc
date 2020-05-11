package tokenizer

import (
	"github.com/nickwallen/quick-calc/internal/io"
	"github.com/nickwallen/quick-calc/internal/tokens"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testCases = map[string][]tokens.Token{
	"": {
		tokens.Error.Token("expected number, but got ''"),
	},
	" ": {
		tokens.Error.Token("expected number, but got ''"),
	},
	"22": {
		tokens.Number.Token("22"),
		tokens.EOF.Token(""),
	},
	"  22": {
		tokens.Number.Token("22"),
		tokens.EOF.Token(""),
	},
	"22    ": {
		tokens.Number.Token("22"),
		tokens.EOF.Token(""),
	},
	"0": {
		tokens.Number.Token("0"),
		tokens.EOF.Token(""),
	},
	"  0": {
		tokens.Number.Token("0"),
		tokens.EOF.Token(""),
	},
	"0    ": {
		tokens.Number.Token("0"),
		tokens.EOF.Token(""),
	},
	"2,200,123": {
		tokens.Number.Token("2,200,123"),
		tokens.EOF.Token(""),
	},
	"  2,200,123": {
		tokens.Number.Token("2,200,123"),
		tokens.EOF.Token(""),
	},
	"2,200,123    ": {
		tokens.Number.Token("2,200,123"),
		tokens.EOF.Token(""),
	},
	",200,200": {
		tokens.Error.Token("expected number, but got ',2'"),
	},
	"+22": {
		tokens.Number.Token("+22"),
		tokens.EOF.Token(""),
	},
	"  +22": {
		tokens.Number.Token("+22"),
		tokens.EOF.Token(""),
	},
	"+22    ": {
		tokens.Number.Token("+22"),
		tokens.EOF.Token(""),
	},
	"+  22": {
		tokens.Number.Token("+22"),
		tokens.EOF.Token(""),
	},
	"-22": {
		tokens.Number.Token("-22"),
		tokens.EOF.Token(""),
	},
	"  -22": {
		tokens.Number.Token("-22"),
		tokens.EOF.Token(""),
	},
	"-22    ": {
		tokens.Number.Token("-22"),
		tokens.EOF.Token(""),
	},
	"  - 22": {
		tokens.Number.Token("-22"),
		tokens.EOF.Token(""),
	},
	"2?": {
		tokens.Number.Token("2"),
		tokens.Error.Token("expected symbol, but got '?'"),
	},
	"   2?": {
		tokens.Number.Token("2"),
		tokens.Error.Token("expected symbol, but got '?'"),
	},
	"2?   ": {
		tokens.Number.Token("2"),
		tokens.Error.Token("expected symbol, but got '?'"),
	},
	"0xAF": {
		tokens.Number.Token("0xAF"),
		tokens.EOF.Token(""),
	},
	"   0xAF": {
		tokens.Number.Token("0xAF"),
		tokens.EOF.Token(""),
	},
	"0xAF   ": {
		tokens.Number.Token("0xAF"),
		tokens.EOF.Token(""),
	},
	"0xG2": {
		tokens.Error.Token("expected number, but got '0xG'"),
	},
	"   0xG2": {
		tokens.Error.Token("expected number, but got '0xG'"),
	},
	"0xG2   ": {
		tokens.Error.Token("expected number, but got '0xG'"),
	},
	"0x": {
		tokens.Error.Token("expected number, but got '0x'"),
	},
	"   0x": {
		tokens.Error.Token("expected number, but got '0x'"),
	},
	"0x   ": {
		tokens.Error.Token("expected number, but got '0x'"),
	},
	"2.22": {
		tokens.Number.Token("2.22"),
		tokens.EOF.Token(""),
	},
	"   2.22": {
		tokens.Number.Token("2.22"),
		tokens.EOF.Token(""),
	},
	"2.22   ": {
		tokens.Number.Token("2.22"),
		tokens.EOF.Token(""),
	},
	"2E10": {
		tokens.Number.Token("2E10"),
		tokens.EOF.Token(""),
	},
	"   2E10": {
		tokens.Number.Token("2E10"),
		tokens.EOF.Token(""),
	},
	"2E10   ": {
		tokens.Number.Token("2E10"),
		tokens.EOF.Token(""),
	},
	"2 + 2": {
		tokens.Number.Token("2"),
		tokens.Plus.Token("+"),
		tokens.Number.Token("2"),
		tokens.EOF.Token(""),
	},
	"   2+2": {
		tokens.Number.Token("2"),
		tokens.Plus.Token("+"),
		tokens.Number.Token("2"),
		tokens.EOF.Token(""),
	},
	"   2 +   2   ": {
		tokens.Number.Token("2"),
		tokens.Plus.Token("+"),
		tokens.Number.Token("2"),
		tokens.EOF.Token(""),
	},
	"2+2": {
		tokens.Number.Token("2"),
		tokens.Plus.Token("+"),
		tokens.Number.Token("2"),
		tokens.EOF.Token(""),
	},
	"2 + -2": {
		tokens.Number.Token("2"),
		tokens.Plus.Token("+"),
		tokens.Number.Token("-2"),
		tokens.EOF.Token(""),
	},
	"   2+-2": {
		tokens.Number.Token("2"),
		tokens.Plus.Token("+"),
		tokens.Number.Token("-2"),
		tokens.EOF.Token(""),
	},
	"   2 +   -2   ": {
		tokens.Number.Token("2"),
		tokens.Plus.Token("+"),
		tokens.Number.Token("-2"),
		tokens.EOF.Token(""),
	},
	"2+-2": {
		tokens.Number.Token("2"),
		tokens.Plus.Token("+"),
		tokens.Number.Token("-2"),
		tokens.EOF.Token(""),
	},
	"2 + +2": {
		tokens.Number.Token("2"),
		tokens.Plus.Token("+"),
		tokens.Number.Token("+2"),
		tokens.EOF.Token(""),
	},
	"   2++2": {
		tokens.Number.Token("2"),
		tokens.Plus.Token("+"),
		tokens.Number.Token("+2"),
		tokens.EOF.Token(""),
	},
	"   2 +   +2   ": {
		tokens.Number.Token("2"),
		tokens.Plus.Token("+"),
		tokens.Number.Token("+2"),
		tokens.EOF.Token(""),
	},
	"2++2": {
		tokens.Number.Token("2"),
		tokens.Plus.Token("+"),
		tokens.Number.Token("+2"),
		tokens.EOF.Token(""),
	},
	"2 +++ 2": {
		tokens.Number.Token("2"),
		tokens.Plus.Token("+"),
		tokens.Error.Token("expected number, but got '++'"),
	},
	"   2+++2": {
		tokens.Number.Token("2"),
		tokens.Plus.Token("+"),
		tokens.Error.Token("expected number, but got '++'"),
	},
	"   2+++   2   ": {
		tokens.Number.Token("2"),
		tokens.Plus.Token("+"),
		tokens.Error.Token("expected number, but got '++'"),
	},
	"2+++2": {
		tokens.Number.Token("2"),
		tokens.Plus.Token("+"),
		tokens.Error.Token("expected number, but got '++'"),
	},
	"2 - 2": {
		tokens.Number.Token("2"),
		tokens.Minus.Token("-"),
		tokens.Number.Token("2"),
		tokens.EOF.Token(""),
	},
	"   2-2": {
		tokens.Number.Token("2"),
		tokens.Minus.Token("-"),
		tokens.Number.Token("2"),
		tokens.EOF.Token(""),
	},
	"   2 -   2   ": {
		tokens.Number.Token("2"),
		tokens.Minus.Token("-"),
		tokens.Number.Token("2"),
		tokens.EOF.Token(""),
	},
	"2-2": {
		tokens.Number.Token("2"),
		tokens.Minus.Token("-"),
		tokens.Number.Token("2"),
		tokens.EOF.Token(""),
	},
	"2 - -2": {
		tokens.Number.Token("2"),
		tokens.Minus.Token("-"),
		tokens.Number.Token("-2"),
		tokens.EOF.Token(""),
	},
	"   2--2": {
		tokens.Number.Token("2"),
		tokens.Minus.Token("-"),
		tokens.Number.Token("-2"),
		tokens.EOF.Token(""),
	},
	"   2 -   -2   ": {
		tokens.Number.Token("2"),
		tokens.Minus.Token("-"),
		tokens.Number.Token("-2"),
		tokens.EOF.Token(""),
	},
	"2--2": {
		tokens.Number.Token("2"),
		tokens.Minus.Token("-"),
		tokens.Number.Token("-2"),
		tokens.EOF.Token(""),
	},
	"2 - +2": {
		tokens.Number.Token("2"),
		tokens.Minus.Token("-"),
		tokens.Number.Token("+2"),
		tokens.EOF.Token(""),
	},
	"   2-+2": {
		tokens.Number.Token("2"),
		tokens.Minus.Token("-"),
		tokens.Number.Token("+2"),
		tokens.EOF.Token(""),
	},
	"   2 -   +2   ": {
		tokens.Number.Token("2"),
		tokens.Minus.Token("-"),
		tokens.Number.Token("+2"),
		tokens.EOF.Token(""),
	},
	"2-+2": {
		tokens.Number.Token("2"),
		tokens.Minus.Token("-"),
		tokens.Number.Token("+2"),
		tokens.EOF.Token(""),
	},
	"2 --- 2": {
		tokens.Number.Token("2"),
		tokens.Minus.Token("-"),
		tokens.Error.Token("expected number, but got '--'"),
	},
	"   2---2": {
		tokens.Number.Token("2"),
		tokens.Minus.Token("-"),
		tokens.Error.Token("expected number, but got '--'"),
	},
	"   2 ---   2   ": {
		tokens.Number.Token("2"),
		tokens.Minus.Token("-"),
		tokens.Error.Token("expected number, but got '--'"),
	},
	"2---2": {
		tokens.Number.Token("2"),
		tokens.Minus.Token("-"),
		tokens.Error.Token("expected number, but got '--'"),
	},
	"2 * 2": {
		tokens.Number.Token("2"),
		tokens.Multiply.Token("*"),
		tokens.Number.Token("2"),
		tokens.EOF.Token(""),
	},
	"   2*2": {
		tokens.Number.Token("2"),
		tokens.Multiply.Token("*"),
		tokens.Number.Token("2"),
		tokens.EOF.Token(""),
	},
	"   2 *   2   ": {
		tokens.Number.Token("2"),
		tokens.Multiply.Token("*"),
		tokens.Number.Token("2"),
		tokens.EOF.Token(""),
	},
	"2*2": {
		tokens.Number.Token("2"),
		tokens.Multiply.Token("*"),
		tokens.Number.Token("2"),
		tokens.EOF.Token(""),
	},
	"2 ** 2": {
		tokens.Number.Token("2"),
		tokens.Multiply.Token("*"),
		tokens.Error.Token("expected number, but got '*'"),
	},
	"   2**2": {
		tokens.Number.Token("2"),
		tokens.Multiply.Token("*"),
		tokens.Error.Token("expected number, but got '*'"),
	},
	"   2 **   2   ": {
		tokens.Number.Token("2"),
		tokens.Multiply.Token("*"),
		tokens.Error.Token("expected number, but got '*'"),
	},
	"2**2": {
		tokens.Number.Token("2"),
		tokens.Multiply.Token("*"),
		tokens.Error.Token("expected number, but got '*'"),
	},
	"2 / 2": {
		tokens.Number.Token("2"),
		tokens.Divide.Token("/"),
		tokens.Number.Token("2"),
		tokens.EOF.Token(""),
	},
	"   2/2": {
		tokens.Number.Token("2"),
		tokens.Divide.Token("/"),
		tokens.Number.Token("2"),
		tokens.EOF.Token(""),
	},
	"   2 /   2   ": {
		tokens.Number.Token("2"),
		tokens.Divide.Token("/"),
		tokens.Number.Token("2"),
		tokens.EOF.Token(""),
	},
	"2/2": {
		tokens.Number.Token("2"),
		tokens.Divide.Token("/"),
		tokens.Number.Token("2"),
		tokens.EOF.Token(""),
	},
	"2 // 2": {
		tokens.Number.Token("2"),
		tokens.Divide.Token("/"),
		tokens.Error.Token("expected number, but got '/'"),
	},
	"   2//2": {
		tokens.Number.Token("2"),
		tokens.Divide.Token("/"),
		tokens.Error.Token("expected number, but got '/'"),
	},
	"   2 //   2   ": {
		tokens.Number.Token("2"),
		tokens.Divide.Token("/"),
		tokens.Error.Token("expected number, but got '/'"),
	},
	"2//2": {
		tokens.Number.Token("2"),
		tokens.Divide.Token("/"),
		tokens.Error.Token("expected number, but got '/'"),
	},
	"245 lbs": {
		tokens.Number.Token("245"),
		tokens.Units.Token("lbs"),
		tokens.EOF.Token(""),
	},
	"    245 lbs": {
		tokens.Number.Token("245"),
		tokens.Units.Token("lbs"),
		tokens.EOF.Token(""),
	},
	"245lbs": {
		tokens.Number.Token("245"),
		tokens.Units.Token("lbs"),
		tokens.EOF.Token(""),
	},
	"245 lbs + 37.50kg": {
		tokens.Number.Token("245"),
		tokens.Units.Token("lbs"),
		tokens.Plus.Token("+"),
		tokens.Number.Token("37.50"),
		tokens.Units.Token("kg"),
		tokens.EOF.Token(""),
	},
	"245   lbs   + 37.50   kg": {
		tokens.Number.Token("245"),
		tokens.Units.Token("lbs"),
		tokens.Plus.Token("+"),
		tokens.Number.Token("37.50"),
		tokens.Units.Token("kg"),
		tokens.EOF.Token(""),
	},
	"20 lbs in kg": {
		tokens.Number.Token("20"),
		tokens.Units.Token("lbs"),
		tokens.In.Token("in"),
		tokens.Units.Token("kg"),
		tokens.EOF.Token(""),
	},
	"   20lbs in   kg   ": {
		tokens.Number.Token("20"),
		tokens.Units.Token("lbs"),
		tokens.In.Token("in"),
		tokens.Units.Token("kg"),
		tokens.EOF.Token(""),
	},
	"20 ints": {
		tokens.Number.Token("20"),
		tokens.Units.Token("ints"),
		tokens.EOF.Token(""),
	},
	"   20ints   ": {
		tokens.Number.Token("20"),
		tokens.Units.Token("ints"),
		tokens.EOF.Token(""),
	},
	"245 lbs + 37.50 kg in kg": {
		tokens.Number.Token("245"),
		tokens.Units.Token("lbs"),
		tokens.Plus.Token("+"),
		tokens.Number.Token("37.50"),
		tokens.Units.Token("kg"),
		tokens.In.Token("in"),
		tokens.Units.Token("kg"),
		tokens.EOF.Token(""),
	},
}

func TestTokens(t *testing.T) {
	for input, expected := range testCases {
		output := io.NewTokenChannel()
		go Tokenize(input, &output)
		for _, expect := range expected {
			actual, err := output.ReadToken()
			assert.Nil(t, err)
			assert.Equal(t, expect, actual)
		}
	}
}
