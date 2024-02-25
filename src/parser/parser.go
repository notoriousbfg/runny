package parser

import (
	"fmt"
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
	Tokens     []token.Token
	Current    int
	Statements []tree.Statement
}

func (p *Parser) Parse() (err error) {
	defer func() {
		if r := recover(); r != nil {
			if str, ok := r.(string); ok {
				err = fmt.Errorf(str)
			} else if e, ok := r.(error); ok {
				err = e
			} else {
				err = fmt.Errorf("unknown panic: %v", r)
			}
		}
	}()
	statements := make([]tree.Statement, 0)
	for !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}
	p.Statements = statements
	return nil
}

func (p *Parser) declaration() tree.Statement {
	if p.match(token.VAR) {
		return p.varDeclaration()
	}
	if p.match(token.TARGET) {
		return p.targetDeclaration()
	}
	if p.match(token.RUN) {
		return p.runDeclaration()
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
		Body: make([]tree.Statement, 0),
	}

	for !p.check(token.RIGHT_BRACE) && !p.isAtEnd() {
		// anything that isn't a keyword will be parsed as an action
		if isKeyword(p.peek()) {
			targetDecl.Body = append(targetDecl.Body, p.declaration())
		} else {
			targetDecl.Body = append(targetDecl.Body, p.actionStatement())
		}
	}

	p.consume(token.RIGHT_BRACE, "expect right brace")

	return targetDecl
}

func (p *Parser) runDeclaration() tree.Statement {
	runDecl := tree.RunStatement{
		Body: make([]tree.Statement, 0),
	}

	fmt.Println(p.peek())

	if p.check(token.IDENTIFIER) {
		name := p.consume(token.IDENTIFIER, "expect target name")
		runDecl.Name = &name
	}

	p.consume(token.LEFT_BRACE, "expect left brace")

	for !p.check(token.RIGHT_BRACE) && !p.isAtEnd() {
		if isKeyword(p.peek()) {
			runDecl.Body = append(runDecl.Body, p.declaration())
		} else {
			runDecl.Body = append(runDecl.Body, p.actionStatement())
		}
	}

	p.consume(token.RIGHT_BRACE, "expect right brace")

	return runDecl
}

func (p *Parser) actionStatement() tree.Statement {
	tokens := make([]token.Token, 0)

	for !isKeyword(p.peek()) && !p.check(token.RIGHT_BRACE) && !p.isAtEnd() {
		// TODO: newline
		tokens = append(tokens, p.peek())
		p.advance()
	}

	return tree.ActionStatement{
		Body: tokens,
	}
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

func isKeyword(t token.Token) (isKeyword bool) {
	_, isKeyword = token.Keywords[t.Text]
	return
}
