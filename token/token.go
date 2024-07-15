package token

type TokenType uint8

type Token struct {
	Type    TokenType
	Literal string
}

// Token types
const (
	// Special
	ILLEGAL TokenType = iota
	EOF

	// Identifier and literals
	IDENTIFIER
	INT

	// Operators
	ASSIGN
	PLUS
	MINUS
	BANG
	ASTERISK
	SLASH

	LT
	GT

	// Delimiters
	COMMA
	SEMICOLON

	LPAREN
	RPAREN
	LBRACE
	RBRACE

	// Keywords
	FUNCTION
	LET
	IF
	ELSE
	RETURN
	TRUE
	FALSE
	EQ
	NOT_EQ
)

// For pretty printing the enum values
var tokenTypeStrings = [...]string{
	"ILLEGAL",
	"EOF",
	"IDENTIFIER",
	"INT",
	"ASSIGN",
	"PLUS",
	"MINUS",
	"BANG",
	"ASTERISK",
	"SLASH",
	"LT",
	"GT",
	"COMMA",
	"SEMICOLON",
	"LPAREN",
	"RPAREN",
	"LBRACE",
	"RBRACE",
	"FUNCTION",
	"LET",
	"IF",
	"ELSE",
	"RETURN",
	"TRUE",
	"FALSE",
	"EQ",
	"NOT_EQ",
}

func (tt TokenType) String() string {
	return tokenTypeStrings[tt]
}

func New(tt TokenType, l string) Token {
	return Token{tt, l}
}

var keywords = map[string]TokenType{
	"fn":  FUNCTION,
	"let": LET,
	"if": IF,
	"else": ELSE,
	"return": RETURN,
	"true": TRUE,
	"false": FALSE,
	"==": EQ,
	"!=": NOT_EQ,
}

func LookupIdentifier(identifier string) TokenType {
	if tok, ok := keywords[identifier]; ok {
		return tok
	}
	return IDENTIFIER
}
