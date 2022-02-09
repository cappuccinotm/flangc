package cmd

import (
	"github.com/cappuccinotm/flangc/app/lexer"
	"os"
	"fmt"
	"github.com/cappuccinotm/flangc/app/parser"
	"errors"
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

	err = errors.New("no errors")
	for ; err != nil; _, err = p.NextExpression() {
	}

	if err != nil {
		return fmt.Errorf("parse next expression: %w", err)
	}

	fmt.Println("built without errors")

	return nil
}
