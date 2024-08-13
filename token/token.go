package token

type TokenType uint8

// TODO add line numbers
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
	STRING

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
	LBRACKET
	RBRACKET

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
var tokenTypeStrings = map[TokenType]string{
	ILLEGAL:    "ILLEGAL",
	EOF:        "EOF",
	IDENTIFIER: "IDENTIFIER",
	INT:        "INT",
	STRING:     "STRING",
	ASSIGN:     "ASSIGN",
	PLUS:       "PLUS",
	MINUS:      "MINUS",
	BANG:       "BANG",
	ASTERISK:   "ASTERISK",
	SLASH:      "SLASH",
	LT:         "LT",
	GT:         "GT",
	COMMA:      "COMMA",
	SEMICOLON:  "SEMICOLON",
	LPAREN:     "LPAREN",
	RPAREN:     "RPAREN",
	LBRACE:     "LBRACE",
	RBRACE:     "RBRACE",
	LBRACKET:   "LBRACKET",
	RBRACKET:   "RBRACKET",
	FUNCTION:   "FUNCTION",
	LET:        "LET",
	IF:         "IF",
	ELSE:       "ELSE",
	RETURN:     "RETURN",
	TRUE:       "TRUE",
	FALSE:      "FALSE",
	EQ:         "EQ",
	NOT_EQ:     "NOT_EQ",
}

func (tt TokenType) String() string {
	return tokenTypeStrings[tt]
}

func New(tt TokenType, l string) Token {
	return Token{tt, l}
}

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"true":   TRUE,
	"false":  FALSE,
	"==":     EQ,
	"!=":     NOT_EQ,
}

func LookupIdentifier(identifier string) TokenType {
	if tok, ok := keywords[identifier]; ok {
		return tok
	}
	return IDENTIFIER
}
