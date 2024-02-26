package interpreter

import (
	"fmt"
	"os"
	"os/exec"
	"runny/src/env"
	"runny/src/tree"
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
		i.Environment.DefineVariable(variable.Name.Text, variable.Initialiser)
	}
	return nil
}

func (i *Interpreter) VisitTargetStatement(stmt tree.TargetStatement) interface{} {
	i.Environment.DefineTarget(stmt.Name.Text, stmt.Body)
	return nil
}

func (i *Interpreter) VisitActionStatement(stmt tree.ActionStatement) interface{} {
	if len(stmt.Body) == 0 {
		return nil
	}
	i.runShellCommand(stmt)
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

func (i *Interpreter) runShellCommand(statement tree.ActionStatement) {
	if len(statement.Body) == 0 {
		return
	}

	cmdString := statement.Body[0].Text
	for _, token := range statement.Body[1:] {
		cmdString += (" " + token.Text)
	}

	cmd := exec.Command("sh", "-c", cmdString)
	cmd.Env = os.Environ()

	variables := i.Environment.GetAll(env.VariableType)
	for name, value := range variables {
		cmd.Env = append(
			cmd.Env,
			fmt.Sprintf("%s=%s", name, value.Accept(i)),
		)
	}

	stdOutStdErr, _ := cmd.CombinedOutput()
	fmt.Print(string(stdOutStdErr))
}
