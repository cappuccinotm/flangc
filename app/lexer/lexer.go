package lexer

import (
	"fmt"
	"strconv"
)

// TokenKind represents the type of the token.
type TokenKind int

// Supported token kinds.
const (
	undefined TokenKind = -100
	SQuote    TokenKind = -1 + iota
	Quote
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
	token := Token{Kind: undefined, Value: nil}
	var err error
	for token.Kind == undefined && err == nil {
		switch TokenKind(a.lexer.next(lvl)) {
		case SQuote:
			token.Kind = SQuote
		case Quote:
			token.Kind = Quote
		case LBrace:
			token.Kind = LBrace
		case RBrace:
			token.Kind = RBrace
		case Number:
			token.Kind = Number
			token.Value, err = strconv.ParseFloat(a.lexer.Text(), 64)
		//case 4:
		//token.Kind = Null
		case Identifier:
			token.Kind = Identifier
			token.Value = a.lexer.Text()
		}
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
