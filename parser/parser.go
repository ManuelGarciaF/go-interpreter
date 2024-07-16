package parser

import (
	"fmt"
	"strconv"

	"github.com/ManuelGarciaF/go-interpreter/ast"
	"github.com/ManuelGarciaF/go-interpreter/lexer"
	"github.com/ManuelGarciaF/go-interpreter/token"
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	currToken token.Token
	peekToken token.Token

	// We associate prefix and infix functions to each token
	// We save them in maps inside the parser to 'bind' the functions to the parser
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns   map[token.TokenType]infixParseFn
}

type precedence int

const (
	LOWEST      precedence = iota
	EQUALS                 // ==
	LESSGREATER            // > or <
	SUM                    // +
	PRODUCT                // *
	PREFIX                 // -x or !x
	CALL                   // x()
)

var precedences = map[token.TokenType]precedence{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:              l,
		errors:         make([]string, 0),
		prefixParseFns: make(map[token.TokenType]prefixParseFn),
		infixParseFns:   make(map[token.TokenType]infixParseFn),
	}

	// Read 2 tokens to populate curToken and peekToken
	p.nextToken()
	p.nextToken()

	// Bind parseFns
	p.prefixParseFns[token.IDENTIFIER] = p.parseIdentifier
	p.prefixParseFns[token.INT] = p.parseIntegerLiteral
	p.prefixParseFns[token.BANG] = p.parsePrefixExpression
	p.prefixParseFns[token.MINUS] = p.parsePrefixExpression
	p.infixParseFns[token.PLUS] = p.parseInfixExpression
	p.infixParseFns[token.MINUS] = p.parseInfixExpression
	p.infixParseFns[token.SLASH] = p.parseInfixExpression
	p.infixParseFns[token.ASTERISK] = p.parseInfixExpression
	p.infixParseFns[token.EQ] = p.parseInfixExpression
	p.infixParseFns[token.NOT_EQ] = p.parseInfixExpression
	p.infixParseFns[token.LT] = p.parseInfixExpression
	p.infixParseFns[token.GT] = p.parseInfixExpression

	return p
}

func (p *Parser) nextToken() {
	p.currToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = make([]ast.Statement, 0)

	// Parse each statement one by one
	for !p.currTokenIs(token.EOF) {
		statement := p.parseStatement()
		if statement != nil { // TODO check for this error
			program.Statements = append(program.Statements, statement)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() ast.Statement {
	statement := &ast.LetStatement{Token: p.currToken}

	// At this point, curr = LET, peek should be an IDENTIFIER.
	if !p.expectPeek(token.IDENTIFIER) {
		return nil
	}
	// expectPeek advanced the currToken to the identifier
	statement.Name = &ast.Identifier{
		Token: p.currToken,
		Value: p.currToken.Literal,
	}

	// After the identifier, we expect an '='
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO temporarily skipping expressions until a semicolon
	for !p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	// After the value, we expect a semicolon
	if !p.expectPeek(token.SEMICOLON) {
		return nil
	}

	return statement
}

func (p *Parser) parseReturnStatement() ast.Statement {
	statement := &ast.ReturnStatement{Token: p.currToken}

	p.nextToken()

	// TODO temporarily skipping expressions until a semicolon
	for !p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	// After the value, we expect a semicolon
	if !p.expectPeek(token.SEMICOLON) {
		return nil
	}

	return statement
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	statement := &ast.ExpressionStatement{Token: p.currToken}
	statement.Expression = p.parseExpression(LOWEST)

	// Semicolons are optional in expression statements, so its easier to use the REPL
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return statement
}

func (p *Parser) parseExpression(precedence precedence) ast.Expression {
	prefix, ok := p.prefixParseFns[p.currToken.Type]
	if !ok {
		p.noPrefixParseFnError(p.currToken.Type)
		return nil
	}
	leftExp := prefix()

	// While we are not at a semicolon and the next expression has a higher precedence
	for !p.peekTokenIs(token.SEMICOLON) && p.peekPrecedence() > precedence {
		infix, ok := p.infixParseFns[p.peekToken.Type]
		if !ok {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}
	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	literal := &ast.IntegerLiteral{Token: p.currToken}
	value, err := strconv.ParseInt(p.currToken.Literal, 10, 64)
	if err != nil {
		p.errors = append(p.errors,
			fmt.Sprintf("Could not parse %q as an integer", p.currToken.Literal),
		)
		return nil
	}

	literal.Value = value

	return literal
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.currToken,
		Operator: p.currToken.Literal,
	}

	// We need to consume the prefix symbol to parse the expression
	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.currToken,
		Left:     left,
		Operator: p.currToken.Literal,
	}
	precedence := p.currPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

//
// Utility functions
//

func (p *Parser) currTokenIs(tt token.TokenType) bool { return p.currToken.Type == tt }
func (p *Parser) peekTokenIs(tt token.TokenType) bool { return p.peekToken.Type == tt }

// Checks if the next token is the expected type and advances to the next token
func (p *Parser) expectPeek(tt token.TokenType) bool {
	if !p.peekTokenIs(tt) {
		p.peekError(tt)
		return false
	}
	p.nextToken()
	return true
}

func (p *Parser) currPrecedence() precedence {
	if p, ok := precedences[p.currToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) peekPrecedence() precedence {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) peekError(expected token.TokenType) {
	p.errors = append(p.errors,
		fmt.Sprintf("Expected next token to be %s, got %s", expected, p.peekToken.Type),
	)
}

func (p *Parser) noPrefixParseFnError(tt token.TokenType) {
	p.errors = append(p.errors,
		fmt.Sprintf("No prefix parse function for %s", tt),
	)
}
