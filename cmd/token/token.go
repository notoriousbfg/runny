package token

type TokenType int

const (
	LEFT_PAREN TokenType = iota
	RIGHT_PAREN
	COLON
	COMMA

	IDENTIFIER
	STRING
	NUMBER

	VAR
	VARS
	TARGET
	RUN

	NEWLINE
	EOF
)
