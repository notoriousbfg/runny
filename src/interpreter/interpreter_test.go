package interpreter

import (
	"runny/src/token"
	"runny/src/tree"
	"testing"
)

const (
	origin = "/origin/path"
)

func TestInterpreter_VisitConfigStatement(t *testing.T) {
	t.Run("config variables are set", func(t *testing.T) {
		i := New(origin)
		i.VisitConfigStatement(tree.ConfigStatement{
			Items: []tree.Config{
				{
					Name: token.Token{Text: "shell", Type: token.STRING},
					Initialiser: tree.ExpressionStatement{
						Expression: tree.Literal{
							Value: "/bin/zsh",
						},
					},
				},
			},
		})
		_, exists := i.Config["shell"]
		if !exists {
			t.Errorf("config value was not set")
		}
	})
}
