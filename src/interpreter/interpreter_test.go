package interpreter

import (
	"runny/src/token"
	"runny/src/tree"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	origin = "/origin/path"
)

// func TestInterpreter_VisitRunStatement(t *testing.T) {
// 	t.Run("simple command printed & executed accurately", func(t *testing.T) {
// 		i := New(origin, true)
// 		output, _ := captureOutput(func() error {
// 			i.VisitRunStatement(tree.RunStatement{
// 				Body: []tree.Statement{
// 					tree.ActionStatement{
// 						Body: token.Token{
// 							Text: "echo \"hello world\"",
// 						},
// 					},
// 				},
// 			})
// 			return nil
// 		})
// 		// hacky
// 		expectedOutput := fmt.Sprintf(`%secho "hello world"%s
// hello world
// `, foreColour, aftColour)
// 		assert.Equal(t, expectedOutput, output)
// 	})
// }

func TestInterpreter_VisitConfigStatement(t *testing.T) {
	t.Run("config variables are set", func(t *testing.T) {
		i := New(origin, true)
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
		assert.NotEmpty(t, i.Config["shell"], "config value was not set")
		assert.Equal(t, i.Config["shell"], "/bin/zsh")
	})
}

// func TestInterpreter_VisitDescribeStatement(t *testing.T) {
// 	t.Run("description statement printed", func(t *testing.T) {
// 		i := New(origin, true)
// 		output, _ := captureOutput(func() error {
// 			i.VisitDescribeStatement(tree.DescribeStatement{
// 				Lines: []tree.Literal{
// 					{
// 						Value: "the command does X",
// 					},
// 				},
// 			})
// 			return nil
// 		})
// 		assert.Equal(t, "> the command does X\n", output)
// 	})
// 	t.Run("multiple description statements printed", func(t *testing.T) {
// 		i := New(origin, true)
// 		output, _ := captureOutput(func() error {
// 			i.VisitDescribeStatement(tree.DescribeStatement{
// 				Lines: []tree.Literal{
// 					{
// 						Value: "the command does X",
// 					},
// 					{
// 						Value: "the command does Y",
// 					},
// 				},
// 			})
// 			return nil
// 		})
// 		assert.Equal(t, "> the command does X\n> the command does Y\n", output)
// 	})
// }

// func captureOutput(f func() error) (string, error) {
// 	orig := os.Stdout
// 	r, w, _ := os.Pipe()
// 	os.Stdout = w
// 	err := f()
// 	os.Stdout = orig
// 	w.Close()
// 	out, _ := io.ReadAll(r)
// 	return string(out), err
// }
