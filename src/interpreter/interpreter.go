package interpreter

import (
	"fmt"
	"os"
	"os/exec"
	"runny/src/env"
	"runny/src/token"
	"runny/src/tree"
	"strings"
)

func New(statements []tree.Statement) *Interpreter {
	return &Interpreter{
		Statements:  statements,
		Environment: env.NewEnvironment(nil),
	}
}

type Interpreter struct {
	Statements  []tree.Statement
	Environment *env.Environment
}

func (i *Interpreter) Evaluate() (result []interface{}) {
	for _, statement := range i.Statements {
		result = append(result, i.Accept(statement))
	}
	return
}

func (i *Interpreter) Accept(statement tree.Statement) interface{} {
	return statement.Accept(i)
}

func (i *Interpreter) VisitVariableDeclaration(stmt tree.VariableDeclaration) interface{} {
	for _, variable := range stmt.Items {
		// evaluate variable now and prevent infinite loop
		if run, ok := variable.Initialiser.(tree.RunStatement); ok {
			var strBuilder strings.Builder
			for _, action := range run.Body {
				if run, ok := action.(tree.ActionStatement); ok {
					stdout := i.runShellCommand(run.Body.Text, nil)
					strBuilder.Write(stdout)
				}
			}
			variable.Initialiser = tree.ExpressionStatement{
				Expression: tree.Literal{Value: strBuilder.String()},
			}
		}
		i.Environment.DefineVariable(variable.Name.Text, variable.Initialiser)
	}
	return nil
}

func (i *Interpreter) VisitTargetStatement(stmt tree.TargetStatement) interface{} {
	i.Environment.DefineTarget(stmt.Name.Text, stmt.Body)
	return nil
}

func (i *Interpreter) VisitActionStatement(stmt tree.ActionStatement) interface{} {
	evaluated := make(map[string]interface{}, 0)
	variables := i.Environment.GetAll(env.VariableType)
	for name, variable := range variables {
		evaluated[name] = i.Accept(variable)
	}
	bytes := i.runShellCommand(stmt.Body.Text, evaluated)
	fmt.Print(string(bytes))
	return nil
}

func (i *Interpreter) VisitRunStatement(stmt tree.RunStatement) interface{} {
	startEnvironment := i.Environment
	i.Environment = env.NewEnvironment(i.Environment)
	defer func() {
		i.Environment = startEnvironment
	}()

	body := stmt.Body

	if stmt.Name != (token.Token{}) {
		targetBody, err := i.Environment.GetTarget(stmt.Name.Text)
		if err != nil {
			panic(err)
		}
		// append contents of target onto end of body
		body = append(body, targetBody...)
	}

	for _, statement := range body {
		i.Accept(statement)
	}

	i.Environment = startEnvironment
	return nil
}

func (i *Interpreter) VisitExpressionStatement(stmt tree.ExpressionStatement) interface{} {
	return stmt.Expression.Accept(i)
}

// TODO
func (i *Interpreter) VisitVariableExpression(expr tree.VariableExpression) interface{} {
	// variable, _ := i.Environment.GetVariable(expr.Name.Text)
	// return variable
	return nil
}

func (i *Interpreter) VisitLiteralExpression(expr tree.Literal) interface{} {
	return expr.Value
}

func (i *Interpreter) runShellCommand(cmdString string, variables map[string]interface{}) []byte {
	if len(cmdString) == 0 {
		return []byte{}
	}

	cmd := exec.Command("sh", "-c", cmdString)
	cmd.Env = os.Environ()

	for name, value := range variables {
		cmd.Env = append(
			cmd.Env,
			fmt.Sprintf("%s=%s", name, trimQuotes(value)),
		)
	}

	stdOutStdErr, _ := cmd.CombinedOutput()
	return stdOutStdErr
}

func trimQuotes(input interface{}) interface{} {
	switch out := input.(type) {
	case string:
		return strings.Trim(out, "\"")
	}
	return input
}
