package ast

import (
	"strings"

	"github.com/ManuelGarciaF/go-interpreter/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode() // Dummy method for compiler error detection
}

// An expression produces a value
type Expression interface {
	Node
	expressionNode() // Dummy method for compiler error detection
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) String() string {
	var sb strings.Builder

	for _, s := range p.Statements {
		sb.WriteString(s.String())
	}

	return sb.String()
}

type LetStatement struct {
	Token token.Token // token.LET
	Name  *Identifier
	Value Expression
}

// Implements Statement
func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) String() string {
	var sb strings.Builder

	sb.WriteString(ls.TokenLiteral())
	sb.WriteString(" ")
	sb.WriteString(ls.Name.String())
	sb.WriteString(" = ")

	if ls.Value != nil { // TODO remove nil check
		sb.WriteString(ls.Value.String())
	}
	sb.WriteString(";")

	return sb.String()
}

type Identifier struct {
	Token token.Token // token.IDENTIFIER
	Value string
}

// Implements Expression, since identifiers do produce values, just not in let statements
func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

type ReturnStatement struct {
	Token token.Token // token.RETURN
	Value Expression
}

// Implements Statement
func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var sb strings.Builder

	sb.WriteString(rs.TokenLiteral())
	sb.WriteString(" ")

	if rs.Value != nil { // TODO remove nil check
		sb.WriteString(rs.Value.String())
	}
	sb.WriteString(";")

	return sb.String()
}

type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

// Implements Statement
func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil { // TODO remove nil check
		return es.Expression.String()
	}
	return ""
}

type IntegerLiteral struct {
	Token token.Token // token.INT
	Value int64
}

// Implements Expression
func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

type PrefixExpression struct {
	Token    token.Token // The prefix token, token.MINUS or token.BANG.
	Operator string      // "-" or "!"
	Right    Expression  // The expression on the right of the operator
}

// Implements Expression
func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	return "(" + pe.Operator + pe.Right.String() + ")"
}

type InfixExpression struct {
	Token    token.Token // The operator's token
	Left     Expression  // The expression on the left of the operator
	Operator string      // "+", "-", "*", etc.
	Right    Expression  // The expression on the right of the operator
}

// Implements Expression
func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	return "(" + ie.Left.String() + " " + ie.Operator + " " + ie.Right.String() + ")"
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string       { return b.Token.Literal }

type IfExpression struct {
	Token       token.Token // token.IF
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

// Implements Expression
func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var sb strings.Builder

	sb.WriteString("if ")
	sb.WriteString(ie.Condition.String())
	sb.WriteString(" ")
	sb.WriteString(ie.Consequence.String())
	if ie.Alternative != nil {
		sb.WriteString("else ")
		sb.WriteString(ie.Alternative.String())
	}

	return sb.String()
}

type BlockStatement struct {
	Token      token.Token // The opening '{'
	Statements []Statement
}

// Implements Statement
func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var sb strings.Builder

	sb.WriteString("{")
	for _, s := range bs.Statements {
		sb.WriteString(s.String())
	}
	sb.WriteString("}")

	return sb.String()
}

type FunctionLiteral struct {
	Token      token.Token // token.FUNCTION
	Parameters []*Identifier
	Body       *BlockStatement
}

// Implements Expression
func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var sb strings.Builder

	params := make([]string, 0, len(fl.Parameters))
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	sb.WriteString("fn(")
	sb.WriteString(strings.Join(params, ", "))
	sb.WriteString(")")
	sb.WriteString(fl.Body.String())

	return sb.String()
}

type CallExpression struct {
	Token     token.Token // token.LPAREN
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

// Implements Expression
func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var sb strings.Builder

	args := make([]string, 0, len(ce.Arguments))
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	sb.WriteString(ce.Function.String())
	sb.WriteString("(")
	sb.WriteString(strings.Join(args, ", "))
	sb.WriteString(")")

	return sb.String()
}
