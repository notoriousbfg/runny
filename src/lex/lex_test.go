package lex_test

import (
	"runny/src/lex"
	"runny/src/token"
	"testing"
)

type TokenCase struct {
	name        string
	inputString string
	want        func() []token.Token
	wantErr     bool
}

func TestLexer(t *testing.T) {
	cases := []TokenCase{
		{
			name:        "basic: string variable declaration",
			inputString: "var { hello \"world\" }",
			want: func() []token.Token {
				return []token.Token{
					{Type: token.VAR, Text: "var"},
					{Type: token.LEFT_BRACE, Text: "{"},
					{Type: token.IDENTIFIER, Text: "hello"},
					{Type: token.STRING, Text: "\"world\""},
					{Type: token.RIGHT_BRACE, Text: "}"},
					{Type: token.EOF, Text: ""},
				}
			},
		},
		{
			name:        "basic: multiple string variable declarations",
			inputString: "var { hello \"world\", name \"tim\" }",
			want: func() []token.Token {
				return []token.Token{
					{Type: token.VAR, Text: "var"},
					{Type: token.LEFT_BRACE, Text: "{"},
					{Type: token.IDENTIFIER, Text: "hello"},
					{Type: token.STRING, Text: "\"world\""},
					{Type: token.COMMA, Text: ","},
					{Type: token.IDENTIFIER, Text: "name"},
					{Type: token.STRING, Text: "\"tim\""},
					{Type: token.RIGHT_BRACE, Text: "}"},
					{Type: token.EOF, Text: ""},
				}
			},
		},
		{
			name:        "basic: target",
			inputString: "target { run { echo \"hello\" } }",
			want: func() []token.Token {
				return []token.Token{
					{Type: token.TARGET, Text: "target"},
					{Type: token.LEFT_BRACE, Text: "{"},
					{Type: token.RUN, Text: "run"},
					{Type: token.LEFT_BRACE, Text: "{"},
					{Type: token.SCRIPT, Text: `echo "hello"`},
					{Type: token.RIGHT_BRACE, Text: "}"},
					{Type: token.RIGHT_BRACE, Text: "}"},
					{Type: token.EOF, Text: ""},
				}
			},
		},
		{
			name:        "basic: nested declarations",
			inputString: "target { run { echo \"hello\" } var { name \"tim\" } run { echo $name } }",
			want: func() []token.Token {
				return []token.Token{
					{Type: token.TARGET, Text: "target"},
					{Type: token.LEFT_BRACE, Text: "{"},
					{Type: token.RUN, Text: "run"},
					{Type: token.LEFT_BRACE, Text: "{"},
					{Type: token.SCRIPT, Text: `echo "hello"`},
					{Type: token.RIGHT_BRACE, Text: "}"},
					{Type: token.VAR, Text: "var"},
					{Type: token.LEFT_BRACE, Text: "{"},
					{Type: token.IDENTIFIER, Text: "name"},
					{Type: token.STRING, Text: "\"tim\""},
					{Type: token.RIGHT_BRACE, Text: "}"},
					{Type: token.RUN, Text: "run"},
					{Type: token.LEFT_BRACE, Text: "{"},
					{Type: token.SCRIPT, Text: `echo $name`},
					{Type: token.RIGHT_BRACE, Text: "}"},
					{Type: token.RIGHT_BRACE, Text: "}"},
					{Type: token.EOF, Text: ""},
				}
			},
		},
		{
			name:        "basic: run declarations",
			inputString: "run { echo \"hello\" echo \"tim\" }",
			want: func() []token.Token {
				return []token.Token{
					{Type: token.RUN, Text: "run"},
					{Type: token.LEFT_BRACE, Text: "{"},
					{Type: token.SCRIPT, Text: `echo "hello" echo "tim"`},
					{Type: token.RIGHT_BRACE, Text: "}"},
					{Type: token.EOF, Text: ""},
				}
			},
		},
		// {
		// 	name: "intermediate: multi line command",
		// 	inputString: `target build_container:private {
		// 		run {
		// 			docker build \ -f .simulacrum/localstack/lambdas/$name.dockerfile \ --build-arg $db_user \ --build-arg $db_password \ --build-arg $db_host \ --build-arg $db_name \ -t "$namespace:$name" \ --no-cache \ .
		// 		}
		// 	}`,
		// 	want: []token.Token{
		// 		{Type: token.TARGET, Text: "target"},
		// 		{Type: token.IDENTIFIER, Text: "build_container:private"},
		// 		{Type: token.LEFT_BRACE, Text: "{"},
		// 		{Type: token.RUN, Text: "run"},
		// 		{Type: token.LEFT_BRACE, Text: "{"},
		// 		{Type: token.SCRIPT, Text: `docker build \ -f .simulacrum/localstack/lambdas/$name.dockerfile \ --build-arg $db_user \ --build-arg $db_password \ --build-arg $db_host \ --build-arg $db_name \ -t "$namespace:$name" \ --no-cache \ .`},
		// 		{Type: token.RIGHT_BRACE, Text: "}"},
		// 		{Type: token.RIGHT_BRACE, Text: "}"},
		// 		{Type: token.EOF, Text: ""},
		// 	},
		// },
		{
			name:        "intermediate: keyword inside braces",
			inputString: "run { run something }",
			want: func() []token.Token {
				return []token.Token{
					{Type: token.RUN, Text: "run"},
					{Type: token.LEFT_BRACE, Text: "{"},
					{Type: token.SCRIPT, Text: "run something"},
					{Type: token.RIGHT_BRACE, Text: "}"},
					{Type: token.EOF, Text: ""},
				}
			},
		},
		{
			name:        "intermediate: run target with no body",
			inputString: "run mytarget {}",
			want: func() []token.Token {
				return []token.Token{
					{Type: token.RUN, Text: "run"},
					{Type: token.IDENTIFIER, Text: "mytarget"},
					{Type: token.LEFT_BRACE, Text: "{"},
					{Type: token.RIGHT_BRACE, Text: "}"},
					{Type: token.EOF, Text: ""},
				}
			},
		},
		{
			name:        "intermediate: run target with var declaration",
			inputString: `run mytarget { var { name "tim" } }`,
			want: func() []token.Token {
				return []token.Token{
					{Type: token.RUN, Text: "run"},
					{Type: token.IDENTIFIER, Text: "mytarget"},
					{Type: token.LEFT_BRACE, Text: "{"},
					{Type: token.VAR, Text: "var"},
					{Type: token.LEFT_BRACE, Text: "{"},
					{Type: token.IDENTIFIER, Text: "name"},
					{Type: token.STRING, Text: "\"tim\""},
					{Type: token.RIGHT_BRACE, Text: "}"},
					{Type: token.RIGHT_BRACE, Text: "}"},
					{Type: token.EOF, Text: ""},
				}
			},
		},
		{
			name:        "intermediate: multiple successive variable scripts",
			inputString: `var { one { run { echo "one" } } two { run { echo "two" } } }`,
			want: func() []token.Token {
				return []token.Token{
					{Type: token.VAR, Text: "var"},
					{Type: token.LEFT_BRACE, Text: "{"},
					{Type: token.IDENTIFIER, Text: "one"},
					{Type: token.LEFT_BRACE, Text: "{"},
					{Type: token.RUN, Text: "run"},
					{Type: token.LEFT_BRACE, Text: "{"},
					{Type: token.SCRIPT, Text: "echo \"one\""},
					{Type: token.RIGHT_BRACE, Text: "}"},
					{Type: token.RIGHT_BRACE, Text: "}"},
					{Type: token.IDENTIFIER, Text: "two"},
					{Type: token.LEFT_BRACE, Text: "{"},
					{Type: token.RUN, Text: "run"},
					{Type: token.LEFT_BRACE, Text: "{"},
					{Type: token.SCRIPT, Text: "echo \"two\""},
					{Type: token.RIGHT_BRACE, Text: "}"},
					{Type: token.RIGHT_BRACE, Text: "}"},
					{Type: token.RIGHT_BRACE, Text: "}"},
					{Type: token.EOF, Text: ""},
				}
			},
		},
		// {
		// 	name:        "intermediate: string containing other strings",
		// 	inputString: "run { `docker run -d --name \"my-container\" MYSQL_ROOT_PASSWORD=$mysql_root_password` }",
		// 	want: []token.Token{
		// 		{Type: token.RUN, Text: "run"},
		// 		{Type: token.LEFT_BRACE, Text: "{"},
		// 		{Type: token.STRING, Text: "\"docker run -d --name \"my-container\" MYSQL_ROOT_PASSWORD=$mysql_root_password\""},
		// 		{Type: token.RIGHT_BRACE, Text: "}"},
		// 		{Type: token.EOF, Text: ""},
		// 	},
		// },
		{
			name:        "basic: run modifier",
			inputString: `run:before { echo "before" }`,
			want: func() []token.Token {
				before := token.BEFORE
				return []token.Token{
					{Type: token.RUN, Text: "run:before", Modifier: &before},
					{Type: token.LEFT_BRACE, Text: "{"},
					{Type: token.SCRIPT, Text: `echo "before"`},
					{Type: token.RIGHT_BRACE, Text: "}"},
					{Type: token.EOF, Text: ""},
				}
			},
		},
		{
			name:        "basic: var modifier",
			inputString: `var:before { name "Jack" }`,
			want: func() []token.Token {
				before := token.BEFORE
				return []token.Token{
					{Type: token.VAR, Text: "var:before", Modifier: &before},
					{Type: token.LEFT_BRACE, Text: "{"},
					{Type: token.IDENTIFIER, Text: "name"},
					{Type: token.STRING, Text: `"Jack"`},
					{Type: token.RIGHT_BRACE, Text: "}"},
					{Type: token.EOF, Text: ""},
				}
			},
		},
	}

	for _, testcase := range cases {
		t.Run(testcase.name, func(t *testing.T) {
			l := lex.New()
			_, err := l.ReadInput(testcase.inputString)
			if (err != nil) != testcase.wantErr {
				t.Fatalf("wantErr '%v', got '%+v', tokens: '%v'", testcase.wantErr, err, l.Tokens)
			}
			if !tokenSlicesMatch(l.Tokens, testcase.want()) {
				t.Fatal("types do not match", lex.TokenNames(testcase.want()), lex.TokenNames(l.Tokens))
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
