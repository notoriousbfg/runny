package parser_test

import (
	"encoding/json"
	"reflect"
	"runny/src/parser"
	"runny/src/token"
	"runny/src/tree"
	"testing"
)

type StatementCase struct {
	name    string
	tokens  func() []token.Token
	want    func() []tree.Statement
	wantErr string
}

func TestStatements(t *testing.T) {
	cases := []StatementCase{
		{
			name: "variable declaration",
			tokens: func() []token.Token {
				return []token.Token{
					{Type: token.VAR, Text: "var"},
					{Type: token.LEFT_BRACE, Text: "{"},
					{Type: token.IDENTIFIER, Text: "name"},
					{Type: token.STRING, Text: "Tim"},
					{Type: token.RIGHT_BRACE, Text: "}"},
					{Type: token.EOF, Text: ""},
				}
			},
			want: func() []tree.Statement {
				return []tree.Statement{
					tree.VariableStatement{
						Items: []tree.Variable{
							{
								Name: token.Token{Type: token.IDENTIFIER, Text: "name"},
								Initialiser: tree.ExpressionStatement{
									Expression: tree.Literal{Value: "Tim"},
								},
							},
						},
					},
				}
			},
		},
		{
			name: "variable declaration with newline",
			tokens: func() []token.Token {
				return []token.Token{
					{Type: token.VAR, Text: "var"},
					{Type: token.LEFT_BRACE, Text: "{"},
					{Type: token.IDENTIFIER, Text: "name"},
					{Type: token.STRING, Text: "Tim"},
					{Type: token.RIGHT_BRACE, Text: "}"},
					{Type: token.EOF, Text: ""},
				}
			},
			want: func() []tree.Statement {
				return []tree.Statement{
					tree.VariableStatement{
						Items: []tree.Variable{
							{
								Name: token.Token{Type: token.IDENTIFIER, Text: "name"},
								Initialiser: tree.ExpressionStatement{
									Expression: tree.Literal{Value: "Tim"},
								},
							},
						},
					},
				}
			},
		},
		{
			name: "variable declaration block with newlines",
			tokens: func() []token.Token {
				return []token.Token{
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
				}
			},
			want: func() []tree.Statement {
				return []tree.Statement{
					tree.VariableStatement{
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
				}
			},
		},
		{
			name: "variable declaration block",
			tokens: func() []token.Token {
				return []token.Token{
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
				}
			},
			want: func() []tree.Statement {
				return []tree.Statement{
					tree.VariableStatement{
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
				}
			},
		},
		{
			name: "variable declaration (multiple)",
			tokens: func() []token.Token {
				return []token.Token{
					{Type: token.VAR, Text: "var"},
					{Type: token.LEFT_BRACE, Text: "{"},
					{Type: token.IDENTIFIER, Text: "name"},
					{Type: token.STRING, Text: "Tim"},
					{Type: token.COMMA, Text: ","},
					{Type: token.IDENTIFIER, Text: "foo"},
					{Type: token.STRING, Text: "bar"},
					{Type: token.RIGHT_BRACE, Text: "}"},
					{Type: token.EOF, Text: ""},
				}
			},
			want: func() []tree.Statement {
				return []tree.Statement{
					tree.VariableStatement{
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
				}
			},
		},
		{
			name: "variable declaration (multiple) with newline",
			tokens: func() []token.Token {
				return []token.Token{
					{Type: token.VAR, Text: "var"},
					{Type: token.LEFT_BRACE, Text: "{"},

					{Type: token.IDENTIFIER, Text: "name"},
					{Type: token.STRING, Text: "Tim"},
					{Type: token.COMMA, Text: ","},
					{Type: token.IDENTIFIER, Text: "foo"},
					{Type: token.STRING, Text: "bar"},
					{Type: token.RIGHT_BRACE, Text: "}"},
					{Type: token.EOF, Text: ""},
				}
			},
			want: func() []tree.Statement {
				return []tree.Statement{
					tree.VariableStatement{
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
				}
			},
		},
		{
			name: "target declaration with single action",
			tokens: func() []token.Token {
				return []token.Token{
					{Type: token.TARGET, Text: "target"},
					{Type: token.IDENTIFIER, Text: "hello_cool_person"},
					{Type: token.LEFT_BRACE, Text: "{"},
					{Type: token.SCRIPT, Text: `echo "hello tim"`},
					{Type: token.RIGHT_BRACE, Text: "}"},
					{Type: token.EOF, Text: ""},
				}
			},
			want: func() []tree.Statement {
				return []tree.Statement{
					tree.TargetStatement{
						Name: token.Token{Type: token.IDENTIFIER, Text: "hello_cool_person"},
						Body: []tree.Statement{
							tree.ActionStatement{
								Body: token.Token{Type: token.SCRIPT, Text: `echo "hello tim"`},
							},
						},
					},
				}
			},
		},
		{
			name: "target declaration with var declaration and action",
			tokens: func() []token.Token {
				return []token.Token{
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
				}
			},
			want: func() []tree.Statement {
				return []tree.Statement{
					tree.TargetStatement{
						Name: token.Token{Type: token.IDENTIFIER, Text: "hello_cool_person"},
						Body: []tree.Statement{
							tree.VariableStatement{
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
				}
			},
		},
		{
			name: "target declaration with run sandwiched between var declarations",
			tokens: func() []token.Token {
				return []token.Token{
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
				}
			},
			want: func() []tree.Statement {
				return []tree.Statement{
					tree.TargetStatement{
						Name: token.Token{Type: token.IDENTIFIER, Text: "hello_cool_person"},
						Body: []tree.Statement{
							tree.VariableStatement{
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
							tree.VariableStatement{
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
				}
			},
		},
		{
			name: "run declaration with no target",
			tokens: func() []token.Token {
				return []token.Token{
					{Type: token.RUN, Text: "run"},
					{Type: token.LEFT_BRACE, Text: "{"},
					{Type: token.SCRIPT, Text: `echo "hello"`},
					{Type: token.RIGHT_BRACE, Text: "}"},
					{Type: token.EOF, Text: ""},
				}
			},
			want: func() []tree.Statement {
				return []tree.Statement{
					tree.RunStatement{
						Body: []tree.Statement{
							tree.ActionStatement{
								Body: token.Token{Type: token.SCRIPT, Text: `echo "hello"`},
							},
						},
					},
				}
			},
		},
		{
			name: "run declaration with target and no body",
			tokens: func() []token.Token {
				return []token.Token{
					{Type: token.RUN, Text: "run"},
					{Type: token.IDENTIFIER, Text: "helloname"},
					{Type: token.EOF, Text: ""},
				}
			},
			want: func() []tree.Statement {
				return []tree.Statement{
					tree.RunStatement{
						Name: token.Token{Type: token.IDENTIFIER, Text: "helloname"},
						Body: []tree.Statement{},
					},
				}
			},
		},
		{
			name: "run declaration with target and var declaration",
			tokens: func() []token.Token {
				return []token.Token{
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
				}
			},
			want: func() []tree.Statement {
				return []tree.Statement{
					tree.RunStatement{
						Name: token.Token{Type: token.IDENTIFIER, Text: "helloname"},
						Body: []tree.Statement{
							tree.VariableStatement{
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
				}
			},
		},
		{
			name: "blank newline between declarations",
			tokens: func() []token.Token {
				return []token.Token{
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
				}
			},
			want: func() []tree.Statement {
				return []tree.Statement{
					tree.VariableStatement{
						Items: []tree.Variable{
							{
								Name: token.Token{Type: token.IDENTIFIER, Text: "name"},
								Initialiser: tree.ExpressionStatement{
									Expression: tree.Literal{Value: "tim"},
								},
							},
						},
					},
					tree.VariableStatement{
						Items: []tree.Variable{
							{
								Name: token.Token{Type: token.IDENTIFIER, Text: "foo"},
								Initialiser: tree.ExpressionStatement{
									Expression: tree.Literal{Value: "bar"},
								},
							},
						},
					},
				}
			},
		},
		{
			name: "variable variable",
			tokens: func() []token.Token {
				return []token.Token{
					{Type: token.VAR, Text: "var"},
					{Type: token.LEFT_BRACE, Text: "{"},
					{Type: token.IDENTIFIER, Text: "name"},
					{Type: token.STRING, Text: "$tim"},
					{Type: token.RIGHT_BRACE, Text: "}"},
					{Type: token.EOF, Text: ""},
				}
			},
			want: func() []tree.Statement {
				return []tree.Statement{
					tree.VariableStatement{
						Items: []tree.Variable{
							{
								Name: token.Token{Type: token.IDENTIFIER, Text: "name"},
								Initialiser: tree.ExpressionStatement{
									Expression: tree.Literal{Value: "$tim"},
								},
							},
						},
					},
				}
			},
		},
		{
			name: "run is not a keyword",
			tokens: func() []token.Token {
				return []token.Token{
					{Type: token.RUN, Text: "run"},
					{Type: token.LEFT_BRACE, Text: "{"},
					{Type: token.SCRIPT, Text: `run something`},
					{Type: token.RIGHT_BRACE, Text: "}"},
					{Type: token.EOF, Text: ""},
				}
			},
			want: func() []tree.Statement {
				return []tree.Statement{
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
				}
			},
		},
		{
			name: "var declaration inside run target context",
			tokens: func() []token.Token {
				return []token.Token{
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
				}
			},
			want: func() []tree.Statement {
				return []tree.Statement{
					tree.RunStatement{
						Name: token.Token{Type: token.IDENTIFIER, Text: "helloname"},
						Body: []tree.Statement{
							tree.VariableStatement{
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
				}
			},
		},
		{
			name: "variable is not a string i.e. not wrapped in quotes",
			tokens: func() []token.Token {
				return []token.Token{
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
				}
			},
			want: func() []tree.Statement {
				return []tree.Statement{
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
									tree.VariableStatement{
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
									tree.VariableStatement{
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
				}
			},
		},
		{
			name: "config declaration",
			tokens: func() []token.Token {
				return []token.Token{
					{Type: token.CONFIG, Text: "config"},
					{Type: token.LEFT_BRACE, Text: "{"},
					{Type: token.IDENTIFIER, Text: "shell"},
					{Type: token.STRING, Text: "/bin/bash"},
					{Type: token.RIGHT_BRACE, Text: "}"},
					{Type: token.EOF, Text: ""},
				}
			},
			want: func() []tree.Statement {
				return []tree.Statement{
					tree.ConfigStatement{
						Items: []tree.Config{
							{
								Name: token.Token{
									Type: token.IDENTIFIER,
									Text: "shell",
								},
								Initialiser: tree.ExpressionStatement{
									Expression: tree.Literal{
										Value: "/bin/bash",
									},
								},
							},
						},
					},
				}
			},
		},
		{
			name: "extends declaration",
			tokens: func() []token.Token {
				return []token.Token{
					{Type: token.EXTENDS, Text: "extends"},
					{Type: token.LEFT_BRACE, Text: "{"},
					{Type: token.STRING, Text: "/some/path"},
					{Type: token.COMMA, Text: ","},
					{Type: token.STRING, Text: "/another/path"},
					{Type: token.RIGHT_BRACE, Text: "}"},
					{Type: token.EOF, Text: ""},
				}
			},
			want: func() []tree.Statement {
				return []tree.Statement{
					tree.ExtendsStatement{
						Paths: []tree.Expression{
							tree.Literal{
								Value: "/some/path",
							},
							tree.Literal{
								Value: "/another/path",
							},
						},
					},
				}
			},
		},
		{
			name: "run statement before stage",
			tokens: func() []token.Token {
				before := token.BEFORE
				return []token.Token{
					{Type: token.RUN, Text: "run:before", Modifier: &before},
					{Type: token.LEFT_BRACE, Text: "{"},
					{Type: token.SCRIPT, Text: `echo "before"`},
					{Type: token.RIGHT_BRACE, Text: "}"},
					{Type: token.EOF, Text: ""},
				}
			},
			want: func() []tree.Statement {
				return []tree.Statement{
					tree.RunStatement{
						Body: []tree.Statement{
							tree.ActionStatement{
								Body: token.Token{Type: token.SCRIPT, Text: `echo "before"`},
							},
						},
						Stage: tree.BEFORE,
					},
				}
			},
		},
	}

	for _, testcase := range cases {
		t.Run(testcase.name, func(t *testing.T) {
			p := parser.New()
			statements, err := p.Parse(testcase.tokens())
			if err != nil && err.Error() != testcase.wantErr {
				t.Fatalf("wantErr '%v', got '%+v', statements: '%v'", testcase.wantErr, err, statements)
			}
			if !reflect.DeepEqual(testcase.want(), statements) {
				t.Fatalf("expressions do not match: expected: %+v, actual: %+v", testcase.want(), statements)
			}
		})
	}
}

func TestParserErrors(t *testing.T) {
	cases := []StatementCase{
		{
			name: "variable declaration: missing identifier",
			tokens: func() []token.Token {
				return []token.Token{
					{Type: token.VAR, Text: "var"},
					{Type: token.LEFT_BRACE, Text: "{"},
					{Type: token.STRING, Text: "Tim"},
					{Type: token.RIGHT_BRACE, Text: "}"},
					{Type: token.EOF, Text: ""},
				}
			},
			wantErr: "[line 0] parse error at 'Tim': expect variable name\n",
			want: func() []tree.Statement {
				return nil
			},
		},
	}

	for _, testcase := range cases {
		t.Run(testcase.name, func(t *testing.T) {
			p := parser.New()
			statements, err := p.Parse(testcase.tokens())
			if err != nil && err.Error() != testcase.wantErr {
				t.Fatalf("wantErr '%v', got '%+v', statements: '%v'", testcase.wantErr, err, statements)
			}
			if !reflect.DeepEqual(testcase.want(), statements) {
				wantJson, _ := json.Marshal(testcase.want())
				gotJson, _ := json.Marshal(statements)
				t.Fatalf("expressions do not match: expected: %+v, actual: %+v", string(wantJson), string(gotJson))
			}
		})
	}
}
