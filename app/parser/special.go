package parser

import (
	"github.com/cappuccinotm/flangc/app/lexer"
	"errors"
	"fmt"
	"github.com/cappuccinotm/flangc/app/eval"
)

var errNoReservedKeyword = errors.New("word is not reserved")

func (p *Parser) findReservedKeyword(tkn lexer.Token) (expr eval.Expression, err error) {
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
func (p *Parser) parseFunc() (eval.Expression, error) {
	tkn, err := p.readAndValidateToken(lexer.Identifier)
	if err != nil {
		return nil, fmt.Errorf("expected identifier, got %s", tkn.Value)
	}

	result := &eval.Call{Name: "func", Args: []eval.Expression{&eval.Identifier{Name: tkn.Value}}}

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
func (p *Parser) parseLambda() (eval.Expression, error) {
	result := &eval.Call{Name: "lambda", Args: nil}

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
func (p *Parser) parseProg() (eval.Expression, error) {
	result := &eval.Call{Name: "prog", Args: nil}

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
