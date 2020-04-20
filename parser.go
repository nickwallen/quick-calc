package toks

import (
	"strconv"
)

// Parser Parses input text and outputs an expression.
type Parser struct {
	tokenizer   *Tokenizer
	expressions chan Expression
}

// NewParser Creates a new parser.
func NewParser(input string) *Parser {
	parser := &Parser{tokenizer: NewTokenizer(input)}
	return parser
}

// Parse Parse the input text and return an expression
func (parser *Parser) Parse() Expression {
	amount1 := parser.expectAmount()

	token := <-parser.tokenizer.Tokens()
	if token.TokenType != EOF {

		// expect an operation
		if token.TokenType == Plus {
			amount2 := parser.expectAmount()
			return Sum{left: amount1, right: amount2, units: amount1.Units}
		}
	}

	return amount1
}

func (parser *Parser) expectAmount() Amount {
	token := <-parser.tokenizer.Tokens()
	if token.TokenType != Number {
		// TODO do something
	}

	// TODO where to handle hexadecimal vs decimal?
	number, err := strconv.ParseInt(token.Value, 10, 64)
	if err != nil {
		// TODO error do something
	}

	units := <-parser.tokenizer.Tokens()
	if units.TokenType != Units {
		// TODO error do something
	}
	// TODO validate the units

	expression := Amount{Value: number, Units: AmountUnits{units: units.Value}}
	return expression
}
