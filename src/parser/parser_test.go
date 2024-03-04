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
				tree.VariableDeclaration{
					Items: []tree.Variable{
						{
							Name: token.Token{Type: token.IDENTIFIER, Text: "name"},
							Initialiser: tree.ExpressionStatement{
								Expression: tree.Literal{Value: "Tim"},
							},
						},
					},
				},
			},
		},
		{
			name: "variable declaration with newline",
			tokens: []token.Token{
				{Type: token.VAR, Text: "var"},
				{Type: token.LEFT_BRACE, Text: "{"},
				{Type: token.IDENTIFIER, Text: "name"},
				{Type: token.STRING, Text: "Tim"},
				{Type: token.RIGHT_BRACE, Text: "}"},
				{Type: token.EOF, Text: ""},
			},
			want: []tree.Statement{
				tree.VariableDeclaration{
					Items: []tree.Variable{
						{
							Name: token.Token{Type: token.IDENTIFIER, Text: "name"},
							Initialiser: tree.ExpressionStatement{
								Expression: tree.Literal{Value: "Tim"},
							},
						},
					},
				},
			},
		},
		{
			name: "variable declaration block with newlines",
			tokens: []token.Token{
				{Type: token.VAR, Text: "var"},
				{Type: token.LEFT_BRACE, Text: "{"},
				{Type: token.IDENTIFIER, Text: "name"},
				{Type: token.LEFT_BRACE, Text: "{"},
				{Type: token.RUN, Text: "run"},
				{Type: token.LEFT_BRACE, Text: "{"},
				{Type: token.SCRIPT, Text: `echo "tim"`},
				{Type: token.RIGHT_BRACE, Text: "}"},
				{Type: token.RIGHT_BRACE, Text: "}"},
				{Type: token.RIGHT_BRACE, Text: "}"},
				{Type: token.EOF, Text: ""},
			},
			want: []tree.Statement{
				tree.VariableDeclaration{
					Items: []tree.Variable{
						{
							Name: token.Token{Type: token.IDENTIFIER, Text: "name"},
							Initialiser: tree.RunStatement{
								Body: []tree.Statement{
									tree.ActionStatement{
										Body: token.Token{Type: token.SCRIPT, Text: `echo "tim"`},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "variable declaration block",
			tokens: []token.Token{
				{Type: token.VAR, Text: "var"},
				{Type: token.LEFT_BRACE, Text: "{"},
				{Type: token.IDENTIFIER, Text: "name"},
				{Type: token.LEFT_BRACE, Text: "{"},
				{Type: token.RUN, Text: "run"},
				{Type: token.LEFT_BRACE, Text: "{"},
				{Type: token.SCRIPT, Text: `echo "tim"`},
				{Type: token.RIGHT_BRACE, Text: "}"},
				{Type: token.RIGHT_BRACE, Text: "}"},
				{Type: token.RIGHT_BRACE, Text: "}"},
				{Type: token.EOF, Text: ""},
			},
			want: []tree.Statement{
				tree.VariableDeclaration{
					Items: []tree.Variable{
						{
							Name: token.Token{Type: token.IDENTIFIER, Text: "name"},
							Initialiser: tree.RunStatement{
								Body: []tree.Statement{
									tree.ActionStatement{
										Body: token.Token{Type: token.SCRIPT, Text: `echo "tim"`},
									},
								},
							},
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
				tree.VariableDeclaration{
					Items: []tree.Variable{
						{
							Name: token.Token{Type: token.IDENTIFIER, Text: "name"},
							Initialiser: tree.ExpressionStatement{
								Expression: tree.Literal{Value: "Tim"},
							},
						},
						{
							Name: token.Token{Type: token.IDENTIFIER, Text: "foo"},
							Initialiser: tree.ExpressionStatement{
								Expression: tree.Literal{Value: "bar"},
							},
						},
					},
				},
			},
		},
		{
			name: "variable declaration (multiple) with newline",
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
				tree.VariableDeclaration{
					Items: []tree.Variable{
						{
							Name: token.Token{Type: token.IDENTIFIER, Text: "name"},
							Initialiser: tree.ExpressionStatement{
								Expression: tree.Literal{Value: "Tim"},
							},
						},
						{
							Name: token.Token{Type: token.IDENTIFIER, Text: "foo"},
							Initialiser: tree.ExpressionStatement{
								Expression: tree.Literal{Value: "bar"},
							},
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
				{Type: token.SCRIPT, Text: `echo "hello tim"`},
				{Type: token.RIGHT_BRACE, Text: "}"},
				{Type: token.EOF, Text: ""},
			},
			want: []tree.Statement{
				tree.TargetStatement{
					Name: token.Token{Type: token.IDENTIFIER, Text: "hello_cool_person"},
					Body: []tree.Statement{
						tree.ActionStatement{
							Body: token.Token{Type: token.SCRIPT, Text: `echo "hello tim"`},
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
				{Type: token.SCRIPT, Text: `echo "hello $name"`},
				{Type: token.RIGHT_BRACE, Text: "}"},
				{Type: token.EOF, Text: ""},
			},
			want: []tree.Statement{
				tree.TargetStatement{
					Name: token.Token{Type: token.IDENTIFIER, Text: "hello_cool_person"},
					Body: []tree.Statement{
						tree.VariableDeclaration{
							Items: []tree.Variable{
								{
									Name: token.Token{Type: token.IDENTIFIER, Text: "name"},
									Initialiser: tree.ExpressionStatement{
										Expression: tree.Literal{Value: "Tim"},
									},
								},
							},
						},
						tree.ActionStatement{
							Body: token.Token{Type: token.SCRIPT, Text: `echo "hello $name"`},
						},
					},
				},
			},
		},
		{
			name: "target declaration with run sandwiched between var declarations",
			tokens: []token.Token{
				{Type: token.TARGET, Text: "target"},
				{Type: token.IDENTIFIER, Text: "hello_cool_person"},
				{Type: token.LEFT_BRACE, Text: "{"},
				{Type: token.VAR, Text: "var"},
				{Type: token.LEFT_BRACE, Text: "{"},
				{Type: token.IDENTIFIER, Text: "name"},
				{Type: token.STRING, Text: "Tim"},
				{Type: token.RIGHT_BRACE, Text: "}"},
				{Type: token.RUN, Text: "run"},
				{Type: token.LEFT_BRACE, Text: "{"},
				{Type: token.SCRIPT, Text: `echo "hello $name"`},
				{Type: token.RIGHT_BRACE, Text: "}"},
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
						tree.VariableDeclaration{
							Items: []tree.Variable{
								{
									Name: token.Token{Type: token.IDENTIFIER, Text: "name"},
									Initialiser: tree.ExpressionStatement{
										Expression: tree.Literal{Value: "Tim"},
									},
								},
							},
						},
						tree.RunStatement{
							Body: []tree.Statement{
								tree.ActionStatement{
									Body: token.Token{Type: token.SCRIPT, Text: `echo "hello $name"`},
								},
							},
						},
						tree.VariableDeclaration{
							Items: []tree.Variable{
								{
									Name: token.Token{Type: token.IDENTIFIER, Text: "foo"},
									Initialiser: tree.ExpressionStatement{
										Expression: tree.Literal{Value: "bar"},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "run declaration with no target",
			tokens: []token.Token{
				{Type: token.RUN, Text: "run"},
				{Type: token.LEFT_BRACE, Text: "{"},
				{Type: token.SCRIPT, Text: `echo "hello"`},
				{Type: token.RIGHT_BRACE, Text: "}"},
				{Type: token.EOF, Text: ""},
			},
			want: []tree.Statement{
				tree.RunStatement{
					Body: []tree.Statement{
						tree.ActionStatement{
							Body: token.Token{Type: token.SCRIPT, Text: `echo "hello"`},
						},
					},
				},
			},
		},
		{
			name: "run declaration with target and no body",
			tokens: []token.Token{
				{Type: token.RUN, Text: "run"},
				{Type: token.IDENTIFIER, Text: "helloname"},
				{Type: token.EOF, Text: ""},
			},
			want: []tree.Statement{
				tree.RunStatement{
					Name: token.Token{Type: token.IDENTIFIER, Text: "helloname"},
					Body: []tree.Statement{},
				},
			},
		},
		{
			name: "run declaration with target and var declaration",
			tokens: []token.Token{
				{Type: token.RUN, Text: "run"},
				{Type: token.IDENTIFIER, Text: "helloname"},
				{Type: token.LEFT_BRACE, Text: "{"},
				{Type: token.VAR, Text: "var"},
				{Type: token.LEFT_BRACE, Text: "{"},
				{Type: token.IDENTIFIER, Text: "name"},
				{Type: token.STRING, Text: "tim"},
				{Type: token.RIGHT_BRACE, Text: "}"},
				{Type: token.RIGHT_BRACE, Text: "}"},
				{Type: token.EOF, Text: ""},
			},
			want: []tree.Statement{
				tree.RunStatement{
					Name: token.Token{Type: token.IDENTIFIER, Text: "helloname"},
					Body: []tree.Statement{
						tree.VariableDeclaration{
							Items: []tree.Variable{
								{
									Name: token.Token{Type: token.IDENTIFIER, Text: "name"},
									Initialiser: tree.ExpressionStatement{
										Expression: tree.Literal{Value: "tim"},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "blank newline between declarations",
			tokens: []token.Token{
				{Type: token.VAR, Text: "var"},
				{Type: token.LEFT_BRACE, Text: "{"},
				{Type: token.IDENTIFIER, Text: "name"},
				{Type: token.STRING, Text: "tim"},
				{Type: token.RIGHT_BRACE, Text: "}"},
				{Type: token.VAR, Text: "var"},
				{Type: token.LEFT_BRACE, Text: "{"},
				{Type: token.IDENTIFIER, Text: "foo"},
				{Type: token.STRING, Text: "bar"},
				{Type: token.RIGHT_BRACE, Text: "}"},
				{Type: token.EOF, Text: ""},
			},
			want: []tree.Statement{
				tree.VariableDeclaration{
					Items: []tree.Variable{
						{
							Name: token.Token{Type: token.IDENTIFIER, Text: "name"},
							Initialiser: tree.ExpressionStatement{
								Expression: tree.Literal{Value: "tim"},
							},
						},
					},
				},
				tree.VariableDeclaration{
					Items: []tree.Variable{
						{
							Name: token.Token{Type: token.IDENTIFIER, Text: "foo"},
							Initialiser: tree.ExpressionStatement{
								Expression: tree.Literal{Value: "bar"},
							},
						},
					},
				},
			},
		},
		{
			name: "variable variable",
			tokens: []token.Token{
				{Type: token.VAR, Text: "var"},
				{Type: token.LEFT_BRACE, Text: "{"},
				{Type: token.IDENTIFIER, Text: "name"},
				{Type: token.STRING, Text: "$tim"},
				{Type: token.RIGHT_BRACE, Text: "}"},
				{Type: token.EOF, Text: ""},
			},
			want: []tree.Statement{
				tree.VariableDeclaration{
					Items: []tree.Variable{
						{
							Name: token.Token{Type: token.IDENTIFIER, Text: "name"},
							Initialiser: tree.ExpressionStatement{
								Expression: tree.Literal{Value: "$tim"},
							},
						},
					},
				},
			},
		},
		{
			name: "run is not a keyword",
			tokens: []token.Token{
				{Type: token.RUN, Text: "run"},
				{Type: token.LEFT_BRACE, Text: "{"},
				{Type: token.SCRIPT, Text: `run something`},
				{Type: token.RIGHT_BRACE, Text: "}"},
				{Type: token.EOF, Text: ""},
			},
			want: []tree.Statement{
				tree.RunStatement{
					Body: []tree.Statement{
						tree.ActionStatement{
							Body: token.Token{
								Type: token.SCRIPT,
								Text: "run something",
							},
						},
					},
				},
			},
		},
		{
			name: "var declaration inside run target context",
			tokens: []token.Token{
				{Type: token.RUN, Text: "run"},
				{Type: token.IDENTIFIER, Text: "helloname"},
				{Type: token.LEFT_BRACE, Text: "{"},
				{Type: token.VAR, Text: "var"},
				{Type: token.LEFT_BRACE, Text: "{"},
				{Type: token.IDENTIFIER, Text: "name"},
				{Type: token.STRING, Text: "James"},
				{Type: token.RIGHT_BRACE, Text: "}"},
				{Type: token.RIGHT_BRACE, Text: "}"},
				{Type: token.EOF, Text: ""},
			},
			want: []tree.Statement{
				tree.RunStatement{
					Name: token.Token{Type: token.IDENTIFIER, Text: "helloname"},
					Body: []tree.Statement{
						tree.VariableDeclaration{
							Items: []tree.Variable{
								{
									Name: token.Token{Type: token.IDENTIFIER, Text: "name"},
									Initialiser: tree.ExpressionStatement{
										Expression: tree.Literal{Value: "James"},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "variable is not a string i.e. not wrapped in quotes",
			tokens: []token.Token{
				{Type: token.TARGET, Text: "target"},
				{Type: token.IDENTIFIER, Text: "build-lambdas"},
				{Type: token.LEFT_BRACE, Text: "{"},
				{Type: token.RUN, Text: "run"},
				{Type: token.IDENTIFIER, Text: "build-lambda-base"},
				{Type: token.RUN, Text: "run"},
				{Type: token.IDENTIFIER, Text: "build-lambda"},
				{Type: token.LEFT_BRACE, Text: "{"},
				{Type: token.VAR, Text: "var"},
				{Type: token.LEFT_BRACE, Text: "{"},
				{Type: token.IDENTIFIER, Text: "LAMBDA"},
				{Type: token.IDENTIFIER, Text: "list-id-providers"},
				{Type: token.RIGHT_BRACE, Text: "}"},
				{Type: token.RIGHT_BRACE, Text: "}"},
				{Type: token.RUN, Text: "run"},
				{Type: token.IDENTIFIER, Text: "build-lambda"},
				{Type: token.LEFT_BRACE, Text: "{"},
				{Type: token.VAR, Text: "var"},
				{Type: token.LEFT_BRACE, Text: "{"},
				{Type: token.IDENTIFIER, Text: "LAMBDA"},
				{Type: token.IDENTIFIER, Text: "post-auth"},
				{Type: token.RIGHT_BRACE, Text: "}"},
				{Type: token.RIGHT_BRACE, Text: "}"},
				{Type: token.RIGHT_BRACE, Text: "}"},
				{Type: token.EOF, Text: ""},
			},
			want: []tree.Statement{
				tree.TargetStatement{
					Name: token.Token{Type: token.IDENTIFIER, Text: "build-lambdas"},
					Body: []tree.Statement{
						tree.RunStatement{
							Name: token.Token{Type: token.IDENTIFIER, Text: "build-lambda-base"},
							Body: []tree.Statement{},
						},
						tree.RunStatement{
							Name: token.Token{Type: token.IDENTIFIER, Text: "build-lambda"},
							Body: []tree.Statement{
								tree.VariableDeclaration{
									Items: []tree.Variable{
										{
											Name: token.Token{Type: token.IDENTIFIER, Text: "LAMBDA"},
											Initialiser: tree.ExpressionStatement{
												Expression: tree.Literal{Value: "list-id-providers"},
											},
										},
									},
								},
							},
						},
						tree.RunStatement{
							Name: token.Token{Type: token.IDENTIFIER, Text: "build-lambda"},
							Body: []tree.Statement{
								tree.VariableDeclaration{
									Items: []tree.Variable{
										{
											Name: token.Token{Type: token.IDENTIFIER, Text: "LAMBDA"},
											Initialiser: tree.ExpressionStatement{
												Expression: tree.Literal{Value: "post-auth"},
											},
										},
									},
								},
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
				// wantJson, _ := json.Marshal(testcase.want)
				// stmtJson, _ := json.Marshal(p.Statements)
				// fmt.Println(string(wantJson))
				// fmt.Println(string(stmtJson))

				t.Fatalf("expressions do not match: expected: %+v, actual: %+v", testcase.want, p.Statements)
			}
		})
	}
}
