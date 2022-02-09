package lexer

import (
	"fmt"
	"strconv"
)

// TokenKind represents the type of the token.
type TokenKind int

// Supported token kinds.
const (
	Less TokenKind = iota
	Greater
	Equal
	NotEqual
	Plus
	Minus
	Multiply
	Divide
	QuoteSign
	LBrace
	RBrace
	Number
	Null
	Identifier
)

// Token represents a basic token.
type Token struct {
	Kind  TokenKind
	Value interface{}
}

// Adapter adapts the generated lexer's tokens into the
// internal token representation.
type Adapter struct {
	lexer         *Lexer
	nextTokenFunc func(*Token, *error)
}

// NextToken returns the next token in the source code.
func (a *Adapter) NextToken(lvl int) (Token, error) {
	var token Token
	var err error
	switch a.lexer.next(lvl) {
	case 0:
		token.Kind = Less
	case 1:
		token.Kind = Greater
	case 2:
		token.Kind = Equal
	case 3:
		token.Kind = NotEqual
	case 4:
		token.Kind = Plus
	case 5:
		token.Kind = Minus
	case 6:
		token.Kind = Multiply
	case 7:
		token.Kind = Divide
	case 8:
		token.Kind = LBrace
	case 9:
		token.Kind = RBrace
	case 10:
		token.Kind = Number
		token.Value, err = strconv.ParseFloat(a.lexer.Text(), 64)
	case 11:
		token.Kind = Null
	case 12:
		token.Kind = Identifier
		token.Value = a.lexer.Text()
	case 13: /* eat up whitespace */
	case 14:
		err = ErrUnrecognizedCharacter(a.lexer.Text())
	case 15: /* eat up one-line comments */
	}
	if err != nil {
		return Token{}, fmt.Errorf("unrecognized token %q at %d: %w", a.lexer.Text(), a.lexer.Line(), err)
	}
	return token, nil
}

// ErrUnrecognizedCharacter represents a parsing error.
type ErrUnrecognizedCharacter string

// Error returns a string representation of the error.
func (e ErrUnrecognizedCharacter) Error() string {
	return fmt.Sprintf("unrecognized character sequence: %q", string(e))
}
