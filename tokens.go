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
	// ErrorToken An error occurred.
	ErrorToken TokenType = iota

	// EOFToken At end of input.
	EOFToken

	// PlusToken Addition as in '+'.
	PlusToken

	// MinusToken Subtraction as in '-'.
	MinusToken

	// MultiplyToken Multiplication as in '*'.
	MultiplyToken

	// DivisionToken Division as in '/'/
	DivisionToken

	// NumberToken A numeral value like 23.
	NumberToken
)

// Of Returns a new Token of this type.
func (t TokenType) Of(value string) Token {
	return Token{typ: t, val: value}
}

func (t Token) String() string {
	switch t.typ {
	case ErrorToken:
		return fmt.Sprintf("ERR(%s)", t.val)
	case EOFToken:
		return "EOF"
	case NumberToken:
		return fmt.Sprintf("NUM(%q)", t.val)
	case PlusToken, MinusToken:
		return fmt.Sprintf("SYM(%q)", t.val)
	default:
		return fmt.Sprintf("TOK(%q)", t.val)
	}
}

func (t TokenType) String() string {
	return fmt.Sprintf("%s", reflect.TypeOf(t))
}

func expectNumber(tok *Tokenizer) stateFn {
	tok.ignoreWhitespaces()

	// optional sign
	// tok.accept("+-")

	// accept hex or base-10
	digits := "0123456789"
	if tok.accept("0") && tok.accept("xX") {
		digits = "0123456789ABCDEF"
	}

	// at least one digit is required
	if strings.IndexRune(digits+"xX", tok.peek()) < 0 {
		tok.next()
		return tok.errorf("expected number, but got %q", tok.current())
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
		return tok.errorf("expected number, but got %q", tok.current())
	}

	// we have the number
	tok.emit(NumberToken)

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
		tok.emit(EOFToken)
		return nil
	}
	return tok.errorf("expected EOF: %q", tok.current())
}

func expectSymbol(tok *Tokenizer) stateFn {
	for {
		switch next := tok.next(); {
		case next == '+':
			tok.emit(PlusToken)
			return expectNumber
		case next == '-':
			tok.emit(MinusToken)
			return expectNumber
		case next == '*':
			tok.emit(MultiplyToken)
			return expectNumber
		case next == '/':
			tok.emit(DivisionToken)
			return expectNumber
		case unicode.IsSpace(next):
			tok.ignore()
		case next == -1:
			tok.backup()
			return expectEOF
		default:
			return tok.errorf("expected symbol, but got '%q'", tok.current())
		}
	}
}
