package parser

import (
	"fmt"
	"runny/src/token"
	"runny/src/tree"
)

func New() *Parser {
	return &Parser{
		Current:    0,
		Depth:      0,
		Statements: make([]tree.Statement, 0),
		Context:    Context{},
	}
}

type Parser struct {
	Tokens     []token.Token
	Current    int
	Depth      int
	Statements []tree.Statement
	Context    Context
}

func (p *Parser) Parse(tokens []token.Token) (statements []tree.Statement, err error) {
	p.Tokens = tokens
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
	} else if p.check(token.RUN) { // this feels hacky
		modifier := p.peek().Modifier
		p.advance()
		return p.runDeclaration(modifier)
	} else if p.match(token.DESCRIBE) {
		return p.describeDeclaration()
	} else if p.match(token.EXTENDS) {
		return p.extendsDeclaration()
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
		Items:  make([]tree.Variable, 0),
		Parent: p.Context.current(),
	}

	p.Context.setContext(varDecl)

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

	p.Context.resetContext()

	p.reduceDepth()

	return varDecl
}

func (p *Parser) targetDeclaration() tree.Statement {
	name := p.consume(token.IDENTIFIER, "expect target name")

	p.consume(token.LEFT_BRACE, "expect left brace")

	depth := p.increaseDepth()

	targetDecl := tree.TargetStatement{
		Name:   name,
		Body:   make([]tree.Statement, 0),
		Parent: p.Context.current(),
	}

	p.Context.setContext(targetDecl)

	for !p.isAtEnd() {
		body := p.declaration()
		targetDecl.Body = append(targetDecl.Body, body)

		if p.check(token.RIGHT_BRACE) && depth == p.Depth {
			break
		}
	}

	p.consume(token.RIGHT_BRACE, "expect right brace")

	p.Context.resetContext()

	p.reduceDepth()

	return targetDecl
}

func (p *Parser) runDeclaration(modifier *token.TokenModifier) tree.Statement {
	runDecl := tree.RunStatement{
		Body:   make([]tree.Statement, 0),
		Parent: p.Context.current(),
	}

	p.Context.setContext(runDecl)

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
		runDecl.Body = append(runDecl.Body, p.declaration()) // is this correct?
		if p.check(token.RIGHT_BRACE) && depth == p.Depth {
			break
		}
	}

	p.consume(token.RIGHT_BRACE, "expect right brace")

	if modifier != nil {
		switch *modifier {
		case token.BEFORE:
			runDecl.Order = tree.BEFORE
		case token.AFTER:
			runDecl.Order = tree.AFTER
		}
	} else {
		runDecl.Order = tree.DURING
	}

	p.Context.resetContext()
	p.reduceDepth()

	return runDecl
}

func (p *Parser) describeDeclaration() tree.Statement {
	p.consume(token.LEFT_BRACE, "expect left brace")

	depth := p.increaseDepth()

	descDecl := tree.DescribeStatement{
		Lines:  make([]tree.Literal, 0),
		Parent: p.Context.current(),
	}

	p.Context.setContext(descDecl)

	for !p.isAtEnd() {
		initialiser := p.declaration()

		var value interface{}
		if statement, ok := initialiser.(tree.ExpressionStatement); ok {
			if literal, ok := statement.Expression.(tree.Literal); ok {
				value = literal.Value
			}
		}

		if value == nil {
			panic(p.error(p.peek(), "description not found"))
		}

		descDecl.Lines = append(descDecl.Lines, tree.Literal{
			Value: value,
		})

		if p.check(token.COMMA) {
			p.advance()
		}

		if p.check(token.RIGHT_BRACE) && depth == p.Depth {
			break
		}
	}

	p.consume(token.RIGHT_BRACE, "expect right brace")

	p.Context.resetContext()

	p.reduceDepth()

	return descDecl
}

func (p *Parser) extendsDeclaration() tree.Statement {
	p.consume(token.LEFT_BRACE, "expect left brace")

	depth := p.increaseDepth()

	extends := tree.ExtendsStatement{}

	for !p.isAtEnd() {
		extends.Paths = append(extends.Paths, p.expression())

		if p.check(token.COMMA) {
			p.advance()
		}

		if p.check(token.RIGHT_BRACE) && depth == p.Depth {
			break
		}
	}

	p.consume(token.RIGHT_BRACE, "expect right brace")

	p.reduceDepth()

	return extends
}

func (p *Parser) actionStatement() tree.Statement {
	script := p.consume(token.SCRIPT, "expect action body")

	actionstatement := tree.ActionStatement{
		Body:   script,
		Parent: p.Context.current(),
	}

	return actionstatement
}

func (p *Parser) expressionStatement() tree.Statement {
	exprstatement := tree.ExpressionStatement{
		Expression: p.expression(),
	}
	return exprstatement
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

type Context struct {
	Stack []tree.Statement
}

func (c *Context) current() tree.Statement {
	if len(c.Stack) == 0 {
		return nil
	}
	if len(c.Stack) == 1 {
		return c.Stack[0]
	}
	return c.Stack[len(c.Stack)-1] // last item
}

func (c *Context) setContext(s tree.Statement) {
	c.Stack = append(c.Stack, s)
}

func (c *Context) replaceContext(s tree.Statement) {
	c.resetContext()
	c.setContext(s)
}

// trims last item
func (c *Context) resetContext() {
	if len(c.Stack) > 0 {
		c.Stack = c.Stack[:len(c.Stack)-1]
	}
}
