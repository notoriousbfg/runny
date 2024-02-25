package token

import "fmt"

type TokenType int

const (
	LEFT_BRACE TokenType = iota
	RIGHT_BRACE
	COLON
	COMMA
	BACKTICK

	IDENTIFIER
	STRING
	NUMBER
	OPERATOR
	FLAG
	COMMENT

	VAR
	TARGET
	RUN

	NEWLINE
	EOF
)

var TokenTypeNames = map[TokenType]string{
	LEFT_BRACE:  "LEFT_BRACE",
	RIGHT_BRACE: "RIGHT_BRACE",
	COLON:       "COLON",
	COMMA:       "COMMA",
	BACKTICK:    "BACKTICK",

	IDENTIFIER: "IDENTIFIER",
	STRING:     "STRING",
	NUMBER:     "NUMBER",
	OPERATOR:   "OPERATOR",
	FLAG:       "FLAG",
	COMMENT:    "COMMENT",

	VAR:    "VAR",
	TARGET: "TARGET",
	RUN:    "RUN",

	NEWLINE: "NEWLINE",
	EOF:     "EOF",
}

var Keywords = map[string]TokenType{
	"var":    VAR,
	"target": TARGET,
	"run":    RUN,
}

// type Context TokenType

// const (
// 	VAR_CTX    = Context(VAR)
// 	TARGET_CTX = Context(TARGET)
// 	RUN_CTX    = Context(RUN)
// )

type Token struct {
	Type     TokenType
	Text     string
	Position int
	Line     int
	Depth    int
}

func (t Token) String() string {
	switch {
	case t.Type == EOF:
		return "EOF"
	case len(t.Text) > 50:
		return fmt.Sprintf("%.10q...", t.Text)
	}
	return fmt.Sprintf("%q", t.Text)
}
