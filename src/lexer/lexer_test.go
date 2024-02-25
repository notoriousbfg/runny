package lexer_test

import (
	"runny/src/lexer"
	"runny/src/token"
	"testing"
)

type TokenCase struct {
	name        string
	inputString string
	want        []token.Token
	wantErr     bool
}

func TestLexer(t *testing.T) {
	cases := []TokenCase{
		{
			name:        "basic: string variable declaration",
			inputString: "var { hello \"world\" }",
			want: []token.Token{
				{Type: token.VAR, Text: "var"},
				{Type: token.LEFT_BRACE, Text: "{"},
				{Type: token.IDENTIFIER, Text: "hello"},
				{Type: token.STRING, Text: "world"},
				{Type: token.RIGHT_BRACE, Text: "}"},
				{Type: token.EOF, Text: ""},
			},
		},
		{
			name:        "basic: multiple string variable declarations",
			inputString: "var { hello \"world\", name \"tim\" }",
			want: []token.Token{
				{Type: token.VAR, Text: "var"},
				{Type: token.LEFT_BRACE, Text: "{"},
				{Type: token.IDENTIFIER, Text: "hello"},
				{Type: token.STRING, Text: "world"},
				{Type: token.COMMA, Text: ","},
				{Type: token.IDENTIFIER, Text: "name"},
				{Type: token.STRING, Text: "tim"},
				{Type: token.RIGHT_BRACE, Text: "}"},
				{Type: token.EOF, Text: ""},
			},
		},
		{
			name:        "basic: target",
			inputString: "target { echo \"hello\" }",
			want: []token.Token{
				{Type: token.TARGET, Text: "target"},
				{Type: token.LEFT_BRACE, Text: "{"},
				{Type: token.IDENTIFIER, Text: "echo"},
				{Type: token.STRING, Text: "hello"},
				{Type: token.RIGHT_BRACE, Text: "}"},
				{Type: token.EOF, Text: ""},
			},
		},
		{
			name:        "basic: nested declarations",
			inputString: "target { echo \"hello\" var { name \"tim\" } echo \"tim\" }",
			want: []token.Token{
				{Type: token.TARGET, Text: "target"},
				{Type: token.LEFT_BRACE, Text: "{"},
				{Type: token.IDENTIFIER, Text: "echo"},
				{Type: token.STRING, Text: "hello"},
				{Type: token.VAR, Text: "var"},
				{Type: token.LEFT_BRACE, Text: "{"},
				{Type: token.IDENTIFIER, Text: "name"},
				{Type: token.STRING, Text: "tim"},
				{Type: token.RIGHT_BRACE, Text: "}"},
				{Type: token.IDENTIFIER, Text: "echo"},
				{Type: token.STRING, Text: "tim"},
				{Type: token.RIGHT_BRACE, Text: "}"},
				{Type: token.EOF, Text: ""},
			},
		},
		{
			name:        "intermediate: command with flag",
			inputString: "docker build -f dev.dockerfile",
			want: []token.Token{
				{Type: token.IDENTIFIER, Text: "docker"},
				{Type: token.IDENTIFIER, Text: "build"},
				{Type: token.FLAG, Text: "-f"},
				{Type: token.IDENTIFIER, Text: "dev.dockerfile"},
				{Type: token.EOF, Text: ""},
			},
		},
		{
			name:        "intermediate: command with double flag",
			inputString: "docker build --f dev.dockerfile",
			want: []token.Token{
				{Type: token.IDENTIFIER, Text: "docker"},
				{Type: token.IDENTIFIER, Text: "build"},
				{Type: token.FLAG, Text: "--f"},
				{Type: token.IDENTIFIER, Text: "dev.dockerfile"},
				{Type: token.EOF, Text: ""},
			},
		},
		{
			name: "intermediate: multi line command",
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
			want: []token.Token{
				{Type: token.TARGET, Text: "target"},
				{Type: token.IDENTIFIER, Text: "build_container"},
				{Type: token.COLON, Text: ":"},
				{Type: token.IDENTIFIER, Text: "private"},
				{Type: token.LEFT_BRACE, Text: "{"},
				{Type: token.IDENTIFIER, Text: "docker"},
				{Type: token.IDENTIFIER, Text: "build"},
				{Type: token.OPERATOR, Text: "\\"},
				{Type: token.FLAG, Text: "-f"},
				{Type: token.IDENTIFIER, Text: ".simulacrum/localstack/lambdas/$name.dockerfile"},
				{Type: token.OPERATOR, Text: "\\"},
				{Type: token.FLAG, Text: "--build-arg"},
				{Type: token.IDENTIFIER, Text: "$db_user"},
				{Type: token.OPERATOR, Text: "\\"},
				{Type: token.FLAG, Text: "--build-arg"},
				{Type: token.IDENTIFIER, Text: "$db_password"},
				{Type: token.OPERATOR, Text: "\\"},
				{Type: token.FLAG, Text: "--build-arg"},
				{Type: token.IDENTIFIER, Text: "$db_host"},
				{Type: token.OPERATOR, Text: "\\"},
				{Type: token.FLAG, Text: "--build-arg"},
				{Type: token.IDENTIFIER, Text: "$db_name"},
				{Type: token.OPERATOR, Text: "\\"},
				{Type: token.FLAG, Text: "-t"},
				{Type: token.STRING, Text: "$namespace:$name"},
				{Type: token.OPERATOR, Text: "\\"},
				{Type: token.FLAG, Text: "--no-cache"},
				{Type: token.OPERATOR, Text: "\\"},
				{Type: token.OPERATOR, Text: "."},
				{Type: token.RIGHT_BRACE, Text: "}"},
				{Type: token.EOF, Text: ""},
			},
		},
	}

	for _, testcase := range cases {
		t.Run(testcase.name, func(t *testing.T) {
			l, err := lexer.New(testcase.inputString)
			if (err != nil) != testcase.wantErr {
				t.Fatalf("wantErr '%v', got '%+v', tokens: '%v'", testcase.wantErr, err, l.Tokens)
			}
			if !tokenSlicesMatch(l.Tokens, testcase.want) {
				t.Fatal("types do not match", lexer.TokenNames(testcase.want), lexer.TokenNames(l.Tokens))
			}
		})
	}
}

func tokenSlicesMatch(a []token.Token, b []token.Token) bool {
	if len(a) != len(b) {
		return false
	}

	for index, aToken := range a {
		bToken := b[index]
		if bToken.Type != aToken.Type || bToken.Text != aToken.Text {
			return false
		}
	}
	return true
}
