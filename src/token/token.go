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

type Token struct {
	Type     TokenType
	Text     string
	Position int
	Line     int
}

func Keywords() map[string]TokenType {
	return map[string]TokenType{
		"var":    VAR,
		"target": TARGET,
		"run":    RUN,
	}
}
