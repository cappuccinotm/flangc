package parser

import (
	"github.com/cappuccinotm/flangc/app/lexer"
	"fmt"
	"log"
)

// Expression defines methods to evaluate them.
type Expression interface {
}

// Parser parses the expression through using the adapter.
type Parser struct {
	lexer  *lexAdapter
	parser parserParser
}

// NewParser makes new instance of Parser.
func NewParser(lexer *lexer.Adapter) *Parser {
	svc := &Parser{lexer: (*lexAdapter)(lexer)}
	svc.parser = parserNewParser()
	return svc
}

// NextExpression returns the next expression in the sequence
func (p *Parser) NextExpression() (Expression, error) {
	p.parser.Parse(p.lexer)
	return nil, nil
}

// ErrUnexpectedToken shows that the token, got from the lexer
// wasn't recognized by the parser.
type ErrUnexpectedToken lexer.Token

// Error returns the string representation of the error.
func (e ErrUnexpectedToken) Error() string {
	return fmt.Sprintf("unexpected token: %v", lexer.Token(e))
}

type lexAdapter lexer.Adapter

// Lex adapts the NextToken func to the parser's needs.
func (l *lexAdapter) Lex(lval *parserSymType) int {
	tkn, err := (*lexer.Adapter)(l).NextToken(lval.yys)
	if err != nil {
		log.Fatalf("[ERROR] failed to parse token: %v", err)
		return 0
	}
	switch tkn.Kind {
	case lexer.Less:
		return LESS
	case lexer.Greater:
		return GREATER
	case lexer.Equal:
		return EQUAL
	case lexer.NotEqual:
		return NOT_EQUAL
	case lexer.Plus:
		return PLUS
	case lexer.Minus:
		return MINUS
	case lexer.Multiply:
		return MULTIPLY
	case lexer.Divide:
		return DIVIDE
	case lexer.Number:
		return NUMBER
	case lexer.Null:
		return NULL
	case lexer.Identifier:
		return IDENTIFIER
	case lexer.LBrace:
		return LBRACE
	case lexer.RBrace:
		return RBRACE
	default:
		return 0
	}
}

func (l *lexAdapter) Error(s string) {}
