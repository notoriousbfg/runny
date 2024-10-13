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
	DESCRIBE

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

	VAR:      "VAR",
	TARGET:   "TARGET",
	RUN:      "RUN",
	CONFIG:   "CONFIG",
	EXTENDS:  "EXTENDS",
	DESCRIBE: "DESCRIBE",

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
	"desc":    DESCRIBE,
}

type TokenModifier int

const (
	BEFORE TokenModifier = iota
	AFTER
)

var TokenModifierNames = map[TokenModifier]string{
	BEFORE: "BEFORE",
	AFTER:  "AFTER",
}

var Modifiers = map[string]TokenModifier{
	"before": BEFORE,
	"after":  AFTER,
}

type Token struct {
	Type     TokenType
	Text     string
	Position int
	Line     int
	Depth    int
	Modifier *TokenModifier
}
