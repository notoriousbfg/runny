package token

import "fmt"

type TokenType int

const (
	LEFT_BRACE TokenType = iota
	RIGHT_BRACE
	COMMA

	IDENTIFIER
	STRING
	NUMBER
	COMMENT
	SCRIPT

	VAR
	TARGET
	RUN

	NEWLINE
	NONE
	EOF
)

var TokenTypeNames = map[TokenType]string{
	LEFT_BRACE:  "LEFT_BRACE",
	RIGHT_BRACE: "RIGHT_BRACE",
	COMMA:       "COMMA",

	IDENTIFIER: "IDENTIFIER",
	STRING:     "STRING",
	NUMBER:     "NUMBER",
	COMMENT:    "COMMENT",
	SCRIPT:     "SCRIPT",

	VAR:    "VAR",
	TARGET: "TARGET",
	RUN:    "RUN",

	NEWLINE: "NEWLINE",
	NONE:    "NONE",
	EOF:     "EOF",
}

var Keywords = map[string]TokenType{
	"var":    VAR,
	"target": TARGET,
	"run":    RUN,
}

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
