package parser

import (
	"github.com/cappuccinotm/flangc/app/lexer"
	"fmt"
	"strconv"
	"strings"
)

type Parser struct {
	l *lexer.Lexer
}

func New(l *lexer.Lexer) *Parser {
	return &Parser{l}
}

func (p *Parser) Parse() (Expression, error) {
	tkn, err := p.l.NextToken()
	if err != nil {
		return nil, fmt.Errorf("get next token: %w", err)
	}

	switch tkn.Type {
	case lexer.SQuote:
		expr, err := p.parseList()
		if err != nil {
			return nil, fmt.Errorf("parse list: %w", err)
		}
		return expr, nil
	case lexer.LParen:
		expr, err := p.parseCall()
		if err != nil {
			return nil, fmt.Errorf("parse call: %w", err)
		}
		return expr, nil
	default:
		return nil, fmt.Errorf("unexpected token at %s: %s", p.l.Cursor(), tkn)
	}
}

func (p *Parser) parseList() (Expression, error) {
	var exprs []Expression

	tkn, err := p.l.NextToken()
	if err != nil {
		return nil, fmt.Errorf("get next token: %w", err)
	}

	if tkn.Type != lexer.LParen {
		return nil, fmt.Errorf("unexpected token at %s: %s", p.l.Cursor(), tkn)
	}

	for {
		tkn, err := p.l.NextToken()
		if err != nil {
			return nil, fmt.Errorf("get next token: %w", err)
		}

		line := p.l.Cursor().Line
		switch tkn.Type {
		case lexer.Number:
			f, err := strconv.ParseFloat(tkn.Value, 64)
			if err != nil {
				return nil, fmt.Errorf("parse number at line %d: %w", line, err)
			}
			exprs = append(exprs, Number{Value: f})
		case lexer.Identifier:
			exprs = append(exprs, Identifier{Name: tkn.Value})
		case lexer.RParen:
			return List{Elements: exprs}, nil
		case lexer.LParen, lexer.SQuote:
			p.l.UnreadToken()
			expr, err := p.Parse()
			if err != nil {
				return nil, fmt.Errorf("parse expression at line %d: %w", line, err)
			}
			exprs = append(exprs, expr)
		default:
			return nil, fmt.Errorf("unexpected token at %s: %s", p.l.Cursor(), tkn)
		}
	}
}

func (p *Parser) parseCall() (Expression, error) {
	tkn, err := p.l.NextToken()
	if err != nil {
		return nil, fmt.Errorf("get next token: %w", err)
	}

	if tkn.Type != lexer.Identifier {
		return nil, fmt.Errorf("unexpected token at %s: %s", p.l.Cursor(), tkn)
	}

	result := Call{Name: tkn.Value}

	for {
		tkn, err := p.l.NextToken()
		if err != nil {
			return nil, fmt.Errorf("get next token: %w", err)
		}

		line := p.l.Cursor().Line
		switch tkn.Type {
		case lexer.Identifier:
			result.Args = append(result.Args, Identifier{Name: tkn.Value})
		case lexer.Number:
			f, err := strconv.ParseFloat(tkn.Value, 64)
			if err != nil {
				return nil, fmt.Errorf("parse number at line %d: %w", line, err)
			}
			result.Args = append(result.Args, Number{Value: f})
		case lexer.RParen:
			return result, nil
		case lexer.LParen, lexer.SQuote:
			p.l.UnreadToken()
			expr, err := p.Parse()
			if err != nil {
				return nil, fmt.Errorf("parse expression at line %d: %w", line, err)
			}
			result.Args = append(result.Args, expr)
		}
	}
}

type Expression interface {
	String() string
}

type Call struct {
	Name string
	Args []Expression
}

func (c Call) String() string {
	var args []string
	for _, arg := range c.Args {
		args = append(args, arg.String())
	}
	return fmt.Sprintf("%s(%s)", c.Name, strings.Join(args, ", "))
}

type Identifier struct {
	Name string
}

func (i Identifier) String() string { return i.Name }

type List struct {
	Elements []Expression
}

func (l List) String() string {
	return fmt.Sprintf("%v", l.Elements)
}

type Number struct {
	Value float64
}

func (n Number) String() string {
	return fmt.Sprintf("%f", n.Value)
}
