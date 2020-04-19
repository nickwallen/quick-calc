package toks

import (
	"fmt"
	"reflect"
	"strings"
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
	tok.ignoreWhitespaces()

	// optional sign
	tok.accept("+-")
	tok.acceptRun(" ")

	// accept hex or base-10
	digits := "0123456789"
	if tok.accept("0") && tok.accept("xX") {
		digits = "0123456789ABCDEF"
	}

	// at least one digit is required
	if strings.IndexRune(digits+"xX", tok.peek()) < 0 {
		tok.next()
		return tok.errorf("expected number, but got '%s'", tok.current())
	}
	tok.acceptRun(digits)

	// floating point number
	if tok.accept(".") {
		tok.acceptRun(digits)
	}

	// scientific notation
	if tok.accept("eE") {
		tok.accept("+-")
		tok.acceptRun("0123456789")
	}

	if unicode.IsLetter(tok.peek()) || unicode.IsNumber(tok.peek()) {
		tok.next()
		return tok.errorf("expected number, but got '%s'", tok.current())
	}

	// we have the number
	tok.emit(Number)

	// what is next?
	switch tok.peek() {
	case -1, '\n':
		return expectEOF
	default:
		return expectSymbol
	}
}

func expectEOF(tok *Tokenizer) stateFn {
	tok.ignoreWhitespaces()
	tok.ignoreRun('\n')
	if tok.next() == -1 {
		tok.emit(EOF)
		return nil
	}
	return tok.errorf("expected EOF, but got '%s'", tok.current())
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
			return tok.errorf("expected symbol, but got '%s'", tok.current())
		}
	}
}
