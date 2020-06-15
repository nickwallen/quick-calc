package tokenizer

import (
	"github.com/nickwallen/quick-calc/internal/io"
	"github.com/nickwallen/quick-calc/internal/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testCases = map[string][]types.Token{
	"": {
		types.Error.TokenAt("expected number, but got ''", 1),
	},
	" ": {
		types.Error.TokenAt("expected number, but got ''", 2),
	},
	"  22": {
		types.Number.TokenAt("22", 3),
		types.EOF.TokenAt("", 5),
	},
	"0": {
		types.Number.TokenAt("0", 1),
		types.EOF.TokenAt("", 2),
	},
	"2,200,123": {
		types.Number.TokenAt("2,200,123", 1),
		types.EOF.TokenAt("", 10),
	},
	",200,200": {
		types.Error.TokenAt("expected number, but got ',2'", 1),
	},
	"  +22": {
		types.Number.TokenAt("+22", 3),
		types.EOF.TokenAt("", 6),
	},
	"  -22": {
		types.Number.TokenAt("-22", 3),
		types.EOF.TokenAt("", 6),
	},
	"2?": {
		types.Number.TokenAt("2", 1),
		types.Error.TokenAt("expected symbol, but got '?'", 2),
	},
	"0xAF   ": {
		types.Number.TokenAt("0xAF", 1),
		types.EOF.TokenAt("", 8),
	},
	"0xG2": {
		types.Error.TokenAt("expected number, but got '0xG'", 1),
	},
	"0x   ": {
		types.Error.TokenAt("expected number, but got '0x'", 1),
	},
	"   2.22": {
		types.Number.TokenAt("2.22", 4),
		types.EOF.TokenAt("", 8),
	},
	"2E10   ": {
		types.Number.TokenAt("2E10", 1),
		types.EOF.TokenAt("", 8),
	},
	"2+2": {
		types.Number.TokenAt("2", 1),
		types.Plus.TokenAt("+", 2),
		types.Number.TokenAt("2", 3),
		types.EOF.TokenAt("", 4),
	},
	"2 + -2": {
		types.Number.TokenAt("2", 1),
		types.Plus.TokenAt("+", 3),
		types.Number.TokenAt("-2", 5),
		types.EOF.TokenAt("", 7),
	},
	"   2+-2": {
		types.Number.TokenAt("2", 4),
		types.Plus.TokenAt("+", 5),
		types.Number.TokenAt("-2", 6),
		types.EOF.TokenAt("", 8),
	},
	"2++2": {
		types.Number.TokenAt("2", 1),
		types.Plus.TokenAt("+", 2),
		types.Number.TokenAt("+2", 3),
		types.EOF.TokenAt("", 5),
	},
	"2+++2": {
		types.Number.TokenAt("2", 1),
		types.Plus.TokenAt("+", 2),
		types.Error.TokenAt("expected number, but got '++'", 3),
	},
	"2 - 2": {
		types.Number.TokenAt("2", 1),
		types.Minus.TokenAt("-", 3),
		types.Number.TokenAt("2", 5),
		types.EOF.TokenAt("", 6),
	},
	"2-2": {
		types.Number.TokenAt("2", 1),
		types.Minus.TokenAt("-", 2),
		types.Number.TokenAt("2", 3),
		types.EOF.TokenAt("", 4),
	},
	"2 - -2": {
		types.Number.TokenAt("2", 1),
		types.Minus.TokenAt("-", 3),
		types.Number.TokenAt("-2", 5),
		types.EOF.TokenAt("", 7),
	},
	"2--2": {
		types.Number.TokenAt("2", 1),
		types.Minus.TokenAt("-", 2),
		types.Number.TokenAt("-2", 3),
		types.EOF.TokenAt("", 5),
	},
	"2-+2": {
		types.Number.TokenAt("2", 1),
		types.Minus.TokenAt("-", 2),
		types.Number.TokenAt("+2", 3),
		types.EOF.TokenAt("", 5),
	},
	"2 --- 2": {
		types.Number.TokenAt("2", 1),
		types.Minus.TokenAt("-", 3),
		types.Error.TokenAt("expected number, but got '--'", 4),
	},
	"   2*2": {
		types.Number.TokenAt("2", 4),
		types.Multiply.TokenAt("*", 5),
		types.Number.TokenAt("2", 6),
		types.EOF.TokenAt("", 7),
	},
	"   2 **   2   ": {
		types.Number.TokenAt("2", 4),
		types.Multiply.TokenAt("*", 6),
		types.Error.TokenAt("expected number, but got '*'", 7),
	},
	"2/2": {
		types.Number.TokenAt("2", 1),
		types.Divide.TokenAt("/", 2),
		types.Number.TokenAt("2", 3),
		types.EOF.TokenAt("", 4),
	},
	"2//2": {
		types.Number.TokenAt("2", 1),
		types.Divide.TokenAt("/", 2),
		types.Error.TokenAt("expected number, but got '/'", 3),
	},
	"245 lbs": {
		types.Number.TokenAt("245", 1),
		types.Units.TokenAt("lbs", 5),
		types.EOF.TokenAt("", 8),
	},
	"    245 lbs": {
		types.Number.TokenAt("245", 5),
		types.Units.TokenAt("lbs", 9),
		types.EOF.TokenAt("", 12),
	},
	"245lbs": {
		types.Number.TokenAt("245", 1),
		types.Units.TokenAt("lbs", 4),
		types.EOF.TokenAt("", 7),
	},
	"245 lbs + 37.50kg": {
		types.Number.TokenAt("245", 1),
		types.Units.TokenAt("lbs", 5),
		types.Plus.TokenAt("+", 9),
		types.Number.TokenAt("37.50", 11),
		types.Units.TokenAt("kg", 16),
		types.EOF.TokenAt("", 18),
	},
	"245   lbs   + 37.50   kg": {
		types.Number.TokenAt("245", 1),
		types.Units.TokenAt("lbs", 7),
		types.Plus.TokenAt("+", 13),
		types.Number.TokenAt("37.50", 15),
		types.Units.TokenAt("kg", 23),
		types.EOF.TokenAt("", 25),
	},
	"20 lbs in kg": {
		types.Number.TokenAt("20", 1),
		types.Units.TokenAt("lbs", 4),
		types.In.TokenAt("in", 8),
		types.Units.TokenAt("kg", 11),
		types.EOF.TokenAt("", 13),
	},
	"   20lbs in   kg   ": {
		types.Number.TokenAt("20", 4),
		types.Units.TokenAt("lbs", 6),
		types.In.TokenAt("in", 10),
		types.Units.TokenAt("kg", 15),
		types.EOF.TokenAt("", 20),
	},
	"20 ints": {
		types.Number.TokenAt("20", 1),
		types.Units.TokenAt("ints", 4),
		types.EOF.TokenAt("", 8),
	},
	"   20ints   ": {
		types.Number.TokenAt("20", 4),
		types.Units.TokenAt("ints", 6),
		types.EOF.TokenAt("", 13),
	},
	"245 lbs + 37.50 kg in kg": {
		types.Number.TokenAt("245", 1),
		types.Units.TokenAt("lbs", 5),
		types.Plus.TokenAt("+", 9),
		types.Number.TokenAt("37.50", 11),
		types.Units.TokenAt("kg", 17),
		types.In.TokenAt("in", 20),
		types.Units.TokenAt("kg", 23),
		types.EOF.TokenAt("", 25),
	},
	"2 oz + 3 oz + 4 oz": {
		types.Number.TokenAt("2", 1),
		types.Units.TokenAt("oz", 3),
		types.Plus.TokenAt("+", 6),
		types.Number.TokenAt("3", 8),
		types.Units.TokenAt("oz", 10),
		types.Plus.TokenAt("+", 13),
		types.Number.TokenAt("4", 15),
		types.Units.TokenAt("oz", 17),
		types.EOF.TokenAt("", 19),
	},
	"  2 + 2": {
		types.Number.TokenAt("2", 3),
		types.Plus.TokenAt("+", 5),
		types.Number.TokenAt("2", 7),
		types.EOF.TokenAt("", 8),
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
