package token

type TokenType int

const (
	LEFT_BRACE TokenType = iota
	RIGHT_BRACE
	COLON
	COMMA

	IDENTIFIER
	ACTION // executable
	STRING
	NUMBER

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

	IDENTIFIER: "IDENTIFIER",
	ACTION:     "ACTION",
	STRING:     "STRING",
	NUMBER:     "NUMBER",

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
}
