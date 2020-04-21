package tokenizer

import (
	"fmt"
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
	WriteToken(token Token) error
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
func (tok *tokenizer) emit(tokenType TokenType) {
	var token Token
	switch tokenType {
	case EOF:
		token = EOF.Token("")
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

// run lexes the input by executing state functions until the state is nil
func (tok *tokenizer) run() {
	startState := tok.state
	for state := startState; state != nil; {
		state = state(tok)
	}
	tok.writer.Close()
}

func (tok *tokenizer) error(format string, args ...interface{}) stateFn {
	token := Token{TokenType: Error, Value: fmt.Sprintf(format, args...)}
	tok.writer.WriteToken(token)
	// stop the tokenizer
	return nil
}
