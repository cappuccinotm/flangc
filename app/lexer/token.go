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

func (t Token) String() string {
	if t.Value == "" {
		return string(t.Type)
	}
	return "(" + string(t.Type) + ": " + t.Value + ")"
}
