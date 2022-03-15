package lexer

import (
	"bufio"
	"fmt"
	"io"
)

// Lexer reads tokens from an input stream.
type Lexer struct {
	rd           *bufio.Reader
	cursor       Cursor
	readComments bool
	lastToken    struct {
		value Token
		use   bool
	}
}

// NewLexer creates a new Lexer.
func NewLexer(rd io.Reader) *Lexer {
	return &Lexer{
		rd:     bufio.NewReader(rd),
		cursor: Cursor{Line: 1, Col: 1},
	}
}

// NextToken returns the next token from the input stream.
func (l *Lexer) NextToken() (Token, error) {
	if l.lastToken.use {
		l.lastToken.use = false
		return l.lastToken.value, nil
	}

	var (
		err error
		r   = ' '
	)

	for r == ' ' || r == '\n' || r == '\t' {
		if r, _, err = l.readRune(); err != nil {
			return Token{}, fmt.Errorf("read next symbol: %w", err)
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
		tkn = l.readComment(r)
		if !l.readComments {
			if tkn, err = l.NextToken(); err != nil {
				return Token{}, err
			}
		}
	default:
		return Token{}, fmt.Errorf("unexpected symbol: %c", r)
	}

	l.lastToken.value = tkn

	return tkn, nil
}

func (l *Lexer) UnreadToken() {
	if l.lastToken.use {
		panic("unread token twice")
	}
	l.lastToken.use = true
}

func (l *Lexer) readIdentifier(r rune) Token {
	var sb = &[]rune{}
	*sb = append(*sb, r)

	for {
		r, _, err := l.readRune()
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
		r, _, err := l.readRune()
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
		r, _, err := l.readRune()
		if err != nil || r == '\n' {
			return Token{Type: Comment, Value: string(*sb)}
		}

		*sb = append(*sb, r)
	}
}

func (l *Lexer) Cursor() Cursor {
	return l.cursor
}

func (l *Lexer) readRune() (r rune, size int, err error) {
	if r, size, err = l.rd.ReadRune(); err != nil {
		return
	}

	l.cursor.Col++

	if r == '\n' {
		l.cursor.Col = 0
		l.cursor.Line++
	}

	return
}

type Cursor struct {
	Line, Col int
}

func (c Cursor) String() string {
	return fmt.Sprintf("%d:%d", c.Line, c.Col)
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func isLetter(r rune) bool {
	return r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z'
}
