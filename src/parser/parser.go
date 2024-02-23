package parser

import (
	"runny/src/token"
	"runny/src/tree"
)

func New(tokens []token.Token) *Parser {
	return &Parser{
		Tokens:  tokens,
		Current: 0,
	}
}

type Parser struct {
	Tokens  []token.Token
	Current int
}

func (p *Parser) Parse() []tree.Statement {
	statements := make([]tree.Statement, 0)
	for !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}
	return statements
}

func (p *Parser) declaration() tree.Statement {
	if p.match(token.VAR) {
		return p.varDeclaration()
	}
	if p.match(token.TARGET) {
		return p.targetDeclaration()
	}
	return p.expressionStatement()
}

func (p *Parser) varDeclaration() tree.Statement {
	p.consume(token.LEFT_BRACE, "expect left brace")

	varDecl := tree.VariableStatement{
		Items: make([]tree.Variable, 0),
	}

	for !p.check(token.RIGHT_BRACE) && !p.isAtEnd() {
		if p.check(token.COMMA) {
			p.advance()
		}

		varDecl.Items = append(varDecl.Items, tree.Variable{
			Name:        p.consume(token.IDENTIFIER, "expect variable name"),
			Initialiser: p.expression(),
		})
	}

	p.consume(token.RIGHT_BRACE, "expect right brace")

	return varDecl
}

func (p *Parser) targetDeclaration() tree.Statement {
	name := p.consume(token.IDENTIFIER, "expect target name")

	p.consume(token.LEFT_BRACE, "expect left brace")

	targetDecl := tree.TargetStatement{
		Name: name,
		Body: make([]token.Token, 0),
	}

	for !p.check(token.RIGHT_BRACE) && !p.isAtEnd() {
		targetDecl.Body = append(targetDecl.Body, p.advance())
	}

	p.consume(token.RIGHT_BRACE, "expect right brace")

	return targetDecl
}

func (p *Parser) expressionStatement() tree.Statement {
	exprStmt := tree.ExpressionStatement{
		Expression: p.expression(),
	}
	return exprStmt
}

func (p *Parser) expression() tree.Expression {
	if p.match(token.NUMBER, token.STRING) {
		return tree.Literal{Value: p.previous().Text}
	}

	panic("expect expression")
}

// check that the current token is any of the types and advance if so
func (p *Parser) match(tokenTypes ...token.TokenType) bool {
	for _, tokenType := range tokenTypes {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}
	return false
}

// check that the current token is of a type
func (p *Parser) check(tokenType token.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == tokenType
}

// get the previous token
func (p *Parser) previous() token.Token {
	return p.Tokens[p.Current-1]
}

func (p *Parser) advance() token.Token {
	if !p.isAtEnd() {
		p.Current++
	}
	return p.previous()
}

// get the token at the current index
func (p *Parser) peek() token.Token {
	return p.Tokens[p.Current]
}

// if the token is of the specified type advance, otherwise panic
func (p *Parser) consume(tokenType token.TokenType, message string) token.Token {
	if p.check(tokenType) {
		return p.advance()
	}

	panic(message)
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == token.EOF
}
