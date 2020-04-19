package toks

import (
	"fmt"
	"reflect"
	"unicode"
)

// Token a tokenizer emits tokens.
type Token struct {
	typ TokenType // the type like numberToken
	val string    // the value, like "46.2"
}

// TokenType the type of an emitted token
type TokenType int

const (
	// Error An error occurred.
	Error TokenType = iota
	// EOF At end of input.
	EOF
	// Plus Addition as in '+'.
	Plus
	// Minus Subtraction as in '-'.
	Minus
	// Multiply Multiplication as in '*'.
	Multiply
	// Divide Division as in '/'/
	Divide
	// Number A numeral value like 23.
	Number
	// Units The units of measure like 'kg' or 'pounds'.
	Units
)

// Token Returns a new Token of this type.
func (t TokenType) Token(value string) Token {
	return Token{typ: t, val: value}
}

func (t Token) String() string {
	switch t.typ {
	case Error:
		return fmt.Sprintf("ERR[%s]", t.val)
	case EOF:
		return "EOF"
	case Number:
		return fmt.Sprintf("NUM[%s]", t.val)
	case Plus, Minus, Multiply, Divide:
		return fmt.Sprintf("SYM[%s]", t.val)
	default:
		return fmt.Sprintf("TOK[%s]", t.val)
	}
}

func (t TokenType) String() string {
	return fmt.Sprintf("%s", reflect.TypeOf(t))
}

func expectNumber(tok *Tokenizer) stateFn {
	tok.ignoreSpaceRun()

	// optional sign
	tok.accept("+-")
	tok.acceptRun(" ")

	// expect decimal values
	digits := "0123456789,"
	decimal := true

	leadZero := tok.accept("0")
	if leadZero && tok.accept("xX") {
		// expect hexadecimal values
		digits = "0123456789ABCDEF"
		decimal = false
	}

	// avoid leading commas
	if tok.accept(",") {
		tok.next()
		return tok.error("expected number, but got '%s'", tok.current())
	}

	// accept a run of digits
	count := tok.acceptRun(digits)

	// validate that we have valid digits
	invalidHex := !decimal && count <= 0
	invalidDec := decimal && !leadZero && count <= 0
	if invalidDec || invalidHex {
		tok.next()
		return tok.error("expected number, but got '%s'", tok.current())
	}

	// floating point number
	if tok.accept(".") {
		tok.acceptRun(digits)
	}

	// scientific notation
	if tok.accept("eE") {
		tok.accept("+-")
		tok.acceptRun("0123456789")
	}

	// we have the number
	tok.emit(Number)

	// what is next?
	tok.ignoreSpaceRun()
	next := tok.peek()
	switch {
	case next == -1, next == '\n':
		return expectEOF
	case unicode.IsLetter(next):
		return expectUnits
	default:
		return expectSymbol
	}
}

func expectEOF(tok *Tokenizer) stateFn {
	tok.ignoreSpaceRun()
	tok.ignoreRun('\n')
	if tok.next() == -1 {
		tok.emit(EOF)
		return nil
	}
	return tok.error("expected EOF, but got '%s'", tok.current())
}

func expectSymbol(tok *Tokenizer) stateFn {
	for {
		switch next := tok.next(); {
		case next == '+':
			tok.emit(Plus)
			return expectNumber
		case next == '-':
			tok.emit(Minus)
			return expectNumber
		case next == '*':
			tok.emit(Multiply)
			return expectNumber
		case next == '/':
			tok.emit(Divide)
			return expectNumber
		case unicode.IsSpace(next):
			tok.ignore()
		case next == -1:
			tok.backup()
			return expectEOF
		default:
			return tok.error("expected symbol, but got '%s'", tok.current())
		}
	}
}

func expectUnits(tok *Tokenizer) stateFn {
	count := tok.acceptLetterRun()
	if count <= 0 {
		tok.next()
		return tok.error("expected units, but got '%s'", tok.current())
	}
	tok.emit(Units)
	return expectSymbol
}
