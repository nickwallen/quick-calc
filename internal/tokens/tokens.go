package tokens

import (
	"fmt"
)

// Token a tokenizer emits tokens.
type Token struct {
	TokenType TokenType // the type like numberToken
	Value     string    // the value, like "46.2"
	Position  int       // the starting position of the token
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

// Token creates a new Token.
func (t TokenType) Token(value string) Token {
	return Token{TokenType: t, Value: value}
}

// TokenAt creates a Token at a fixed position.
func (t TokenType) TokenAt(value string, position int) Token {
	return Token{TokenType: t, Value: value, Position: position}
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
