package lexer_test

import (
	"runny/src/lexer"
	"runny/src/token"
	"testing"
)

type TokenCase struct {
	inputString string
	types       []token.TokenType
	wantErr     bool
}

func TestLexer(t *testing.T) {
	cases := map[string]TokenCase{
		"basic: string variable declaration": {
			inputString: "var { hello \"world\" }",
			types: []token.TokenType{
				token.VAR,
				token.LEFT_BRACE,
				token.IDENTIFIER,
				token.STRING,
				token.RIGHT_BRACE,
				token.EOF,
			},
		},
		"basic: multiple string variable declarations": {
			inputString: "var { hello \"world\", name \"tim\" }",
			types: []token.TokenType{
				token.VAR,
				token.LEFT_BRACE,
				token.IDENTIFIER,
				token.STRING,
				token.COMMA,
				token.IDENTIFIER,
				token.STRING,
				token.RIGHT_BRACE,
				token.EOF,
			},
		},
		"basic: target": {
			inputString: "target { echo \"hello\" }",
			types: []token.TokenType{
				token.TARGET,
				token.LEFT_BRACE,
				token.IDENTIFIER,
				token.STRING,
				token.RIGHT_BRACE,
				token.EOF,
			},
		},
		"basic: nested declarations": {
			inputString: "target { echo \"hello\" var { name \"tim\" } echo \"tim\" }",
			types: []token.TokenType{
				token.TARGET,
				token.LEFT_BRACE,
				token.IDENTIFIER,
				token.STRING,
				token.VAR,
				token.LEFT_BRACE,
				token.IDENTIFIER,
				token.STRING,
				token.RIGHT_BRACE,
				token.IDENTIFIER,
				token.STRING,
				token.RIGHT_BRACE,
				token.EOF,
			},
		},
		"intermediate: multi line command": {
			inputString: `
				target build_container:private {
					docker build \
						-f .simulacrum/localstack/lambdas/$name.dockerfile \
						--build-arg $db_user \
						--build-arg $db_password \
						--build-arg $db_host \
						--build-arg $db_name \
						-t "$namespace:$name" \
						--no-cache \
						.
				}
			`,
			types: []token.TokenType{
				token.TARGET,
				token.LEFT_BRACE,
				token.IDENTIFIER,
				token.STRING,
				token.VAR,
				token.LEFT_BRACE,
				token.IDENTIFIER,
				token.STRING,
				token.RIGHT_BRACE,
				token.IDENTIFIER,
				token.STRING,
				token.RIGHT_BRACE,
				token.EOF,
			},
		},
	}

	for name, testcase := range cases {
		t.Run(name, func(t *testing.T) {
			l, err := lexer.New(testcase.inputString)
			if (err != nil) != testcase.wantErr {
				t.Fatalf("wantErr '%v', got '%+v', tokens: '%v'", testcase.wantErr, err, l.Tokens)
			}
			if !slicesMatch(l.TokenTypes(), testcase.types) {
				t.Fatal("types do not match", lexer.TokenTypeNames(testcase.types), lexer.TokenTypeNames(l.TokenTypes()))
			}
		})
	}
}

func slicesMatch(a []token.TokenType, b []token.TokenType) bool {
	if len(a) != len(b) {
		return false
	}

	for index, aType := range a {
		bType := b[index]
		if bType != aType {
			return false
		}
	}
	return true
}
