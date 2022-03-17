package parser

import (
	"github.com/cappuccinotm/flangc/app/lexer"
	"fmt"
	"strconv"
	"errors"
	"github.com/cappuccinotm/flangc/app/eval"
)

// Parser parsers expressions, read through the lexer
type Parser struct {
	l *lexer.Lexer
}

// NewParser creates a new parser
func NewParser(l *lexer.Lexer) *Parser {
	return &Parser{l}
}

// ParseNext parses next expression.
func (p *Parser) ParseNext() (eval.Expression, error) {
	cursor := p.l.Cursor()
	tkn, err := p.l.NextToken()
	if err != nil {
		return nil, fmt.Errorf("get next token: %w", err)
	}

	switch tkn.Type {
	case lexer.SQuote:
		expr, err := p.parseTuple()
		if err != nil {
			return nil, fmt.Errorf("parse list: %w", err)
		}
		return expr, nil
	case lexer.LParen:
		cursor = p.l.Cursor()
		if tkn, err = p.readAndValidateToken(lexer.Identifier); err != nil {
			return nil, fmt.Errorf("get next token at %s: %w", cursor, err)
		}

		expr, err := p.parseCall(tkn)
		if err != nil {
			return nil, fmt.Errorf("parse call at %s: %w", cursor, err)
		}
		return expr, nil
	default:
		return nil, fmt.Errorf("unexpected token at %s: %s", cursor, tkn)
	}
}

// parses (el1 el2 el3) as list, without counting quote sign '
func (p *Parser) parseTuple() (eval.Expression, error) {
	var exprs []eval.Expression

	tkn, err := p.readAndValidateToken(lexer.LParen)
	if err != nil {
		return nil, err
	}

	for {
		if tkn, err = p.l.NextToken(); err != nil {
			return nil, fmt.Errorf("get next token: %w", err)
		}

		switch tkn.Type {
		case lexer.Number:
			f, err := strconv.ParseFloat(tkn.Value, 64)
			if err != nil {
				return nil, fmt.Errorf("parse number: %w", err)
			}
			exprs = append(exprs, &eval.Number{Value: f})
		case lexer.Identifier:
			exprs = append(exprs, parseIdentifier(tkn.Value))
		case lexer.RParen:
			return &eval.List{Values: exprs}, nil
		case lexer.LParen, lexer.SQuote:
			p.l.UnreadToken()
			expr, err := p.ParseNext()
			if err != nil {
				return nil, fmt.Errorf("parse expression: %w", err)
			}
			exprs = append(exprs, expr)
		default:
			return nil, fmt.Errorf("unexpected token: %s", tkn)
		}
	}
}

func parseIdentifier(value string) eval.Expression {
	switch value {
	case "true":
		return &eval.Boolean{Value: true}
	case "false":
		return &eval.Boolean{Value: false}
	case "null":
		return eval.Null{}
	default:
		return &eval.Identifier{Name: value}
	}
}

func (p *Parser) parseCall(tkn lexer.Token) (eval.Expression, error) {
	expr, err := p.findReservedKeyword(tkn)
	switch {
	case errors.Is(err, errNoReservedKeyword):
	case err != nil:
		return nil, err
	case err == nil:
		if _, ok := expr.(*eval.Call); !ok {
			return nil, fmt.Errorf("expected function call, got %s", expr.String())
		}
		return expr, nil
	}

	result := &eval.Call{Name: tkn.Value}

	for {
		tkn, err := p.l.NextToken()
		if err != nil {
			return nil, fmt.Errorf("get next token: %w", err)
		}

		switch tkn.Type {
		case lexer.Identifier:
			result.Args = append(result.Args, parseIdentifier(tkn.Value))
		case lexer.Number:
			f, err := strconv.ParseFloat(tkn.Value, 64)
			if err != nil {
				return nil, fmt.Errorf("parse number: %w", err)
			}
			result.Args = append(result.Args, &eval.Number{Value: f})
		case lexer.RParen:
			return result, nil
		case lexer.LParen, lexer.SQuote:
			p.l.UnreadToken()
			expr, err := p.ParseNext()
			if err != nil {
				return nil, fmt.Errorf("parse expression: %w", err)
			}
			result.Args = append(result.Args, expr)
		}
	}
}

func (p *Parser) readAndValidateToken(typ lexer.TokenType) (lexer.Token, error) {
	tkn, err := p.l.NextToken()
	if err != nil {
		return lexer.Token{}, fmt.Errorf("get next token: %w", err)
	}

	if tkn.Type != typ {
		return lexer.Token{}, fmt.Errorf("expected %s, got: %s", typ, tkn)
	}

	return tkn, nil
}
