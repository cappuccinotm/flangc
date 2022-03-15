package lexer

type TokenType string

const (
	SQuote     TokenType = "SQuote"
	LParen     TokenType = "LParen"
	RParen     TokenType = "RParen"
	Number     TokenType = "Number"
	Identifier TokenType = "Identifier"
	Comment    TokenType = "Comment"
)

type Token struct {
	Type  TokenType
	Value string
}
