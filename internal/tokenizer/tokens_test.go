package tokenizer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testCases = map[string][]Token{
	"":                         {Error.Token("expected number, but got ''")},
	" ":                        {Error.Token("expected number, but got ''")},
	"22":                       {Number.Token("22"), EOF.Token("")},
	"  22":                     {Number.Token("22"), EOF.Token("")},
	"22    ":                   {Number.Token("22"), EOF.Token("")},
	"0":                        {Number.Token("0"), EOF.Token("")},
	"  0":                      {Number.Token("0"), EOF.Token("")},
	"0    ":                    {Number.Token("0"), EOF.Token("")},
	"2,200,123":                {Number.Token("2,200,123"), EOF.Token("")},
	"  2,200,123":              {Number.Token("2,200,123"), EOF.Token("")},
	"2,200,123    ":            {Number.Token("2,200,123"), EOF.Token("")},
	",200,200":                 {Error.Token("expected number, but got ',2'")},
	"+22":                      {Number.Token("+22"), EOF.Token("")},
	"  +22":                    {Number.Token("+22"), EOF.Token("")},
	"+22    ":                  {Number.Token("+22"), EOF.Token("")},
	"+  22":                    {Number.Token("+22"), EOF.Token("")},
	"-22":                      {Number.Token("-22"), EOF.Token("")},
	"  -22":                    {Number.Token("-22"), EOF.Token("")},
	"-22    ":                  {Number.Token("-22"), EOF.Token("")},
	"  - 22":                   {Number.Token("-22"), EOF.Token("")},
	"2?":                       {Number.Token("2"), Error.Token("expected symbol, but got '?'")},
	"   2?":                    {Number.Token("2"), Error.Token("expected symbol, but got '?'")},
	"2?   ":                    {Number.Token("2"), Error.Token("expected symbol, but got '?'")},
	"0xAF":                     {Number.Token("0xAF"), EOF.Token("")},
	"   0xAF":                  {Number.Token("0xAF"), EOF.Token("")},
	"0xAF   ":                  {Number.Token("0xAF"), EOF.Token("")},
	"0xG2":                     {Error.Token("expected number, but got '0xG'")},
	"   0xG2":                  {Error.Token("expected number, but got '0xG'")},
	"0xG2   ":                  {Error.Token("expected number, but got '0xG'")},
	"0x":                       {Error.Token("expected number, but got '0x'")},
	"   0x":                    {Error.Token("expected number, but got '0x'")},
	"0x   ":                    {Error.Token("expected number, but got '0x'")},
	"2.22":                     {Number.Token("2.22"), EOF.Token("")},
	"   2.22":                  {Number.Token("2.22"), EOF.Token("")},
	"2.22   ":                  {Number.Token("2.22"), EOF.Token("")},
	"2E10":                     {Number.Token("2E10"), EOF.Token("")},
	"   2E10":                  {Number.Token("2E10"), EOF.Token("")},
	"2E10   ":                  {Number.Token("2E10"), EOF.Token("")},
	"2 + 2":                    {Number.Token("2"), Plus.Token("+"), Number.Token("2"), EOF.Token("")},
	"   2+2":                   {Number.Token("2"), Plus.Token("+"), Number.Token("2"), EOF.Token("")},
	"   2 +   2   ":            {Number.Token("2"), Plus.Token("+"), Number.Token("2"), EOF.Token("")},
	"2+2":                      {Number.Token("2"), Plus.Token("+"), Number.Token("2"), EOF.Token("")},
	"2 + -2":                   {Number.Token("2"), Plus.Token("+"), Number.Token("-2"), EOF.Token("")},
	"   2+-2":                  {Number.Token("2"), Plus.Token("+"), Number.Token("-2"), EOF.Token("")},
	"   2 +   -2   ":           {Number.Token("2"), Plus.Token("+"), Number.Token("-2"), EOF.Token("")},
	"2+-2":                     {Number.Token("2"), Plus.Token("+"), Number.Token("-2"), EOF.Token("")},
	"2 + +2":                   {Number.Token("2"), Plus.Token("+"), Number.Token("+2"), EOF.Token("")},
	"   2++2":                  {Number.Token("2"), Plus.Token("+"), Number.Token("+2"), EOF.Token("")},
	"   2 +   +2   ":           {Number.Token("2"), Plus.Token("+"), Number.Token("+2"), EOF.Token("")},
	"2++2":                     {Number.Token("2"), Plus.Token("+"), Number.Token("+2"), EOF.Token("")},
	"2 +++ 2":                  {Number.Token("2"), Plus.Token("+"), Error.Token("expected number, but got '++'")},
	"   2+++2":                 {Number.Token("2"), Plus.Token("+"), Error.Token("expected number, but got '++'")},
	"   2+++   2   ":           {Number.Token("2"), Plus.Token("+"), Error.Token("expected number, but got '++'")},
	"2+++2":                    {Number.Token("2"), Plus.Token("+"), Error.Token("expected number, but got '++'")},
	"2 - 2":                    {Number.Token("2"), Minus.Token("-"), Number.Token("2"), EOF.Token("")},
	"   2-2":                   {Number.Token("2"), Minus.Token("-"), Number.Token("2"), EOF.Token("")},
	"   2 -   2   ":            {Number.Token("2"), Minus.Token("-"), Number.Token("2"), EOF.Token("")},
	"2-2":                      {Number.Token("2"), Minus.Token("-"), Number.Token("2"), EOF.Token("")},
	"2 - -2":                   {Number.Token("2"), Minus.Token("-"), Number.Token("-2"), EOF.Token("")},
	"   2--2":                  {Number.Token("2"), Minus.Token("-"), Number.Token("-2"), EOF.Token("")},
	"   2 -   -2   ":           {Number.Token("2"), Minus.Token("-"), Number.Token("-2"), EOF.Token("")},
	"2--2":                     {Number.Token("2"), Minus.Token("-"), Number.Token("-2"), EOF.Token("")},
	"2 - +2":                   {Number.Token("2"), Minus.Token("-"), Number.Token("+2"), EOF.Token("")},
	"   2-+2":                  {Number.Token("2"), Minus.Token("-"), Number.Token("+2"), EOF.Token("")},
	"   2 -   +2   ":           {Number.Token("2"), Minus.Token("-"), Number.Token("+2"), EOF.Token("")},
	"2-+2":                     {Number.Token("2"), Minus.Token("-"), Number.Token("+2"), EOF.Token("")},
	"2 --- 2":                  {Number.Token("2"), Minus.Token("-"), Error.Token("expected number, but got '--'")},
	"   2---2":                 {Number.Token("2"), Minus.Token("-"), Error.Token("expected number, but got '--'")},
	"   2 ---   2   ":          {Number.Token("2"), Minus.Token("-"), Error.Token("expected number, but got '--'")},
	"2---2":                    {Number.Token("2"), Minus.Token("-"), Error.Token("expected number, but got '--'")},
	"2 * 2":                    {Number.Token("2"), Multiply.Token("*"), Number.Token("2"), EOF.Token("")},
	"   2*2":                   {Number.Token("2"), Multiply.Token("*"), Number.Token("2"), EOF.Token("")},
	"   2 *   2   ":            {Number.Token("2"), Multiply.Token("*"), Number.Token("2"), EOF.Token("")},
	"2*2":                      {Number.Token("2"), Multiply.Token("*"), Number.Token("2"), EOF.Token("")},
	"2 ** 2":                   {Number.Token("2"), Multiply.Token("*"), Error.Token("expected number, but got '*'")},
	"   2**2":                  {Number.Token("2"), Multiply.Token("*"), Error.Token("expected number, but got '*'")},
	"   2 **   2   ":           {Number.Token("2"), Multiply.Token("*"), Error.Token("expected number, but got '*'")},
	"2**2":                     {Number.Token("2"), Multiply.Token("*"), Error.Token("expected number, but got '*'")},
	"2 / 2":                    {Number.Token("2"), Divide.Token("/"), Number.Token("2"), EOF.Token("")},
	"   2/2":                   {Number.Token("2"), Divide.Token("/"), Number.Token("2"), EOF.Token("")},
	"   2 /   2   ":            {Number.Token("2"), Divide.Token("/"), Number.Token("2"), EOF.Token("")},
	"2/2":                      {Number.Token("2"), Divide.Token("/"), Number.Token("2"), EOF.Token("")},
	"2 // 2":                   {Number.Token("2"), Divide.Token("/"), Error.Token("expected number, but got '/'")},
	"   2//2":                  {Number.Token("2"), Divide.Token("/"), Error.Token("expected number, but got '/'")},
	"   2 //   2   ":           {Number.Token("2"), Divide.Token("/"), Error.Token("expected number, but got '/'")},
	"2//2":                     {Number.Token("2"), Divide.Token("/"), Error.Token("expected number, but got '/'")},
	"245 lbs":                  {Number.Token("245"), Units.Token("lbs"), EOF.Token("")},
	"    245 lbs":              {Number.Token("245"), Units.Token("lbs"), EOF.Token("")},
	"245lbs":                   {Number.Token("245"), Units.Token("lbs"), EOF.Token("")},
	"245 lbs + 37.50kg":        {Number.Token("245"), Units.Token("lbs"), Plus.Token("+"), Number.Token("37.50"), Units.Token("kg"), EOF.Token("")},
	"245   lbs   + 37.50   kg": {Number.Token("245"), Units.Token("lbs"), Plus.Token("+"), Number.Token("37.50"), Units.Token("kg"), EOF.Token("")},
	"20 lbs in kg":             {Number.Token("20"), Units.Token("lbs"), In.Token("in"), Units.Token("kg"), EOF.Token("")},
	"   20lbs in   kg   ":      {Number.Token("20"), Units.Token("lbs"), In.Token("in"), Units.Token("kg"), EOF.Token("")},
	"20 ints":                  {Number.Token("20"), Units.Token("ints"), EOF.Token("")},
	"   20ints   ":             {Number.Token("20"), Units.Token("ints"), EOF.Token("")},
	"245 lbs + 37.50 kg in kg": {Number.Token("245"), Units.Token("lbs"), Plus.Token("+"), Number.Token("37.50"), Units.Token("kg"), In.Token("in"), Units.Token("kg"), EOF.Token("")},
}

func TestTokens(t *testing.T) {
	for input, expected := range testCases {
		output := make(chan Token, 2)
		go Tokenize(input, output)
		expect(t, expected, output)
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
