package parser_test

import (
	"reflect"
	"runny/src/parser"
	"runny/src/token"
	"runny/src/tree"
	"testing"
)

type StatementCase struct {
	tokens []token.Token
	want   []tree.Statement
}

func TestStatements(t *testing.T) {
	cases := map[string]StatementCase{
		"variable declaration": {
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
		"variable declaration (multiple)": {
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
	}

	for name, testcase := range cases {
		t.Run(name, func(t *testing.T) {
			p := parser.New(testcase.tokens)
			parsedExpression := p.Parse()
			if !reflect.DeepEqual(testcase.want, parsedExpression) {
				t.Fatalf("expressions do not match: expected: %+v, actual: %+v", testcase.want, parsedExpression)
			}
		})
	}
}
