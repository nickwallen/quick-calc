package toks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testCases = map[string][]Token{
	"":                         {Error.token("expected number, but got ''")},
	" ":                        {Error.token("expected number, but got ''")},
	"22":                       {Number.token("22"), EOF.token("")},
	"  22":                     {Number.token("22"), EOF.token("")},
	"22    ":                   {Number.token("22"), EOF.token("")},
	"0":                        {Number.token("0"), EOF.token("")},
	"  0":                      {Number.token("0"), EOF.token("")},
	"0    ":                    {Number.token("0"), EOF.token("")},
	"2,200,123":                {Number.token("2,200,123"), EOF.token("")},
	"  2,200,123":              {Number.token("2,200,123"), EOF.token("")},
	"2,200,123    ":            {Number.token("2,200,123"), EOF.token("")},
	",200,200":                 {Error.token("expected number, but got ',2'")},
	"+22":                      {Number.token("+22"), EOF.token("")},
	"  +22":                    {Number.token("+22"), EOF.token("")},
	"+22    ":                  {Number.token("+22"), EOF.token("")},
	"+  22":                    {Number.token("+22"), EOF.token("")},
	"-22":                      {Number.token("-22"), EOF.token("")},
	"  -22":                    {Number.token("-22"), EOF.token("")},
	"-22    ":                  {Number.token("-22"), EOF.token("")},
	"  - 22":                   {Number.token("-22"), EOF.token("")},
	"2?":                       {Number.token("2"), Error.token("expected symbol, but got '?'")},
	"   2?":                    {Number.token("2"), Error.token("expected symbol, but got '?'")},
	"2?   ":                    {Number.token("2"), Error.token("expected symbol, but got '?'")},
	"0xAF":                     {Number.token("0xAF"), EOF.token("")},
	"   0xAF":                  {Number.token("0xAF"), EOF.token("")},
	"0xAF   ":                  {Number.token("0xAF"), EOF.token("")},
	"0xG2":                     {Error.token("expected number, but got '0xG'")},
	"   0xG2":                  {Error.token("expected number, but got '0xG'")},
	"0xG2   ":                  {Error.token("expected number, but got '0xG'")},
	"0x":                       {Error.token("expected number, but got '0x'")},
	"   0x":                    {Error.token("expected number, but got '0x'")},
	"0x   ":                    {Error.token("expected number, but got '0x'")},
	"2.22":                     {Number.token("2.22"), EOF.token("")},
	"   2.22":                  {Number.token("2.22"), EOF.token("")},
	"2.22   ":                  {Number.token("2.22"), EOF.token("")},
	"2E10":                     {Number.token("2E10"), EOF.token("")},
	"   2E10":                  {Number.token("2E10"), EOF.token("")},
	"2E10   ":                  {Number.token("2E10"), EOF.token("")},
	"2 + 2":                    {Number.token("2"), Plus.token("+"), Number.token("2"), EOF.token("")},
	"   2+2":                   {Number.token("2"), Plus.token("+"), Number.token("2"), EOF.token("")},
	"   2 +   2   ":            {Number.token("2"), Plus.token("+"), Number.token("2"), EOF.token("")},
	"2+2":                      {Number.token("2"), Plus.token("+"), Number.token("2"), EOF.token("")},
	"2 + -2":                   {Number.token("2"), Plus.token("+"), Number.token("-2"), EOF.token("")},
	"   2+-2":                  {Number.token("2"), Plus.token("+"), Number.token("-2"), EOF.token("")},
	"   2 +   -2   ":           {Number.token("2"), Plus.token("+"), Number.token("-2"), EOF.token("")},
	"2+-2":                     {Number.token("2"), Plus.token("+"), Number.token("-2"), EOF.token("")},
	"2 + +2":                   {Number.token("2"), Plus.token("+"), Number.token("+2"), EOF.token("")},
	"   2++2":                  {Number.token("2"), Plus.token("+"), Number.token("+2"), EOF.token("")},
	"   2 +   +2   ":           {Number.token("2"), Plus.token("+"), Number.token("+2"), EOF.token("")},
	"2++2":                     {Number.token("2"), Plus.token("+"), Number.token("+2"), EOF.token("")},
	"2 +++ 2":                  {Number.token("2"), Plus.token("+"), Error.token("expected number, but got '++'")},
	"   2+++2":                 {Number.token("2"), Plus.token("+"), Error.token("expected number, but got '++'")},
	"   2+++   2   ":           {Number.token("2"), Plus.token("+"), Error.token("expected number, but got '++'")},
	"2+++2":                    {Number.token("2"), Plus.token("+"), Error.token("expected number, but got '++'")},
	"2 - 2":                    {Number.token("2"), Minus.token("-"), Number.token("2"), EOF.token("")},
	"   2-2":                   {Number.token("2"), Minus.token("-"), Number.token("2"), EOF.token("")},
	"   2 -   2   ":            {Number.token("2"), Minus.token("-"), Number.token("2"), EOF.token("")},
	"2-2":                      {Number.token("2"), Minus.token("-"), Number.token("2"), EOF.token("")},
	"2 - -2":                   {Number.token("2"), Minus.token("-"), Number.token("-2"), EOF.token("")},
	"   2--2":                  {Number.token("2"), Minus.token("-"), Number.token("-2"), EOF.token("")},
	"   2 -   -2   ":           {Number.token("2"), Minus.token("-"), Number.token("-2"), EOF.token("")},
	"2--2":                     {Number.token("2"), Minus.token("-"), Number.token("-2"), EOF.token("")},
	"2 - +2":                   {Number.token("2"), Minus.token("-"), Number.token("+2"), EOF.token("")},
	"   2-+2":                  {Number.token("2"), Minus.token("-"), Number.token("+2"), EOF.token("")},
	"   2 -   +2   ":           {Number.token("2"), Minus.token("-"), Number.token("+2"), EOF.token("")},
	"2-+2":                     {Number.token("2"), Minus.token("-"), Number.token("+2"), EOF.token("")},
	"2 --- 2":                  {Number.token("2"), Minus.token("-"), Error.token("expected number, but got '--'")},
	"   2---2":                 {Number.token("2"), Minus.token("-"), Error.token("expected number, but got '--'")},
	"   2 ---   2   ":          {Number.token("2"), Minus.token("-"), Error.token("expected number, but got '--'")},
	"2---2":                    {Number.token("2"), Minus.token("-"), Error.token("expected number, but got '--'")},
	"2 * 2":                    {Number.token("2"), Multiply.token("*"), Number.token("2"), EOF.token("")},
	"   2*2":                   {Number.token("2"), Multiply.token("*"), Number.token("2"), EOF.token("")},
	"   2 *   2   ":            {Number.token("2"), Multiply.token("*"), Number.token("2"), EOF.token("")},
	"2*2":                      {Number.token("2"), Multiply.token("*"), Number.token("2"), EOF.token("")},
	"2 ** 2":                   {Number.token("2"), Multiply.token("*"), Error.token("expected number, but got '*'")},
	"   2**2":                  {Number.token("2"), Multiply.token("*"), Error.token("expected number, but got '*'")},
	"   2 **   2   ":           {Number.token("2"), Multiply.token("*"), Error.token("expected number, but got '*'")},
	"2**2":                     {Number.token("2"), Multiply.token("*"), Error.token("expected number, but got '*'")},
	"2 / 2":                    {Number.token("2"), Divide.token("/"), Number.token("2"), EOF.token("")},
	"   2/2":                   {Number.token("2"), Divide.token("/"), Number.token("2"), EOF.token("")},
	"   2 /   2   ":            {Number.token("2"), Divide.token("/"), Number.token("2"), EOF.token("")},
	"2/2":                      {Number.token("2"), Divide.token("/"), Number.token("2"), EOF.token("")},
	"2 // 2":                   {Number.token("2"), Divide.token("/"), Error.token("expected number, but got '/'")},
	"   2//2":                  {Number.token("2"), Divide.token("/"), Error.token("expected number, but got '/'")},
	"   2 //   2   ":           {Number.token("2"), Divide.token("/"), Error.token("expected number, but got '/'")},
	"2//2":                     {Number.token("2"), Divide.token("/"), Error.token("expected number, but got '/'")},
	"245 lbs":                  {Number.token("245"), Units.token("lbs"), EOF.token("")},
	"    245 lbs":              {Number.token("245"), Units.token("lbs"), EOF.token("")},
	"245lbs":                   {Number.token("245"), Units.token("lbs"), EOF.token("")},
	"245 lbs + 37.50kg":        {Number.token("245"), Units.token("lbs"), Plus.token("+"), Number.token("37.50"), Units.token("kg"), EOF.token("")},
	"245   lbs   + 37.50   kg": {Number.token("245"), Units.token("lbs"), Plus.token("+"), Number.token("37.50"), Units.token("kg"), EOF.token("")},
	"20 lbs in kg":             {Number.token("20"), Units.token("lbs"), In.token("in"), Units.token("kg"), EOF.token("")},
	"   20lbs in   kg   ":      {Number.token("20"), Units.token("lbs"), In.token("in"), Units.token("kg"), EOF.token("")},
	"20 ints":                  {Number.token("20"), Units.token("ints"), EOF.token("")},
	"   20ints   ":             {Number.token("20"), Units.token("ints"), EOF.token("")},
	"245 lbs + 37.50 kg in kg": {Number.token("245"), Units.token("lbs"), Plus.token("+"), Number.token("37.50"), Units.token("kg"), In.token("in"), Units.token("kg"), EOF.token("")},
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
