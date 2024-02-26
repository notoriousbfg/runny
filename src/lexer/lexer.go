package lexer

import (
	"fmt"
	"runny/src/token"
	"strconv"
	"unicode"
)

func New(input string) (*Lexer, error) {
	lexer := Lexer{
		Input:   input,
		Line:    1,
		Start:   0,
		Current: 0,
		Depth:   0,
	}
	err := lexer.readInput()
	if err != nil {
		return &lexer, err
	}
	return &lexer, nil
}

type Lexer struct {
	Input   string
	Tokens  []token.Token
	Start   int
	Current int
	Line    int
	Depth   int // number of braces deep
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
		l.Depth++
	case "}":
		l.addToken(token.RIGHT_BRACE, char)
		l.Depth--
	case ":":
		l.addToken(token.COLON, char)
	case ";":
		l.addToken(token.SEMICOLON, char)
	case ",":
		l.addToken(token.COMMA, char)
	case ".":
		if isAlphaNumeric(l.peek()) || isAllowedIdentChar(l.peek()) {
			l.matchIdentifier()
		} else {
			l.addToken(token.OPERATOR, char)
		}
	case "-":
		if l.matchNext("-") {
			if isAlphaNumeric(l.peek()) {
				l.matchFlag()
			} else {
				l.addToken(token.OPERATOR, "--")
			}
		} else if isAlphaNumeric(l.peek()) {
			l.matchFlag()
		} else {
			l.addToken(token.OPERATOR, "-")
		}
	case "+", "*", "/", "\\", "=":
		l.addToken(token.OPERATOR, char)
	case "$":
		l.matchIdentifier()
	case "\"", "`":
		l.matchString(char)
	case "\n":
		l.addToken(token.NEWLINE, "\\n")
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
		Depth:    l.Depth,
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

func (l *Lexer) matchString(delimiter string) {
	for l.peek() != delimiter && !l.isAtEnd() {
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
	if tokenType, keywordExists := token.Keywords[text]; keywordExists {
		l.addToken(tokenType, text)
	} else {
		l.addToken(token.IDENTIFIER, text)
	}
}

func (l *Lexer) matchFlag() {
	for (isLetter(l.peek()) || isHyphen(l.peek())) && !l.isAtEnd() {
		l.nextChar()
	}
	l.addToken(token.FLAG, l.Input[l.Start:l.Current])
}

func (l *Lexer) matchNext(expected string) bool {
	if string(l.Input[l.Current]) != expected {
		return false
	}
	l.nextChar()
	return true
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

// is this horrendously verbose?
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

func isHyphen(ch string) bool {
	return ch == "-"
}
