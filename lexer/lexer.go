package lexer

import (
	"github.com/ManuelGarciaF/go-interpreter/token"
)

type Lexer struct {
	input        string
	position     int  // Pos of current char (ch)
	readPosition int  // Pos of next char to read
	ch           byte // Only ascii for now
}

const EOF byte = 0

func New(input string) *Lexer {
	l := &Lexer{
		input:        input,
		position:     0,
		readPosition: 0,
		ch:           0,
	}
	l.readChar() // Have to initialize with a first read
	return l
}

func (l *Lexer) readChar() {
	// If out of bounds
	if l.readPosition >= len(l.input) {
		l.ch = EOF
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1

}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	// NOTE: Could simplify this by extracting all the simple cases into a map
	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			first := l.ch
			// Advance a char.
			l.readChar()
			tok = token.New(token.EQ, string(first)+string(l.ch))
		} else {
			tok = token.New(token.ASSIGN, string(l.ch))
		}
	case '+':
		tok = token.New(token.PLUS, string(l.ch))
	case '-':
		tok = token.New(token.MINUS, string(l.ch))
	case '!':
		if l.peekChar() == '=' {
			first := l.ch
			// Advance a char.
			l.readChar()
			tok = token.New(token.NOT_EQ, string(first)+string(l.ch))
		} else {
			tok = token.New(token.BANG, string(l.ch))
		}
	case '/':
		tok = token.New(token.SLASH, string(l.ch))
	case '*':
		tok = token.New(token.ASTERISK, string(l.ch))
	case '<':
		tok = token.New(token.LT, string(l.ch))
	case '>':
		tok = token.New(token.GT, string(l.ch))
	case ',':
		tok = token.New(token.COMMA, string(l.ch))
	case ';':
		tok = token.New(token.SEMICOLON, string(l.ch))
	case ':':
		tok = token.New(token.COLON, string(l.ch))
	case '(':
		tok = token.New(token.LPAREN, string(l.ch))
	case ')':
		tok = token.New(token.RPAREN, string(l.ch))
	case '{':
		tok = token.New(token.LBRACE, string(l.ch))
	case '}':
		tok = token.New(token.RBRACE, string(l.ch))
	case '[':
		tok = token.New(token.LBRACKET, string(l.ch))
	case ']':
		tok = token.New(token.RBRACKET, string(l.ch))
	case '"':
		// There are cases in which we don't find a complete string.
		string, ok := l.readString()
		if ok {
			tok = token.New(token.STRING, string)
		} else {
			tok = token.New(token.EOF, "")
		}

	case EOF:
		tok = token.New(token.EOF, "")
	default:
		// We have to check for identifiers if there are letters.
		if isLetter(l.ch) {
			literal := l.readIdentifier()
			// We return early so we don't advance an extra character.
			return token.New(token.LookupIdentifier(literal), literal)
		} else if isDigit(l.ch) { // Check for ints.
			num := l.readNumber()
			// We return early so we don't advance an extra character.
			return token.New(token.INT, num)
		} else { // If it does not start with a letter it's not a valid token.
			tok = token.New(token.ILLEGAL, string(l.ch))
		}

	}

	l.readChar()
	return tok

}

func (l *Lexer) readIdentifier() string {
	initialPos := l.position

	// We already checked the first one is exclusively a letter before.
	for isValidInIdentifier(l.ch) {
		l.readChar()
	}

	// The current ch is not part of the identifier, so we use l.position.
	return l.input[initialPos:l.position]
}

func (l *Lexer) readNumber() string {
	initialPos := l.position

	// No support for floats
	for isDigit(l.ch) {
		l.readChar()
	}
	// The current ch is not part of the identifier, so we use l.position.
	return l.input[initialPos:l.position]
}

// Returns the string, and ok. ok is false if a closing '"' couldnt' be found.
func (l *Lexer) readString() (string, bool) {
	// Advance the first '"'
	l.readChar()

	start := l.position
	for l.ch != '"' {
		l.readChar()
		if l.ch == EOF {
			return "", false
		}
	}
	return l.input[start:l.position], true
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) peekChar() byte {
	// If out of bounds.
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition] // The following character.
}

func isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') ||
		(ch >= 'A' && ch <= 'Z')
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

// For characters after the first one, we allow underscores and numbers
func isValidInIdentifier(ch byte) bool {
	return isLetter(ch) || ch == '_' || isDigit(ch)
}
