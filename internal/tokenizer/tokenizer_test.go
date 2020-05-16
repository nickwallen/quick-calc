package tokenizer

import (
	"github.com/nickwallen/quick-calc/internal/io"
	"github.com/nickwallen/quick-calc/internal/tokens"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testCases = map[string][]tokens.Token{
	"": {
		tokens.Error.TokenAt("expected number, but got ''", 1),
	},
	" ": {
		tokens.Error.TokenAt("expected number, but got ''", 2),
	},
	"  22": {
		tokens.Number.TokenAt("22", 3),
		tokens.EOF.TokenAt("", 5),
	},
	"0": {
		tokens.Number.TokenAt("0", 1),
		tokens.EOF.TokenAt("", 2),
	},
	"2,200,123": {
		tokens.Number.TokenAt("2,200,123", 1),
		tokens.EOF.TokenAt("", 10),
	},
	",200,200": {
		tokens.Error.TokenAt("expected number, but got ',2'", 1),
	},
	"  +22": {
		tokens.Number.TokenAt("+22", 3),
		tokens.EOF.TokenAt("", 6),
	},
	"  -22": {
		tokens.Number.TokenAt("-22", 3),
		tokens.EOF.TokenAt("", 6),
	},
	"2?": {
		tokens.Number.TokenAt("2", 1),
		tokens.Error.TokenAt("expected symbol, but got '?'", 2),
	},
	"0xAF   ": {
		tokens.Number.TokenAt("0xAF", 1),
		tokens.EOF.TokenAt("", 8),
	},
	"0xG2": {
		tokens.Error.TokenAt("expected number, but got '0xG'", 1),
	},
	"0x   ": {
		tokens.Error.TokenAt("expected number, but got '0x'", 1),
	},
	"   2.22": {
		tokens.Number.TokenAt("2.22", 4),
		tokens.EOF.TokenAt("", 8),
	},
	"2E10   ": {
		tokens.Number.TokenAt("2E10", 1),
		tokens.EOF.TokenAt("", 8),
	},
	"2+2": {
		tokens.Number.TokenAt("2", 1),
		tokens.Plus.TokenAt("+", 2),
		tokens.Number.TokenAt("2", 3),
		tokens.EOF.TokenAt("", 4),
	},
	"2 + -2": {
		tokens.Number.TokenAt("2", 1),
		tokens.Plus.TokenAt("+", 3),
		tokens.Number.TokenAt("-2", 5),
		tokens.EOF.TokenAt("", 7),
	},
	"   2+-2": {
		tokens.Number.TokenAt("2", 4),
		tokens.Plus.TokenAt("+", 5),
		tokens.Number.TokenAt("-2", 6),
		tokens.EOF.TokenAt("", 8),
	},
	"2++2": {
		tokens.Number.TokenAt("2", 1),
		tokens.Plus.TokenAt("+", 2),
		tokens.Number.TokenAt("+2", 3),
		tokens.EOF.TokenAt("", 5),
	},
	"2+++2": {
		tokens.Number.TokenAt("2", 1),
		tokens.Plus.TokenAt("+", 2),
		tokens.Error.TokenAt("expected number, but got '++'", 3),
	},
	"2 - 2": {
		tokens.Number.TokenAt("2", 1),
		tokens.Minus.TokenAt("-", 3),
		tokens.Number.TokenAt("2", 5),
		tokens.EOF.TokenAt("", 6),
	},
	"2-2": {
		tokens.Number.TokenAt("2", 1),
		tokens.Minus.TokenAt("-", 2),
		tokens.Number.TokenAt("2", 3),
		tokens.EOF.TokenAt("", 4),
	},
	"2 - -2": {
		tokens.Number.TokenAt("2", 1),
		tokens.Minus.TokenAt("-", 3),
		tokens.Number.TokenAt("-2", 5),
		tokens.EOF.TokenAt("", 7),
	},
	"2--2": {
		tokens.Number.TokenAt("2", 1),
		tokens.Minus.TokenAt("-", 2),
		tokens.Number.TokenAt("-2", 3),
		tokens.EOF.TokenAt("", 5),
	},
	"2-+2": {
		tokens.Number.TokenAt("2", 1),
		tokens.Minus.TokenAt("-", 2),
		tokens.Number.TokenAt("+2", 3),
		tokens.EOF.TokenAt("", 5),
	},
	"2 --- 2": {
		tokens.Number.TokenAt("2", 1),
		tokens.Minus.TokenAt("-", 3),
		tokens.Error.TokenAt("expected number, but got '--'", 4),
	},
	"   2*2": {
		tokens.Number.TokenAt("2", 4),
		tokens.Multiply.TokenAt("*", 5),
		tokens.Number.TokenAt("2", 6),
		tokens.EOF.TokenAt("", 7),
	},
	"   2 **   2   ": {
		tokens.Number.TokenAt("2", 4),
		tokens.Multiply.TokenAt("*", 6),
		tokens.Error.TokenAt("expected number, but got '*'", 7),
	},
	"2/2": {
		tokens.Number.TokenAt("2", 1),
		tokens.Divide.TokenAt("/", 2),
		tokens.Number.TokenAt("2", 3),
		tokens.EOF.TokenAt("", 4),
	},
	"2//2": {
		tokens.Number.TokenAt("2", 1),
		tokens.Divide.TokenAt("/", 2),
		tokens.Error.TokenAt("expected number, but got '/'", 3),
	},
	"245 lbs": {
		tokens.Number.TokenAt("245", 1),
		tokens.Units.TokenAt("lbs", 5),
		tokens.EOF.TokenAt("", 8),
	},
	"    245 lbs": {
		tokens.Number.TokenAt("245", 5),
		tokens.Units.TokenAt("lbs", 9),
		tokens.EOF.TokenAt("", 12),
	},
	"245lbs": {
		tokens.Number.TokenAt("245", 1),
		tokens.Units.TokenAt("lbs", 4),
		tokens.EOF.TokenAt("", 7),
	},
	"245 lbs + 37.50kg": {
		tokens.Number.TokenAt("245", 1),
		tokens.Units.TokenAt("lbs", 5),
		tokens.Plus.TokenAt("+", 9),
		tokens.Number.TokenAt("37.50", 11),
		tokens.Units.TokenAt("kg", 16),
		tokens.EOF.TokenAt("", 18),
	},
	"245   lbs   + 37.50   kg": {
		tokens.Number.TokenAt("245", 1),
		tokens.Units.TokenAt("lbs", 7),
		tokens.Plus.TokenAt("+", 13),
		tokens.Number.TokenAt("37.50", 15),
		tokens.Units.TokenAt("kg", 23),
		tokens.EOF.TokenAt("", 25),
	},
	"20 lbs in kg": {
		tokens.Number.TokenAt("20", 1),
		tokens.Units.TokenAt("lbs", 4),
		tokens.In.TokenAt("in", 8),
		tokens.Units.TokenAt("kg", 11),
		tokens.EOF.TokenAt("", 13),
	},
	"   20lbs in   kg   ": {
		tokens.Number.TokenAt("20", 4),
		tokens.Units.TokenAt("lbs", 6),
		tokens.In.TokenAt("in", 10),
		tokens.Units.TokenAt("kg", 15),
		tokens.EOF.TokenAt("", 20),
	},
	"20 ints": {
		tokens.Number.TokenAt("20", 1),
		tokens.Units.TokenAt("ints", 4),
		tokens.EOF.TokenAt("", 8),
	},
	"   20ints   ": {
		tokens.Number.TokenAt("20", 4),
		tokens.Units.TokenAt("ints", 6),
		tokens.EOF.TokenAt("", 13),
	},
	"245 lbs + 37.50 kg in kg": {
		tokens.Number.TokenAt("245", 1),
		tokens.Units.TokenAt("lbs", 5),
		tokens.Plus.TokenAt("+", 9),
		tokens.Number.TokenAt("37.50", 11),
		tokens.Units.TokenAt("kg", 17),
		tokens.In.TokenAt("in", 20),
		tokens.Units.TokenAt("kg", 23),
		tokens.EOF.TokenAt("", 25),
	},
	"2 oz + 3 oz + 4 oz": {
		tokens.Number.TokenAt("2", 1),
		tokens.Units.TokenAt("oz", 3),
		tokens.Plus.TokenAt("+", 6),
		tokens.Number.TokenAt("3", 8),
		tokens.Units.TokenAt("oz", 10),
		tokens.Plus.TokenAt("+", 13),
		tokens.Number.TokenAt("4", 15),
		tokens.Units.TokenAt("oz", 17),
		tokens.EOF.TokenAt("", 19),
	},
	"  2 + 2": {
		tokens.Number.TokenAt("2", 3),
		tokens.Plus.TokenAt("+", 5),
		tokens.Number.TokenAt("2", 7),
		tokens.EOF.TokenAt("", 8),
	},
}

func TestTokens(t *testing.T) {
	for input, expected := range testCases {
		output := io.NewTokenChannel()
		go Tokenize(input, &output)
		for _, expect := range expected {
			actual, err := output.ReadToken()
			assert.Nil(t, err)
			assert.Equal(t, expect, actual, "'%s'", input)
		}
	}
}
