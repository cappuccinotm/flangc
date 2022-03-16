package parser

import (
	"github.com/cappuccinotm/flangc/app/lexer"
	"errors"
	"fmt"
)

var errNoReservedKeyword = errors.New("word is not reserved")

func (p *Parser) findReservedKeyword(tkn lexer.Token) (expr Expression, err error) {
	switch tkn.Value {
	case "func":
		expr, err = p.parseFunc()
	case "lambda":
		expr, err = p.parseLambda()
	case "prog":
		expr, err = p.parseProg()
	default:
		return nil, errNoReservedKeyword
	}
	if err != nil {
		return nil, fmt.Errorf("scan %s: %s", tkn.Value, err)
	}

	return expr, nil
}

// (func name (args) (body))
func (p *Parser) parseFunc() (Expression, error) {
	tkn, err := p.readAndValidateToken(lexer.Identifier)
	if err != nil {
		return nil, fmt.Errorf("expected identifier, got %s", tkn.Value)
	}

	result := &Call{Name: "func", Args: []Expression{&Identifier{Name: tkn.Value}}}

	expr, err := p.parseTuple()
	if err != nil {
		return nil, err
	}

	result.Args = append(result.Args, expr)

	if _, err = p.readAndValidateToken(lexer.LParen); err != nil {
		return nil, err
	}

	if tkn, err = p.readAndValidateToken(lexer.Identifier); err != nil {
		return nil, err
	}

	if expr, err = p.parseCall(tkn); err != nil {
		return nil, err
	}

	result.Args = append(result.Args, expr)

	if _, err = p.readAndValidateToken(lexer.RParen); err != nil {
		return nil, err
	}

	return result, nil
}

// (lambda (args) (body))
func (p *Parser) parseLambda() (Expression, error) {
	result := &Call{Name: "lambda", Args: nil}

	expr, err := p.parseTuple()
	if err != nil {
		return nil, err
	}

	result.Args = append(result.Args, expr)

	if _, err = p.readAndValidateToken(lexer.LParen); err != nil {
		return nil, err
	}

	tkn, err := p.readAndValidateToken(lexer.Identifier)
	if err != nil {
		return nil, err
	}

	if expr, err = p.parseCall(tkn); err != nil {
		return nil, err
	}

	result.Args = append(result.Args, expr)

	if _, err = p.readAndValidateToken(lexer.RParen); err != nil {
		return nil, err
	}

	return result, nil
}

// (prog (scopevars) (statements))
func (p *Parser) parseProg() (Expression, error) {
	result := &Call{Name: "prog", Args: nil}

	expr, err := p.parseTuple()
	if err != nil {
		return nil, err
	}

	result.Args = append(result.Args, expr)

	if expr, err = p.parseTuple(); err != nil {
		return nil, err
	}

	result.Args = append(result.Args, expr)

	if _, err = p.readAndValidateToken(lexer.RParen); err != nil {
		return nil, err
	}

	return result, nil
}
