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
		Environment: *env.NewEnvironment(nil),
	}
}

type Interpreter struct {
	Statements  []tree.Statement
	Environment env.Environment
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
		i.Environment.Define(variable.Name.Text, variable.Initialiser.Accept(i), env.VariableType)
	}

	return nil
}

func (i *Interpreter) VisitTargetStatement(stmt tree.TargetStatement) interface{} {
	for _, target := range stmt.Body {
		i.Environment.Define(stmt.Name.Text, target.Accept(i), env.TargetType)
	}
	return nil
}

func (i *Interpreter) VisitActionStatement(stmt tree.ActionStatement) interface{} {
	return nil
}

func (i *Interpreter) VisitRunStatement(stmt tree.RunStatement) interface{} {
	// TODO: running specific target
	for _, statement := range stmt.Body {
		switch statementType := statement.(type) {
		case tree.VariableStatement:
			// define scoped vars
		case tree.ActionStatement:
			i.runShellCommand(statementType)
		case tree.RunStatement:
			i.VisitRunStatement(statementType)
		}
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
			fmt.Sprintf("%s=%s", name, value),
		)
	}

	stdOutStdErr, _ := cmd.CombinedOutput()
	fmt.Println(string(stdOutStdErr))
}
