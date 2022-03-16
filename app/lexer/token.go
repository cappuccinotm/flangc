package lexer

type TokenType string

const (
	SQuote     TokenType = "quote"
	LParen     TokenType = "("
	RParen     TokenType = ")"
	Number     TokenType = "number"
	Identifier TokenType = "identifier"
	Comment    TokenType = "comment"
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
