package toks

import (
	"fmt"
	"unicode"
)

// Token a Tokenizer emits tokens.
type Token struct {
	TokenType TokenType // the type like numberToken
	Value     string    // the value, like "46.2"
}

// TokenType the type of an emitted Token
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
	// In A symbol for conversions; 23 lbs in kg.
	In
	// Number A numeral value like 23.
	Number
	// Units The units of measure like 'kg' or 'pounds'.
	Units
)

func (t TokenType) String() string {
	switch t {
	case Error:
		return "error"
	case EOF:
		return "end-of-line"
	case Plus:
		return "addition (+)"
	case Minus:
		return "minus (-)"
	case Multiply:
		return "multiply (*)"
	case Divide:
		return "division (/)"
	case In:
		return "keyword 'in'"
	case Number:
		return "number"
	case Units:
		return "units"
	default:
		return "unknown"
	}
}

// Token Returns a new Token of this type.
func (t TokenType) token(value string) Token {
	return Token{TokenType: t, Value: value}
}

func (t Token) String() string {
	switch t.TokenType {
	case Error:
		return fmt.Sprintf("ERR[%s]", t.Value)
	case EOF:
		return "EOF"
	case Number:
		return fmt.Sprintf("NUM[%s]", t.Value)
	case Plus, Minus, Multiply, Divide:
		return fmt.Sprintf("SYM[%s]", t.Value)
	case Units:
		return fmt.Sprintf("UNI[%s]", t.Value)
	default:
		return fmt.Sprintf("TOK[%s]", t.Value)
	}
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
	case next == eofRune, next == '\n':
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
	if tok.next() == eofRune {
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
		case next == eofRune:
			tok.backup()
			return expectEOF
		default:
			return tok.error("expected symbol, but got '%s'", tok.current())
		}
	}
}

func expectUnits(tok *Tokenizer) stateFn {
	tok.ignoreSpaceRun()
	count := tok.acceptLetterRun()
	if count <= 0 {
		tok.next()
		return tok.error("expected units, but got '%s'", tok.current())
	}
	tok.emit(Units)

	// what is next?
	tok.ignoreSpaceRun()
	switch {
	case tok.accept("iI") && tok.accept("nN") && tok.accept(" "):
		tok.backup()
		tok.backup()
		tok.backup()
		return expectIn
	case unicode.IsLetter(tok.peek()):
		return expectUnits
	default:
		return expectSymbol
	}
}

func expectIn(tok *Tokenizer) stateFn {
	tok.ignoreSpaceRun()
	if tok.accept("iI") && tok.accept("nN") && tok.accept(" ") {
		tok.emit(In)
		return expectUnits
	}
	// error
	tok.next()
	return tok.error("expected 'in' keyword, but got '%s'", tok.current())
}
