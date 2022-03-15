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
		return p.parseList()
	case lexer.LParen:
		return p.parseCall()
	default:
		return nil, fmt.Errorf("unexpected token at line %d: %s", p.l.Line(), tkn.Type)
	}
}

func (p *Parser) parseList() (Expression, error) {
	var exprs []Expression

	for {
		tkn, err := p.l.NextToken()
		if err != nil {
			return nil, fmt.Errorf("get next token: %w", err)
		}

		switch tkn.Type {
		case lexer.Identifier:
			exprs = append(exprs, Identifier{Name: tkn.Value})
		case lexer.Number:
			f, err := strconv.ParseFloat(tkn.Value, 64)
			if err != nil {
				return nil, fmt.Errorf("parse number at line %d: %w", p.l.Line(), err)
			}
			exprs = append(exprs, Number{Value: f})
		case lexer.RParen:
			return List{Elements: exprs}, nil
		case lexer.LParen, lexer.SQuote:
			p.l.UnreadToken()
			expr, err := p.Parse()
			if err != nil {
				return nil, fmt.Errorf("parse expression at line %d: %w", p.l.Line(), err)
			}
			exprs = append(exprs, expr)
		default:
			return nil, fmt.Errorf("unexpected token at line %d: %s", p.l.Line(), tkn.Type)
		}
	}
}

func (p *Parser) parseCall() (Expression, error) {
	tkn, err := p.l.NextToken()
	if err != nil {
		return nil, fmt.Errorf("get next token: %w", err)
	}

	if tkn.Type != lexer.Identifier {
		return nil, fmt.Errorf("unexpected token at line %d: %s", p.l.Line(), tkn.Type)
	}

	result := Call{Name: tkn.Value}

	for {
		tkn, err := p.l.NextToken()
		if err != nil {
			return nil, fmt.Errorf("get next token: %w", err)
		}

		switch tkn.Type {
		case lexer.Identifier:
			result.Args = append(result.Args, Identifier{Name: tkn.Value})
		case lexer.Number:
			f, err := strconv.ParseFloat(tkn.Value, 64)
			if err != nil {
				return nil, fmt.Errorf("parse number at line %d: %w", p.l.Line(), err)
			}
			result.Args = append(result.Args, Number{Value: f})
		case lexer.RParen:
			return result, nil
		case lexer.LParen, lexer.SQuote:
			p.l.UnreadToken()
			expr, err := p.Parse()
			if err != nil {
				return nil, fmt.Errorf("parse expression at line %d: %w", p.l.Line(), err)
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
