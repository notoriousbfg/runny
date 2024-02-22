package lexer

import (
	"fmt"
	"runny/src/token"
	"strconv"
	"unicode"
)

func New(input string) (Lexer, error) {
	lexer := Lexer{
		Input:   input,
		Line:    1,
		Start:   0,
		Current: 0,
	}
	err := lexer.readInput()
	if err != nil {
		return lexer, err
	}
	return lexer, nil
}

type Lexer struct {
	Input   string
	Tokens  []token.Token
	Start   int
	Current int
	Line    int
}

func (l *Lexer) readInput() error {
	for !l.isAtEnd() {
		l.Start = l.Current
		err := l.readChar()
		if err != nil {
			return err
		}
	}
	l.Start++
	l.addToken(token.EOF, "")
	return nil
}

func (l *Lexer) readChar() error {
	char := l.nextChar()
	switch char {
	case "{":
		l.addToken(token.LEFT_BRACE, char)
	case "}":
		l.addToken(token.RIGHT_BRACE, char)
	case ":":
		l.addToken(token.COLON, char)
	case ",":
		l.addToken(token.COMMA, char)
	case "\"":
		l.matchString()
	case "\n":
		l.Line++
	case " ", "\r", "\t":
		break
	default:
		if isDigit(char) {
			l.matchNumber()
		} else if isLetter(char) {
			l.matchIdentifier()
		} else {
			return fmt.Errorf("unsupported type: %s", char)
		}
	}
	return nil
}

func (l *Lexer) nextChar() string {
	char := string(l.Input[l.Current])
	l.Current++
	return char
}

func (l *Lexer) addToken(tokenType token.TokenType, text string) {
	l.Tokens = append(l.Tokens, token.Token{
		Type:     tokenType,
		Text:     text,
		Position: l.Start,
		Line:     l.Line,
	})
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

func (l *Lexer) peek() string {
	if l.isAtEnd() {
		return ""
	}
	return string(l.Input[l.Current])
}

func (l *Lexer) matchString() {
	for l.peek() != "\"" && !l.isAtEnd() {
		l.nextChar()
	}
	l.nextChar()
	text := l.Input[l.Start+1 : l.Current-1]
	l.addToken(token.STRING, text)
}

func (l *Lexer) matchNumber() {
	for isDigit(l.peek()) {
		l.nextChar()
	}

	// if l.peek() == "." && isDigit(l.peekNext()) {
	// 	l.nextChar()

	// 	for isDigit(l.peek()) {
	// 		l.nextChar()
	// 	}
	// }

	text := l.Input[l.Start:l.Current]

	// var val interface{}
	// if strings.Contains(text, ".") {
	// 	val, _ = strconv.ParseFloat(text, 64)
	// } else {
	// 	intVal, _ := strconv.ParseInt(text, 10, 0)
	// 	val = int(intVal)
	// }

	l.addToken(token.NUMBER, text)
}

// func (l *Lexer) matchAction() {
// 	for l.peek() != "}" && !l.isAtEnd() {
// 		l.nextChar()
// 	}
// 	text := l.Input[l.Start:l.Current]
// 	l.addToken(token.ACTION, text)
// }

func (l *Lexer) matchIdentifier() {
	for (isAlphaNumeric(l.peek()) || isAllowedIdentChar(l.peek())) && !l.isAtEnd() {
		l.nextChar()
	}

	text := l.Input[l.Start:l.Current]
	// var, target, run etc
	if tokenType, ok := token.Keywords[text]; ok {
		l.addToken(tokenType, text)
	} else {
		l.addToken(token.IDENTIFIER, text)
	}
}

func TokenTypeNames(types []token.TokenType) []string {
	var typeNames []string
	for _, t := range types {
		typeNames = append(typeNames, token.TokenTypeNames[t])
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

// is this horrendously verbose?
func isAllowedIdentChar(ch string) bool {
	allowed := map[string]bool{
		"_": true,
	}
	return allowed[ch]
}
