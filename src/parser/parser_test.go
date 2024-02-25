package parser_test

import (
	"reflect"
	"runny/src/parser"
	"runny/src/token"
	"runny/src/tree"
	"testing"
)

type StatementCase struct {
	name    string
	tokens  []token.Token
	want    []tree.Statement
	wantErr bool
}

func TestStatements(t *testing.T) {
	cases := []StatementCase{
		{
			name: "variable declaration",
			tokens: []token.Token{
				{Type: token.VAR, Text: "var"},
				{Type: token.LEFT_BRACE, Text: "{"},
				{Type: token.IDENTIFIER, Text: "name"},
				{Type: token.STRING, Text: "Tim"},
				{Type: token.RIGHT_BRACE, Text: "}"},
				{Type: token.EOF, Text: ""},
			},
			want: []tree.Statement{
				tree.VariableStatement{
					Items: []tree.Variable{
						{
							Name:        token.Token{Type: token.IDENTIFIER, Text: "name"},
							Initialiser: tree.Literal{Value: "Tim"},
						},
					},
				},
			},
		},
		{
			name: "variable declaration (multiple)",
			tokens: []token.Token{
				{Type: token.VAR, Text: "var"},
				{Type: token.LEFT_BRACE, Text: "{"},
				{Type: token.IDENTIFIER, Text: "name"},
				{Type: token.STRING, Text: "Tim"},
				{Type: token.COMMA, Text: ","},
				{Type: token.IDENTIFIER, Text: "foo"},
				{Type: token.STRING, Text: "bar"},
				{Type: token.RIGHT_BRACE, Text: "}"},
				{Type: token.EOF, Text: ""},
			},
			want: []tree.Statement{
				tree.VariableStatement{
					Items: []tree.Variable{
						{
							Name:        token.Token{Type: token.IDENTIFIER, Text: "name"},
							Initialiser: tree.Literal{Value: "Tim"},
						},
						{
							Name:        token.Token{Type: token.IDENTIFIER, Text: "foo"},
							Initialiser: tree.Literal{Value: "bar"},
						},
					},
				},
			},
		},
		{
			name: "target declaration with single action",
			tokens: []token.Token{
				{Type: token.TARGET, Text: "target"},
				{Type: token.IDENTIFIER, Text: "hello_cool_person"},
				{Type: token.LEFT_BRACE, Text: "{"},
				{Type: token.IDENTIFIER, Text: "echo"},
				{Type: token.STRING, Text: "hello tim"},
				{Type: token.RIGHT_BRACE, Text: "}"},
				{Type: token.EOF, Text: ""},
			},
			want: []tree.Statement{
				tree.TargetStatement{
					Name: token.Token{Type: token.IDENTIFIER, Text: "hello_cool_person"},
					Body: []tree.Statement{
						tree.ActionStatement{
							Body: []token.Token{
								{Type: token.IDENTIFIER, Text: "echo"},
								{Type: token.STRING, Text: "hello tim"},
							},
						},
					},
				},
			},
		},
		{
			name: "target declaration with var declaration and action",
			tokens: []token.Token{
				{Type: token.TARGET, Text: "target"},
				{Type: token.IDENTIFIER, Text: "hello_cool_person"},
				{Type: token.LEFT_BRACE, Text: "{"},
				{Type: token.VAR, Text: "var"},
				{Type: token.LEFT_BRACE, Text: "{"},
				{Type: token.IDENTIFIER, Text: "name"},
				{Type: token.STRING, Text: "Tim"},
				{Type: token.RIGHT_BRACE, Text: "}"},
				{Type: token.IDENTIFIER, Text: "echo"},
				{Type: token.STRING, Text: "hello $name"},
				{Type: token.RIGHT_BRACE, Text: "}"},
				{Type: token.EOF, Text: ""},
			},
			want: []tree.Statement{
				tree.TargetStatement{
					Name: token.Token{Type: token.IDENTIFIER, Text: "hello_cool_person"},
					Body: []tree.Statement{
						tree.VariableStatement{
							Items: []tree.Variable{
								{
									Name:        token.Token{Type: token.IDENTIFIER, Text: "name"},
									Initialiser: tree.Literal{Value: "Tim"},
								},
							},
						},
						tree.ActionStatement{
							Body: []token.Token{
								{Type: token.IDENTIFIER, Text: "echo"},
								{Type: token.STRING, Text: "hello $name"},
							},
						},
					},
				},
			},
		},
		{
			name: "target declaration with action sandwiched between var declarations",
			tokens: []token.Token{
				{Type: token.TARGET, Text: "target"},
				{Type: token.IDENTIFIER, Text: "hello_cool_person"},
				{Type: token.LEFT_BRACE, Text: "{"},
				{Type: token.VAR, Text: "var"},
				{Type: token.LEFT_BRACE, Text: "{"},
				{Type: token.IDENTIFIER, Text: "name"},
				{Type: token.STRING, Text: "Tim"},
				{Type: token.RIGHT_BRACE, Text: "}"},
				{Type: token.IDENTIFIER, Text: "echo"},
				{Type: token.STRING, Text: "hello $name"},
				{Type: token.VAR, Text: "var"},
				{Type: token.LEFT_BRACE, Text: "{"},
				{Type: token.IDENTIFIER, Text: "foo"},
				{Type: token.STRING, Text: "bar"},
				{Type: token.RIGHT_BRACE, Text: "}"},
				{Type: token.RIGHT_BRACE, Text: "}"},
				{Type: token.EOF, Text: ""},
			},
			want: []tree.Statement{
				tree.TargetStatement{
					Name: token.Token{Type: token.IDENTIFIER, Text: "hello_cool_person"},
					Body: []tree.Statement{
						tree.VariableStatement{
							Items: []tree.Variable{
								{
									Name:        token.Token{Type: token.IDENTIFIER, Text: "name"},
									Initialiser: tree.Literal{Value: "Tim"},
								},
							},
						},
						tree.ActionStatement{
							Body: []token.Token{
								{Type: token.IDENTIFIER, Text: "echo"},
								{Type: token.STRING, Text: "hello $name"},
							},
						},
						tree.VariableStatement{
							Items: []tree.Variable{
								{
									Name:        token.Token{Type: token.IDENTIFIER, Text: "foo"},
									Initialiser: tree.Literal{Value: "bar"},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "run declaration with no target and single action",
			tokens: []token.Token{
				{Type: token.RUN, Text: "run"},
				{Type: token.LEFT_BRACE, Text: "{"},
				{Type: token.IDENTIFIER, Text: "echo"},
				{Type: token.STRING, Text: "hello"},
				{Type: token.RIGHT_BRACE, Text: "}"},
				{Type: token.EOF, Text: ""},
			},
			want: []tree.Statement{
				tree.RunStatement{
					Name: nil,
					Body: []tree.Statement{
						tree.ActionStatement{
							Body: []token.Token{
								{Type: token.IDENTIFIER, Text: "echo"},
								{Type: token.STRING, Text: "hello"},
							},
						},
					},
				},
			},
		},
	}

	for _, testcase := range cases {
		t.Run(testcase.name, func(t *testing.T) {
			p := parser.New(testcase.tokens)
			err := p.Parse()
			if (err != nil) != testcase.wantErr {
				t.Fatalf("wantErr '%v', got '%+v', statements: '%v'", testcase.wantErr, err, p.Statements)
			}
			if !reflect.DeepEqual(testcase.want, p.Statements) {
				t.Fatalf("expressions do not match: expected: %+v, actual: %+v", testcase.want, p.Statements)
			}
		})
	}
}
