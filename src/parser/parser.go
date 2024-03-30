package parser

import (
	"fmt"
	"runny/src/token"
	"runny/src/tree"
)

func New(tokens []token.Token) *Parser {
	return &Parser{
		Tokens:     tokens,
		Current:    0,
		Depth:      0,
		Statements: make([]tree.Statement, 0),
	}
}

type Parser struct {
	Tokens     []token.Token
	Current    int
	Depth      int
	Statements []tree.Statement
}

func (p *Parser) Parse() (statements []tree.Statement, err error) {
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
	for !p.isAtEnd() {
		p.Statements = append(p.Statements, p.declaration())
	}
	statements = p.Statements
	return
}

func (p *Parser) declaration() tree.Statement {
	if p.match(token.CONFIG) {
		return p.configDeclaration()
	} else if p.match(token.VAR) {
		return p.varDeclaration()
	} else if p.match(token.TARGET) {
		return p.targetDeclaration()
	} else if p.match(token.RUN) {
		return p.runDeclaration()
	} else if p.check(token.SCRIPT) {
		return p.actionStatement()
	}
	return p.expressionStatement()
}

func (p *Parser) configDeclaration() tree.Statement {
	p.consume(token.LEFT_BRACE, "expect left brace")

	depth := p.increaseDepth()

	configDecl := tree.ConfigStatement{
		Items: make([]tree.Config, 0),
	}

	for !p.isAtEnd() {
		name := p.consume(token.IDENTIFIER, "expect config variable")
		initialiser := p.declaration()

		configDecl.Items = append(configDecl.Items, tree.Config{
			Name:        name,
			Initialiser: initialiser,
		})

		if p.check(token.COMMA) {
			p.advance()
		}

		if p.check(token.RIGHT_BRACE) && depth == p.Depth {
			break
		}
	}

	p.consume(token.RIGHT_BRACE, "expect right brace")

	p.reduceDepth()

	return configDecl
}

func (p *Parser) varDeclaration() tree.Statement {
	p.consume(token.LEFT_BRACE, "expect left brace")

	depth := p.increaseDepth()

	varDecl := tree.VariableStatement{
		Items: make([]tree.Variable, 0),
	}

	for !p.isAtEnd() {
		name := p.consume(token.IDENTIFIER, "expect variable name")

		var initialiser tree.Statement
		if p.match(token.LEFT_BRACE) {
			initialiser = p.declaration() // var is the output of an evaluated block e.g. var name { run { echo "tim" } }
			p.consume(token.RIGHT_BRACE, "expect right brace")
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

		if p.check(token.RIGHT_BRACE) && depth == p.Depth {
			break
		}
	}

	p.consume(token.RIGHT_BRACE, "expect right brace")

	p.reduceDepth()

	return varDecl
}

func (p *Parser) targetDeclaration() tree.Statement {
	name := p.consume(token.IDENTIFIER, "expect target name")

	p.consume(token.LEFT_BRACE, "expect left brace")

	depth := p.increaseDepth()

	targetDecl := tree.TargetStatement{
		Name: name,
		Body: make([]tree.Statement, 0),
	}

	for !p.isAtEnd() {
		body := p.declaration()
		targetDecl.Body = append(targetDecl.Body, body)

		if p.check(token.RIGHT_BRACE) && depth == p.Depth {
			break
		}
	}

	p.consume(token.RIGHT_BRACE, "expect right brace")

	p.reduceDepth()

	return targetDecl
}

func (p *Parser) runDeclaration() tree.Statement {
	runDecl := tree.RunStatement{
		Body: make([]tree.Statement, 0),
	}

	if p.check(token.IDENTIFIER) {
		name := p.consume(token.IDENTIFIER, "expect target name")
		runDecl.Name = name

		if !p.check(token.LEFT_BRACE) {
			return runDecl
		}
	}

	p.consume(token.LEFT_BRACE, "expect left brace")

	depth := p.increaseDepth()

	for !p.isAtEnd() {
		runDecl.Body = append(runDecl.Body, p.declaration())
		if p.check(token.RIGHT_BRACE) && depth == p.Depth {
			break
		}
	}

	p.consume(token.RIGHT_BRACE, "expect right brace")

	p.reduceDepth()

	return runDecl
}

func (p *Parser) actionStatement() tree.Statement {
	script := p.consume(token.SCRIPT, "expect action body")

	return tree.ActionStatement{
		Body: script,
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
	if p.match(token.IDENTIFIER) {
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

// if the token is of the specified type advance, otherwise panic
func (p *Parser) consume(tokenType token.TokenType, message string) token.Token {
	if p.check(tokenType) {
		return p.advance()
	}

	panic(p.error(p.peek(), message))
}

func (p *Parser) increaseDepth() int {
	p.Depth++
	return p.Depth
}

func (p *Parser) reduceDepth() int {
	p.Depth--
	return p.Depth
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
		Message: fmt.Sprintf("[line %d] parse error %s: %s\n", thisToken.Line, where, message),
	}
	return err
}
