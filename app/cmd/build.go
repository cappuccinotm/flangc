package cmd

import (
	"github.com/cappuccinotm/flangc/app/lexer"
	"os"
	"fmt"
	"log"
	"github.com/cappuccinotm/flangc/app/parser"
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

	lex := lexer.NewLexer(f)
	//var tkn lexer.Token
	//for err == nil {
	//	if tkn, err = lex.NextToken(); errors.Is(err, io.EOF) {
	//		err = nil
	//		break
	//	} else if err != nil {
	//		return fmt.Errorf("next token: %w", err)
	//	}
	//	log.Printf("[INFO] received token: %s", tkn)
	//}

	expr, err := parser.New(lex).Parse()
	if err != nil {
		return fmt.Errorf("parse: %w", err)
	}

	log.Printf("[INFO] parsed expression: %s", expr)

	if expr, err = parser.New(lex).Parse(); err != nil {
		return fmt.Errorf("parse: %w", err)
	}

	log.Printf("[INFO] parsed expression: %s", expr)

	log.Printf("[INFO] built without errors")

	return nil
}
