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
	"encoding/json"
)

// Run command builds the program at the specified path.
type Run struct {
	FileLocation string `short:"f" long:"file" env:"FILE"`
	FailOnError  bool   `short:"e" long:"error" env:"ERROR"`
	PrintAST     bool   `short:"a" long:"ast" env:"AST"`
	PrintJSONAST bool   `short:"j" long:"json-ast" env:"JSON_AST"`
}

// Execute runs the command.
func (b Run) Execute(_ []string) error {
	var rd io.Reader = os.Stdin
	if b.FileLocation != "" {
		f, err := os.Open(b.FileLocation)
		if err != nil {
			return fmt.Errorf("open file: %w", err)
		}
		defer f.Close()
		rd = f
	}

	lex := lexer.NewLexer(rd)

	p := parser.NewParser(lex)
	scope := eval.NewScope("", nil, false)

	for {
		if b.FileLocation == "" {
			fmt.Printf(">>> ")
		}
		expr, err := p.ParseNext()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			log.Printf("[WARN] parse: %v", err)
			if !b.FailOnError {
				continue
			}
			return fmt.Errorf("parse: %w", err)
		}

		if b.PrintAST {
			log.Printf("[INFO] ast: %v", expr)
		}
		if b.PrintJSONAST {
			bts, err := json.MarshalIndent(expr, "", "  ")
			if err != nil {
				log.Printf("[WARN] failed to marshal AST: %v", err)
			}
			log.Printf("[INFO] ast in json representation: \n%s", bts)
		}

		res, err := scope.Eval(expr)
		if err != nil {
			log.Printf("[WARN] execute statement %s: %v", expr.String(), err)
			if !b.FailOnError {
				continue
			}
			return fmt.Errorf("eval: %w", err)
		}

		if _, err = scope.Print(&eval.Call{Name: "print", Args: []eval.Expression{res}}); err != nil {
			log.Printf("[WARN] print result %s: %v", res.String(), err)
			if !b.FailOnError {
				continue
			}
			return fmt.Errorf("print result: %w", err)
		}
	}

	return nil
}
