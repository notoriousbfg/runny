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
	var statement tree.Statement
	if p.match(token.VAR) {
		statement = p.varDeclaration()
	}
	if p.match(token.TARGET) {
		statement = p.targetDeclaration()
	}
	if p.match(token.RUN) {
		statement = p.runDeclaration()
	}
	// there could be any number of newlines after a block
	if statement != nil {
		for p.check(token.NEWLINE) {
			p.advance()
		}
		return statement
	}
	return p.expressionStatement()
}

func (p *Parser) varDeclaration() tree.Statement {
	p.consume(token.LEFT_BRACE, "expect left brace")

	p.skipNewline()

	varDecl := tree.VariableStatement{
		Items: make([]tree.Variable, 0),
	}

	for !p.check(token.RIGHT_BRACE) && !p.isAtEnd() {
		name := p.consume(token.IDENTIFIER, "expect variable name")

		var initialiser tree.Statement
		if p.match(token.LEFT_BRACE) {
			p.skipNewline()
			initialiser = p.parseBlock() // var is the output of an evaluated block e.g. var name { echo "tim" }
			p.skipNewline()
			p.consume(token.RIGHT_BRACE, "expect right brace")
		} else if p.match(token.IDENTIFIER) {
			p.advance()
		} else {
			initialiser = p.declaration()
		}

		varDecl.Items = append(varDecl.Items, tree.Variable{
			Name:        name,
			Initialiser: initialiser,
		})

		if p.check(token.COMMA) {
			p.advance()
		}

		p.skipNewline()
	}

	p.consume(token.RIGHT_BRACE, "expect right brace")

	return varDecl
}

func (p *Parser) targetDeclaration() tree.Statement {
	name := p.consume(token.IDENTIFIER, "expect target name")

	p.consume(token.LEFT_BRACE, "expect left brace")

	p.skipNewline()

	targetDecl := tree.TargetStatement{
		Name: name,
		Body: make([]tree.Statement, 0),
	}

	for !p.check(token.RIGHT_BRACE) && !p.isAtEnd() {
		targetDecl.Body = append(targetDecl.Body, p.parseBlock())
		p.skipNewline()
	}

	p.consume(token.RIGHT_BRACE, "expect right brace")

	return targetDecl
}

func (p *Parser) runDeclaration() tree.Statement {
	runDecl := tree.RunStatement{}

	if p.check(token.IDENTIFIER) {
		name := p.consume(token.IDENTIFIER, "expect target name")
		runDecl.Name = &name
		// TODO: no braces
	}

	p.consume(token.LEFT_BRACE, "expect left brace")

	p.skipNewline()

	for !p.check(token.RIGHT_BRACE) && !p.isAtEnd() {
		runDecl.Body = append(runDecl.Body, p.parseBlock())
	}

	p.consume(token.RIGHT_BRACE, "expect right brace")

	return runDecl
}

func (p *Parser) actionStatement() tree.Statement {
	script := p.consume(token.SCRIPT, "expect action body")

	return tree.ActionStatement{
		Body: script,
	}
}

func (p *Parser) parseBlock() tree.Statement {
	if isKeyword(p.peek()) {
		return p.declaration()
	} else {
		return p.actionStatement()
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

	panic(p.error(p.peek(), "expect expression"))
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

func (p *Parser) checkSequence(tokenTypes ...token.TokenType) bool {
	start := p.Current
	defer func() {
		p.Current = start
	}()
	for _, tokenType := range tokenTypes {
		if !p.check(tokenType) {
			return false
		}
		p.advance()
	}
	return true
}

func (p *Parser) advance() token.Token {
	if !p.isAtEnd() {
		p.Current++
	}
	return p.previous()
}

// get the previous token
func (p *Parser) previous() token.Token {
	return p.Tokens[p.Current-1]
}

// get the token at the current index
func (p *Parser) peek() token.Token {
	return p.Tokens[p.Current]
}

// get the token at the next index
func (p *Parser) next() token.Token {
	return p.Tokens[p.Current+1]
}

// if the token is of the specified type advance, otherwise panic
func (p *Parser) consume(tokenType token.TokenType, message string) token.Token {
	if p.check(tokenType) {
		return p.advance()
	}

	panic(p.error(p.peek(), message))
}

func (p *Parser) skipNewline() {
	if p.check(token.NEWLINE) {
		p.advance()
	}
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == token.EOF
}

type ParseError struct {
	Message string
}

func (pe *ParseError) Error() string {
	return pe.Message
}

func (p *Parser) error(thisToken token.Token, message string) *ParseError {
	var where string
	if thisToken.Type == token.EOF {
		where = "at end"
	} else if thisToken.Type == token.NEWLINE {
		where = "at '\\n'"
	} else {
		where = "at '" + thisToken.Text + "'"
	}
	err := &ParseError{
		Message: fmt.Sprintf("[line %d] error %s: %s\n", thisToken.Line, where, message),
	}
	return err
}

func isKeyword(t token.Token) bool {
	for _, keyword := range token.Keywords {
		if t.Type == keyword {
			return true
		}
	}
	return false
}
