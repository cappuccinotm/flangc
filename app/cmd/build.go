package cmd

import (
	"github.com/cappuccinotm/flangc/app/lexer"
	"os"
	"fmt"
	"log"
	"github.com/cappuccinotm/flangc/app/parser"
	"io"
	"errors"
	"github.com/cappuccinotm/flangc/app/eval"
)

// Run command builds the program at the specified path.
type Run struct {
	FileLocation string `short:"f" long:"file" env:"FILE"`
	FailOnError  bool   `short:"e" long:"error" env:"ERROR"`
}

// Execute runs the command.
func (b Run) Execute(_ []string) error {
	f, err := os.Open(b.FileLocation)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}

	lex := lexer.NewLexer(f)

	p := parser.NewParser(lex)
	scope := eval.NewScope(nil, nil, false)

	for {
		expr, err := p.ParseNext()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			log.Printf("[ERROR] parse: %v", err)
			if !b.FailOnError {
				continue
			}
			return fmt.Errorf("parse: %w", err)
		}

		res, err := scope.Eval(expr)
		if err != nil {
			log.Printf("[ERROR] execute statement %s: %v", expr.String(), err)
			if !b.FailOnError {
				continue
			}
			return fmt.Errorf("eval: %w", err)
		}
		_, err = scope.Print(&eval.Call{Name: "print", Args: []eval.Expression{res}})
		if err != nil {
			log.Printf("[ERROR] print result %s: %v", res.String(), err)
			if !b.FailOnError {
				continue
			}
			return fmt.Errorf("print result: %w", err)
		}
		//log.Printf("[INFO] parsed expression: %s", expr)
	}

	//log.Printf("[INFO] built without errors")

	return nil
}
