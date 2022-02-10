package cmd

import (
	"github.com/cappuccinotm/flangc/app/lexer"
	"os"
	"fmt"
	"github.com/cappuccinotm/flangc/app/parser"
	"log"
)

// Run command builds the program at the specified path.
type Run struct {
	FileLocation string `short:"f" long:"file" env:"FILE"`
}

// Execute runs the command.
func (b Run) Execute(_ []string) error {
	f, err := os.Open(b.FileLocation)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}

	lex := lexer.NewAdapter(f)
	p := parser.NewParser(lex)
	var code parser.Expression

	defer func() {
		if rvr := recover(); rvr != nil {
			log.Printf("[INFO] received one error: %v", rvr)
		}
	}()

	for ; err == nil; code, err = p.NextExpression() {
		log.Printf("[INFO] received code from parser %d", code)
	}

	if err != nil {
		return fmt.Errorf("parse next expression: %w", err)
	}

	log.Printf("[INFO] built without errors")

	return nil
}
