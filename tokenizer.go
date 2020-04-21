package toks

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	eofRune = rune(0)
)

// Tokenizer A Tokenizer performs lexical analysis on an input string.
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

// Tokenize Tokenizes the input string and writes each token to the output channel.
func Tokenize(input string, output chan Token) {
	tok := &Tokenizer{
		state:  expectNumber,
		input:  input,
		tokens: output,
	}
	tok.run()
}

// TokenizeToSlice Tokenizes the input string and returns each token as a slice.
func TokenizeToSlice(input string) []Token {
	// tokenize in a separate goroutine
	output := make(chan Token)
	go Tokenize(input, output)

	// fetch the tokens into a slice
	tokens := make([]Token, 0)
	for token := range output {
		tokens = append(tokens, token) // TODO probably a more efficient way to do this
	}
	return tokens
}

// returns what is currently being scanned
func (tok *Tokenizer) current() string {
	// remove leading, trailing, and embedded whitespace
	value := tok.input[tok.start:tok.pos]
	return strings.ReplaceAll(value, " ", "")
}

// emits a Token to be consumed by the client
func (tok *Tokenizer) emit(tokenType TokenType) {
	var token Token
	switch tokenType {
	case EOF:
		token = EOF.token("")
	default:
		token = tokenType.token(tok.current())
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
func (tok *Tokenizer) ignoreSpaceRun() {
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
		return eofRune
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
func (tok *Tokenizer) acceptRun(valid string) (count int) {
	for strings.IndexRune(valid, tok.next()) >= 0 {
		// keep consuming runes
		count++
	}
	tok.backup()
	return count
}

// acceptLetterRun consumes a run of alphabetic characters
func (tok *Tokenizer) acceptLetterRun() (count int) {
	for unicode.IsLetter(tok.next()) {
		// keep consuming runes
		count++
	}
	tok.backup()
	return count
}

// run lexes the input by executing state functions until the state is nil
func (tok *Tokenizer) run() {
	startState := tok.state
	for state := startState; state != nil; {
		state = state(tok)
	}
	close(tok.tokens)
}

func (tok *Tokenizer) error(format string, args ...interface{}) stateFn {
	tok.tokens <- Token{TokenType: Error, Value: fmt.Sprintf(format, args...)}
	// stop the Tokenizer
	return nil
}
