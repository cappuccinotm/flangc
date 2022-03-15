package lexer

import (
	"bufio"
	"fmt"
	"io"
)

// Lexer reads tokens from an input stream.
type Lexer struct {
	rd           *bufio.Reader
	currentLine  int
	readComments bool
	lastToken    Token
	useLastToken bool
}

// NewLexer creates a new Lexer.
func NewLexer(rd io.Reader) *Lexer {
	return &Lexer{
		rd:          bufio.NewReader(rd),
		currentLine: 1,
	}
}

// NextToken returns the next token from the input stream.
func (l *Lexer) NextToken() (Token, error) {
	if l.useLastToken {
		l.useLastToken = false
		return l.lastToken, nil
	}

	var (
		err error
		r   = ' '
	)

	for r == ' ' || r == '\n' || r == '\t' {
		if r, _, err = l.rd.ReadRune(); err != nil {
			return Token{}, fmt.Errorf("read next symbol: %w", err)
		}
		if r == '\n' {
			l.currentLine++
		}
	}

	var tkn Token

	switch {
	case r == '\'':
		tkn = Token{Type: SQuote}
	case r == '(':
		tkn = Token{Type: LParen}
	case r == ')':
		tkn = Token{Type: RParen}
	case isDigit(r):
		tkn = l.readNumber(r)
	case r == '_', isLetter(r):
		tkn = l.readIdentifier(r)
	case r == '/':
		if l.readComments {
			tkn = l.readComment(r)
		}
		l.currentLine++
		return l.NextToken()
	default:
		return Token{}, fmt.Errorf("unexpected symbol: %c", r)
	}

	l.lastToken = tkn

	return tkn, nil
}

func (l *Lexer) UnreadToken() {
	if l.useLastToken {
		panic("unread token twice")
	}
	l.useLastToken = true
}

func (l *Lexer) readIdentifier(r rune) Token {
	var sb = &[]rune{}
	*sb = append(*sb, r)

	for {
		r, _, err := l.rd.ReadRune()
		if err != nil {
			return Token{Type: Identifier, Value: string(*sb)}
		}

		if !isLetter(r) && !isDigit(r) && r != '_' {
			l.rd.UnreadRune()
			return Token{Type: Identifier, Value: string(*sb)}
		}

		*sb = append(*sb, r)
	}
}

func (l *Lexer) readNumber(r rune) Token {
	var sb = &[]rune{}
	*sb = append(*sb, r)

	for {
		r, _, err := l.rd.ReadRune()
		if err != nil {
			return Token{Type: Number, Value: string(*sb)}
		}

		if !isDigit(r) && r != '.' {
			l.rd.UnreadRune()
			return Token{Type: Number, Value: string(*sb)}
		}

		*sb = append(*sb, r)
	}
}

func (l *Lexer) readComment(r rune) Token {
	var sb = &[]rune{}
	*sb = append(*sb, r)

	for {
		r, _, err := l.rd.ReadRune()
		if err != nil || r == '\n' {
			return Token{Type: Comment, Value: string(*sb)}
		}

		*sb = append(*sb, r)
	}
}

func (l *Lexer) Line() int {
	return l.currentLine
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func isLetter(r rune) bool {
	return r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z'
}
