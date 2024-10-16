package lex

import (
	"fmt"
	"runny/src/token"
	"strconv"
	"strings"
	"unicode"
)

func New() *Lexer {
	return &Lexer{
		Line:    1,
		Start:   0,
		Current: 0,
		Depth:   0,
		Context: Context{},
	}
}

type Lexer struct {
	Input   string
	Tokens  []token.Token
	Start   int
	Current int
	Line    int
	Depth   int // number of braces deep
	Context Context
}

func (l *Lexer) ReadInput(input string) ([]token.Token, error) {
	l.Input = input
	for !l.isAtEnd() {
		l.Start = l.Current
		err := l.readChar()
		if err != nil {
			return []token.Token{}, err
		}
	}
	l.Start++
	l.addToken(token.EOF, "")
	return l.Tokens, nil
}

func (l *Lexer) readChar() error {
	char := l.nextChar()

	switch char {
	case "{":
		l.addToken(token.LEFT_BRACE, char)
		l.Depth++
		if l.Context.current() == token.RUN {
			l.matchScript()
		}
	case "}":
		l.addToken(token.RIGHT_BRACE, char)
		l.Depth--
		l.Context.resetContext()
	case ",":
		l.addToken(token.COMMA, char)
	case "$":
		l.matchIdentifier()
	case "#":
		l.matchComment()
	case "\n":
		l.Line++
	case " ", "\r", "\t":
		break
	default:
		if isDigit(char) {
			l.matchNumber()
		} else if isLetter(char) {
			l.matchIdentifier()
		} else if char == "`" || char == "\"" {
			l.matchString(char)
		} else {
			return l.error(char, "unsupported type")
		}
	}
	return nil
}

func (l *Lexer) nextChar() string {
	char := string(l.Input[l.Current])
	l.Current++
	return char
}

type TokenOptionFunc func(*token.Token)

func withModifier(modifier token.TokenModifier) func(*token.Token) {
	return func(token *token.Token) {
		token.Modifier = &modifier
	}
}

func withPosition(position, line, depth int) func(*token.Token) {
	return func(token *token.Token) {
		token.Position = position
		token.Depth = depth
		token.Line = line
	}
}

func (l *Lexer) addToken(tokenType token.TokenType, text string, options ...TokenOptionFunc) {
	token := token.Token{
		Type:     tokenType,
		Text:     text,
		Position: l.Start,
		Line:     l.Line,
		Depth:    l.Depth,
	}
	for _, opt := range options {
		opt(&token)
	}
	l.Tokens = append(l.Tokens, token)
}

func (l *Lexer) TokenTypes() []token.TokenType {
	var types []token.TokenType
	for _, token := range l.Tokens {
		types = append(types, token.Type)
	}
	return types
}

func (l *Lexer) isAtEnd() bool {
	return l.Current >= len(l.Input)
}

func (l *Lexer) lastToken() token.Token {
	if len(l.Tokens) == 0 {
		return token.Token{}
	}
	if len(l.Tokens) == 1 {
		return l.Tokens[0]
	}
	return l.Tokens[len(l.Tokens)-1]
}

func (l *Lexer) peek() string {
	if l.isAtEnd() {
		return ""
	}
	return string(l.Input[l.Current])
}

func (l *Lexer) matchComment() {
	for !l.isAtEnd() && l.peek() != "\n" {
		l.nextChar()
	}
}

func (l *Lexer) matchScript() {
	start := l.Start + 1
	firstLine := l.Line
	bracesCount := 1
	for !l.isAtEnd() {
		if l.peek() == "\n" {
			l.Line++
		} else if l.peek() == "{" {
			bracesCount++
		} else if l.peek() == "}" {
			bracesCount--
			if bracesCount == 0 {
				break
			}
		}
		l.nextChar()
	}
	text := l.Input[start:l.Current]
	if len(text) > 0 {
		l.addToken(
			token.SCRIPT,
			strings.TrimSpace(text),
			withPosition(start, firstLine, l.Depth),
		)
	}
}

func (l *Lexer) matchString(delimiter string) {
	for l.peek() != delimiter && !l.isAtEnd() {
		if l.peek() == "\n" {
			l.Line++
		}
		l.nextChar()
	}
	l.nextChar()
	text := l.Input[l.Start+1 : l.Current-1]
	l.addToken(token.STRING, fmt.Sprintf("\"%s\"", text))
}

func (l *Lexer) matchNumber() {
	for isDigit(l.peek()) {
		l.nextChar()
	}
	l.addToken(token.NUMBER, l.Input[l.Start:l.Current])
}

func (l *Lexer) matchIdentifier() {
	identifier := l.readIdentifier()
	if keyword, isKeyword := l.isKeyword(identifier); isKeyword {
		l.addToken(keyword, identifier)
		l.Context.setContext(keyword)
	} else if keyword, mod, _, hasTag := l.hasModifier(identifier); hasTag {
		l.addToken(*keyword, identifier, withModifier(*mod))
		l.Context.setContext(*keyword)
	} else {
		// we're in target context if running target
		if l.lastToken().Type == token.RUN {
			l.Context.replaceContext(token.TARGET)
		}
		l.addToken(token.IDENTIFIER, identifier)
	}
}

func (l *Lexer) isKeyword(identifier string) (token.TokenType, bool) {
	keyword, isKeyword := token.Keywords[identifier]
	return keyword, isKeyword
}

func (l *Lexer) hasModifier(identifier string) (*token.TokenType, *token.TokenModifier, string, bool) {
	if strings.Contains(identifier, ":") {
		parts := strings.Split(identifier, ":")
		var keywordType *token.TokenType
		if keyword, ok := l.isKeyword(parts[0]); ok {
			keywordType = &keyword
		}
		var modifierType *token.TokenModifier
		if modifier, ok := token.Modifiers[parts[1]]; ok {
			modifierType = &modifier
		}
		if keywordType != nil && modifierType != nil {
			return keywordType, modifierType, parts[0], true
		}
	}
	return nil, nil, "", false
}

func (l *Lexer) readIdentifier() string {
	for (isAlphaNumeric(l.peek()) || isAllowedIdentChar(l.peek())) && !l.isAtEnd() {
		l.nextChar()
	}
	text := l.Input[l.Start:l.Current]
	return text
}

type LexError struct {
	Message string
}

func (le *LexError) Error() string {
	return le.Message
}

func (l *Lexer) error(ch string, message string) *LexError {
	var where string
	if ch == "\n" {
		where = "at '\\n'"
	} else {
		where = "at '" + ch + "'"
	}
	err := &LexError{
		Message: fmt.Sprintf("[line %d] lex error %s: %s\n", l.Line, where, message),
	}
	return err
}

type Context struct {
	Stack []token.TokenType
}

func (c *Context) current() token.TokenType {
	if len(c.Stack) == 0 {
		return token.NONE
	}
	if len(c.Stack) == 1 {
		return c.Stack[0]
	}
	return c.Stack[len(c.Stack)-1] // last item
}

func (c *Context) setContext(t token.TokenType) {
	c.Stack = append(c.Stack, t)
}

func (c *Context) replaceContext(t token.TokenType) {
	c.resetContext()
	c.setContext(t)
}

// trims last item
func (c *Context) resetContext() {
	if len(c.Stack) > 0 {
		c.Stack = c.Stack[:len(c.Stack)-1]
	}
}

func TokenNames(types []token.Token) []string {
	var typeNames []string
	for _, t := range types {
		typeNames = append(typeNames, fmt.Sprintf("%s(%s)", token.TokenTypeNames[t.Type], t.Text))
	}
	return typeNames
}

func isDigit(ch string) bool {
	_, err := strconv.Atoi(ch)
	return err == nil
}

func isLetter(ch string) bool {
	for _, r := range ch {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func isAlphaNumeric(ch string) bool {
	return isDigit(ch) || isLetter(ch)
}

func isAllowedIdentChar(ch string) bool {
	allowed := map[string]bool{
		"_": true,
		"-": true,
		".": true,
		"/": true,
		"$": true,
		":": true,
	}
	return allowed[ch]
}

// helpful for creating parser tests
func TokenGenerator(input string) {
	lexer := New()
	tokens, err := lexer.ReadInput(input)
	if err != nil {
		fmt.Println(err)
	}
	for _, t := range tokens {
		fmt.Printf("{Type: token.%s, Text: \"%s\"},\n", token.TokenTypeNames[t.Type], t.Text)
	}
}
