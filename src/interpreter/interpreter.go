package interpreter

import (
	"fmt"
	"os"
	"os/exec"
	"runny/src/env"
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

func (i *Interpreter) VisitVariableStatement(stmt tree.VariableStatement) interface{} {
	for _, variable := range stmt.Items {
		// evaluate variable now and prevent infinite loop
		if action, ok := variable.Initialiser.(tree.ActionStatement); ok {
			stdout := i.runShellCommand(action.Body.Text, nil)
			variable.Initialiser = tree.ExpressionStatement{
				Expression: tree.Literal{Value: stdout},
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
	// if len(stmt.Body) == 0 {
	// 	return nil
	// }
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
	body := stmt.Body

	if stmt.Name != nil {
		targetBody, err := i.Environment.GetTarget(stmt.Name.Text)
		if err != nil {
			panic(err)
		}
		body = append(body, targetBody...)
	}

	for _, statement := range body {
		if _, ok := statement.(tree.VariableStatement); ok {
			i.Environment = env.NewEnvironment(i.Environment)
		}
		i.Accept(statement)
	}

	return nil
}

func (i *Interpreter) VisitExpressionStatement(stmt tree.ExpressionStatement) interface{} {
	return stmt.Expression.Accept(i)
}

func (i *Interpreter) VisitLiteralExpr(expr tree.Literal) interface{} {
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
