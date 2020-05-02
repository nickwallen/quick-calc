package tokenizer

import (
	"fmt"
	"github.com/nickwallen/quick-calc/internal/tokens"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	eofRune = rune(0)
)

// tokenizer A tokenizer performs lexical analysis on an input string.
type tokenizer struct {
	state  stateFn     // the current state function
	input  string      // the string to scan
	start  int         // start position for this item
	pos    int         // current position in the input
	width  int         // width of the last rune read
	writer tokenWriter // allows the tokenizer to write tokens that it finds
}

// the state of the scanner as a function that returns the next state.
type stateFn func(*tokenizer) stateFn

// the interface used by the tokenizer to write tokens.
type tokenWriter interface {
	WriteToken(token tokens.Token) error
	Close()
}

// Tokenize Tokenize the input string and writes each Token to the output channel.
func Tokenize(input string, writer tokenWriter) {
	tok := &tokenizer{
		state:  expectNumber, // at the start, expect a number
		input:  input,
		writer: writer,
	}
	tok.run()
}

// returns what is currently being scanned
func (tok *tokenizer) current() string {
	// remove leading, trailing, and embedded whitespace
	value := tok.input[tok.start:tok.pos]
	return strings.ReplaceAll(value, " ", "")
}

// emits a Token to be consumed by the client
func (tok *tokenizer) emit(tokenType tokens.TokenType) {
	var token tokens.Token
	switch tokenType {
	case tokens.EOF:
		token = tokens.EOF.Token("")
	default:
		token = tokenType.Token(tok.current())
	}
	tok.writer.WriteToken(token)
	tok.start = tok.pos
}

// skips over the pending input
func (tok *tokenizer) ignore() {
	tok.start = tok.pos
}

// skips over a run of values
func (tok *tokenizer) ignoreRun(ignore rune) {
	for tok.next() == ignore {
		tok.ignore()
	}
}

// skips over any whitespace
func (tok *tokenizer) ignoreSpaceRun() {
	for unicode.IsSpace(tok.next()) {
		tok.ignore()
	}
	tok.backup()
}

// steps back one
func (tok *tokenizer) backup() {
	tok.pos -= tok.width
}

// peek returns, but does not consume the next rune in the input.
func (tok *tokenizer) peek() rune {
	rune := tok.next()
	tok.backup()
	return rune
}

func (tok *tokenizer) next() rune {
	if tok.pos >= len(tok.input) {
		tok.width = 0
		return eofRune
	}
	var r rune
	r, tok.width = utf8.DecodeRuneInString(tok.input[tok.pos:])
	tok.pos += tok.width
	return r
}

// accept consumes the next rune if it is valid.
func (tok *tokenizer) accept(valid string) bool {
	if strings.IndexRune(valid, tok.next()) >= 0 {
		return true
	}
	tok.backup()
	return false
}

// acceptRun consumes a run of strings
func (tok *tokenizer) acceptRun(valid string) (count int) {
	for strings.IndexRune(valid, tok.next()) >= 0 {
		// keep consuming runes
		count++
	}
	tok.backup()
	return count
}

// acceptLetterRun consumes a run of alphabetic characters
func (tok *tokenizer) acceptLetterRun() (count int) {
	for unicode.IsLetter(tok.next()) {
		// keep consuming runes
		count++
	}
	tok.backup()
	return count
}

// acceptLetterRun consumes a run of alphabetic characters
func (tok *tokenizer) acceptAlphaNumRun() (count int) {
	next := tok.next()
	for unicode.IsLetter(next) || unicode.IsNumber(next) {
		// keep consuming runes
		next = tok.next()
		count++
	}
	tok.backup()
	return count
}

// run lexes the input by executing state functions until the state is nil
func (tok *tokenizer) run() {
	startState := tok.state
	for state := startState; state != nil; {
		state = state(tok)
	}
	tok.writer.Close()
}

func (tok *tokenizer) error(format string, args ...interface{}) stateFn {
	token := tokens.Token{TokenType: tokens.Error, Value: fmt.Sprintf(format, args...)}
	tok.writer.WriteToken(token)
	// stop the tokenizer
	return nil
}

// the state function expecting a number
func expectNumber(tok *tokenizer) stateFn {
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
	tok.emit(tokens.Number)

	// what is next?
	tok.ignoreSpaceRun()
	next := tok.peek()
	switch {
	case next == eofRune, next == '\n':
		return expectEOF
	case unicode.IsLetter(next) || unicode.IsNumber(next):
		return expectUnits
	default:
		return expectSymbol
	}
}

// the state function where an EOF is expected
func expectEOF(tok *tokenizer) stateFn {
	tok.ignoreSpaceRun()
	tok.ignoreRun('\n')
	if tok.next() == eofRune {
		tok.emit(tokens.EOF)
		return nil
	}
	return tok.error("expected EOF, but got '%s'", tok.current())
}

// the state function where a symbol is expected
func expectSymbol(tok *tokenizer) stateFn {
	for {
		switch next := tok.next(); {
		case next == '+':
			tok.emit(tokens.Plus)
			return expectNumber
		case next == '-':
			tok.emit(tokens.Minus)
			return expectNumber
		case next == '*':
			tok.emit(tokens.Multiply)
			return expectNumber
		case next == '/':
			tok.emit(tokens.Divide)
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

// the state function where units are expected
func expectUnits(tok *tokenizer) stateFn {
	tok.ignoreSpaceRun()
	count := tok.acceptAlphaNumRun()
	if count <= 0 {
		tok.next()
		return tok.error("expected units, but got '%s'", tok.current())
	}
	tok.emit(tokens.Units)

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

// the state function where 'in' is expected
func expectIn(tok *tokenizer) stateFn {
	tok.ignoreSpaceRun()
	if tok.accept("iI") && tok.accept("nN") && tok.accept(" ") {
		tok.emit(tokens.In)
		return expectUnits
	}
	// error
	tok.next()
	return tok.error("expected 'in' keyword, but got '%s'", tok.current())
}
