package toks

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Tokenizer A tokenizer performs lexical analysis on an input string.
type Tokenizer struct {
	state  stateFn    // the current state function
	input  string     // the string to scan
	start  int        // start position for this item
	pos    int        // current position in the input
	width  int        // width of the last rune read
	tokens chan Token // the channel where items are emitted
}

// the state of the scanner as a function that returns the next state.
type stateFn func(*Tokenizer) stateFn

// New Creates a tokenizer that can tokenize a string.
func New(input string) *Tokenizer {
	tok := &Tokenizer{
		state:  expectNumber,
		input:  input,
		tokens: make(chan Token),
	}
	go tok.run()
	return tok
}

// NextToken fetch the next token.
func (tok *Tokenizer) NextToken() Token {
	for {
		select {
		case token := <-tok.tokens:
			return token
		default:
			tok.state = tok.state(tok)
		}
	}
}

// Tokens returns a channel that will be filled with tokens.
func (tok *Tokenizer) Tokens() chan Token {
	return tok.tokens
}

// returns what is currently being scanned
func (tok *Tokenizer) current() string {
	return tok.input[tok.start:tok.pos]
}

// emits a token to be consumed by the client
func (tok *Tokenizer) emit(tokenType TokenType) {
	var token Token
	switch tokenType {
	case EOFToken:
		token = EOFToken.Of("")
	default:
		token = tokenType.Of(tok.current())
	}
	tok.tokens <- token
	tok.start = tok.pos
}

// skips over the pending input
func (tok *Tokenizer) ignore() {
	tok.start = tok.pos
}

// skips over a run of values
func (tok *Tokenizer) ignoreRun(ignore rune) {
	for tok.next() == ignore {
		tok.ignore()
	}
}

// skips over any whitespace
func (tok *Tokenizer) ignoreWhitespaces() {
	for unicode.IsSpace(tok.next()) {
		tok.ignore()
	}
	tok.backup()
}

// steps back one
func (tok *Tokenizer) backup() {
	tok.pos -= tok.width
}

// peek returns, but does not consume the next rune in the input.
func (tok *Tokenizer) peek() rune {
	rune := tok.next()
	tok.backup()
	return rune
}

func (tok *Tokenizer) next() rune {
	if tok.pos >= len(tok.input) {
		tok.width = 0
		return -1 // eof
	}
	var r rune
	r, tok.width = utf8.DecodeRuneInString(tok.input[tok.pos:])
	tok.pos += tok.width
	return r
}

// accept consumes the next rune if it is valid.
func (tok *Tokenizer) accept(valid string) bool {
	if strings.IndexRune(valid, tok.next()) >= 0 {
		return true
	}
	tok.backup()
	return false
}

// acceptRun consumes a run of strings
func (tok *Tokenizer) acceptRun(valid string) {
	for strings.IndexRune(valid, tok.next()) >= 0 {
		// keep consuming runes
	}
	tok.backup()
}

// run lexes the input by executing state functions until the state is nil
func (tok *Tokenizer) run() {
	startState := expectNumber
	for state := startState; state != nil; {
		state = state(tok)
	}
	close(tok.tokens)
}

func (tok *Tokenizer) errorf(format string, args ...interface{}) stateFn {
	tok.tokens <- Token{typ: ErrorToken, val: fmt.Sprintf(format, args...)}
	// stop the tokenizer
	return nil
}
