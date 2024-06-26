package token

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
	CONFIG
	EXTENDS

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

	VAR:     "VAR",
	TARGET:  "TARGET",
	RUN:     "RUN",
	CONFIG:  "CONFIG",
	EXTENDS: "EXTENDS",

	NEWLINE: "NEWLINE",
	NONE:    "NONE",
	EOF:     "EOF",
}

var Keywords = map[string]TokenType{
	"var":     VAR,
	"target":  TARGET,
	"run":     RUN,
	"config":  CONFIG,
	"extends": EXTENDS,
}

type Token struct {
	Type     TokenType
	Text     string
	Position int
	Line     int
	Depth    int
}
