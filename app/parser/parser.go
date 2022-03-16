package parser

import (
	"github.com/cappuccinotm/flangc/app/lexer"
	"fmt"
	"strconv"
	"errors"
)

type Parser struct {
	l *lexer.Lexer
}

func New(l *lexer.Lexer) *Parser {
	return &Parser{l}
}

func (p *Parser) Parse() (Expression, error) {
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
func (p *Parser) parseTuple() (Expression, error) {
	var exprs []Expression

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
			exprs = append(exprs, &Number{Value: f})
		case lexer.Identifier:
			exprs = append(exprs, &Identifier{Name: tkn.Value})
		case lexer.RParen:
			return &List{Elements: exprs}, nil
		case lexer.LParen, lexer.SQuote:
			p.l.UnreadToken()
			expr, err := p.Parse()
			if err != nil {
				return nil, fmt.Errorf("parse expression: %w", err)
			}
			exprs = append(exprs, expr)
		default:
			return nil, fmt.Errorf("unexpected token: %s", tkn)
		}
	}
}

func (p *Parser) parseCall(tkn lexer.Token) (Expression, error) {
	expr, err := p.findReservedKeyword(tkn)
	switch {
	case errors.Is(err, errNoReservedKeyword):
	case err != nil:
		return nil, err
	case err == nil:
		if _, ok := expr.(*Call); !ok {
			return nil, fmt.Errorf("expected function call, got %s", expr.String())
		}
		return expr, nil
	}

	result := &Call{Name: tkn.Value}

	for {
		tkn, err := p.l.NextToken()
		if err != nil {
			return nil, fmt.Errorf("get next token: %w", err)
		}

		switch tkn.Type {
		case lexer.Identifier:
			result.Args = append(result.Args, &Identifier{Name: tkn.Value})
		case lexer.Number:
			f, err := strconv.ParseFloat(tkn.Value, 64)
			if err != nil {
				return nil, fmt.Errorf("parse number: %w", err)
			}
			result.Args = append(result.Args, &Number{Value: f})
		case lexer.RParen:
			return result, nil
		case lexer.LParen, lexer.SQuote:
			p.l.UnreadToken()
			expr, err := p.Parse()
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
