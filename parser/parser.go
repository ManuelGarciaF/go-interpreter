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

	// We associate prefix and infix functions to each token.
	// We save them in maps inside the parser to 'bind' the functions to the parser.
	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
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
	token.LPAREN:   CALL,
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:              l,
		errors:         make([]string, 0),
		prefixParseFns: make(map[token.TokenType]prefixParseFn),
		infixParseFns:  make(map[token.TokenType]infixParseFn),
	}

	// Read 2 tokens to populate curToken and peekToken
	p.nextToken()
	p.nextToken()

	// Bind parseFns
	p.prefixParseFns[token.IDENTIFIER] = p.parseIdentifier
	p.prefixParseFns[token.INT] = p.parseIntegerLiteral
	p.prefixParseFns[token.BANG] = p.parsePrefixExpression
	p.prefixParseFns[token.MINUS] = p.parsePrefixExpression
	p.prefixParseFns[token.TRUE] = p.parseBoolean
	p.prefixParseFns[token.FALSE] = p.parseBoolean
	p.prefixParseFns[token.LPAREN] = p.parseGroupedExpression
	p.prefixParseFns[token.IF] = p.parseIfExpression
	p.prefixParseFns[token.FUNCTION] = p.parseFunctionLiteral

	p.infixParseFns[token.PLUS] = p.parseInfixExpression
	p.infixParseFns[token.MINUS] = p.parseInfixExpression
	p.infixParseFns[token.SLASH] = p.parseInfixExpression
	p.infixParseFns[token.ASTERISK] = p.parseInfixExpression
	p.infixParseFns[token.EQ] = p.parseInfixExpression
	p.infixParseFns[token.NOT_EQ] = p.parseInfixExpression
	p.infixParseFns[token.LT] = p.parseInfixExpression
	p.infixParseFns[token.GT] = p.parseInfixExpression
	p.infixParseFns[token.LPAREN] = p.parseCallExpression

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

	p.nextToken()
	statement.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return statement
}

func (p *Parser) parseReturnStatement() ast.Statement {
	statement := &ast.ReturnStatement{Token: p.currToken}

	p.nextToken()

	statement.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
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
	prefixFn, ok := p.prefixParseFns[p.currToken.Type]
	if !ok {
		p.noPrefixParseFnError(p.currToken.Type)
		return nil
	}
	leftExp := prefixFn()

	// While we are not at a semicolon and the next expression has a higher precedence,
	// parse the next expression.
	for !p.peekTokenIs(token.SEMICOLON) && p.peekPrecedence() > precedence {
		infixFn, ok := p.infixParseFns[p.peekToken.Type]
		if !ok {
			return leftExp
		}

		p.nextToken()

		leftExp = infixFn(leftExp)
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

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.currToken, Value: p.currTokenIs(token.TRUE)}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken() // Advance the starting LPAREN

	// Parse an expression with the lowest precedence, since we inside parethesis
	exp := p.parseExpression(LOWEST)

	// After parsing the expression there must be a closing parens
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
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

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.currToken}

	// After the "if", there must be a opening paren
	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	// Inside the parens, there is an expression for the condition
	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	// Expect a closing parens and brace
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	// Check for else and opening brace
	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	literal := &ast.FunctionLiteral{Token: p.currToken}

	// There should be a paren after the "fn" token
	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	literal.Parameters = p.parseFunctionParameters()

	// there should be an opening brace after the parameters
	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	literal.Body = p.parseBlockStatement()

	return literal
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	b := &ast.BlockStatement{
		Token:      p.currToken,
		Statements: make([]ast.Statement, 0),
	}

	// Skip over the '{'
	p.nextToken()

	for !p.currTokenIs(token.RBRACE) {
		statement := p.parseStatement()

		if statement != nil {
			b.Statements = append(b.Statements, statement)
		}
		p.nextToken()
	}

	return b
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := make([]*ast.Identifier, 0)

	// Special case where there are no parameters
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}

	// Skip over the '('
	p.nextToken()
	// Parse first parameter
	identifiers = append(identifiers, &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal})

	// As long as there is a comma after the current parameter
	for p.peekTokenIs(token.COMMA) {
		// Skip over the current parameter and the comma
		p.nextToken() // ',' is currToken
		p.nextToken() // Identifier is currToken

		identifiers = append(identifiers, &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal})
	}

	// Expect a closing parens
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
}

// The left side of the parens is the function
func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.currToken, Function: function}
	exp.Arguments = p.parseCallArguments()
	return exp
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := make([]ast.Expression, 0)

	// Special case if there are no arguments
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return args
	}

	// Skip over the '('
	p.nextToken()
	// Parse first argument
	args = append(args, p.parseExpression(LOWEST))

	// As long as there is a comma after the current argument
	for p.peekTokenIs(token.COMMA) {
		// Skip over the current argument and the comma
		p.nextToken() // ',' is currToken
		p.nextToken() // Identifier is currToken

		args = append(args, p.parseExpression(LOWEST))
	}

	// Expect a closing parens
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return args
}

//
// Utility functions
//

func (p *Parser) currTokenIs(tt token.TokenType) bool { return p.currToken.Type == tt }
func (p *Parser) peekTokenIs(tt token.TokenType) bool { return p.peekToken.Type == tt }

// Checks if the next token is the expected type and advances to it
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
