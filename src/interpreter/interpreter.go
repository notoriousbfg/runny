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
		Environment: env.NewEnvironment(nil),
		locals:      make(locals),
	}
}

type locals map[string]int

func (l locals) put(expr tree.Expression, depth int) {
	l[l.key(expr)] = depth
}

func (l locals) get(expr tree.Expression) (int, bool) {
	distance, ok := l[l.key(expr)]
	return distance, ok
}

func (l *locals) key(expr tree.Expression) string {
	return fmt.Sprintf("%#v", expr)
}

type Interpreter struct {
	Statements  []tree.Statement
	Environment *env.Environment
	locals      locals
}

func (i *Interpreter) Evaluate(statements []tree.Statement) (result []interface{}) {
	for _, statement := range statements {
		result = append(result, i.Accept(statement))
	}
	return
}

func (i *Interpreter) Accept(statement tree.Statement) interface{} {
	return statement.Accept(i)
}

func (i *Interpreter) VisitVariableStatement(stmt tree.VariableStatement) interface{} {
	for _, variable := range stmt.Items {
		i.Environment.Define(variable.Name.Text, variable.Initialiser)
	}
	return nil
}

func (i *Interpreter) VisitTargetStatement(stmt tree.TargetStatement) interface{} {
	i.Environment.Define(stmt.Name.Text, stmt.Body)
	return nil
}

func (i *Interpreter) VisitActionStatement(stmt tree.ActionStatement) interface{} {
	evaluated := make(map[string]interface{}, 0)
	for _, variable := range stmt.Variables {
		evaluated[variable.Name.Text] = i.VisitVariableExpr(variable)
	}

	bytes := runShellCommand(stmt.Body.Text, evaluated)
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
		targetBodyInt, err := i.Environment.Get(stmt.Name.Text)
		if err != nil {
			panic(err)
		}
		if targetBody, ok := targetBodyInt.([]tree.Statement); ok {
			// append contents of target onto end of body
			body = append(body, targetBody...)
		}
	}

	for _, statement := range body {
		i.Accept(statement)
	}

	return nil
}

func (i *Interpreter) VisitExpressionStatement(stmt tree.ExpressionStatement) interface{} {
	return stmt.Expression.Accept(i)
}

func (i *Interpreter) VisitVariableExpr(expr tree.VariableExpression) interface{} {
	val, err := i.lookupVariable(expr.Name.Text, expr)
	if err != nil {
		panic(err)
	}
	switch typedVal := val.(type) {
	// if run statement evaluate its actions now
	case tree.RunStatement:
		var strBuilder strings.Builder
		for _, action := range typedVal.Body {
			if run, ok := action.(tree.ActionStatement); ok {
				stdout := runShellCommand(run.Body.Text, nil)
				strBuilder.Write(stdout)
			}
		}
		return strBuilder.String()
	case tree.Statement:
		return i.Accept(typedVal)
	default:
		return ""
	}
}

func (i *Interpreter) VisitLiteralExpr(expr tree.Literal) interface{} {
	return expr.Value
}

func (i *Interpreter) Resolve(expr tree.Expression, depth int) {
	i.locals.put(expr, depth)
}

// this isn't working properly
func (i *Interpreter) lookupVariable(name string, expr tree.Expression) (interface{}, error) {
	// if distance, ok := i.locals.get(expr); ok {
	// 	return i.Environment.GetAt(distance, name), nil
	// }
	return i.Environment.Get(name)
}

func runShellCommand(cmdString string, variables map[string]interface{}) []byte {
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
