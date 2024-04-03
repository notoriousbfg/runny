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
		Config:      make(map[string]interface{}, 0),
		Environment: env.NewEnvironment(nil),
	}
}

type Config map[string]interface{}

func (c Config) getShell() string {
	shell, ok := c["shell"]
	if ok {
		trimmedShell := trimQuotes(shell)
		if shellStr, ok := trimmedShell.(string); ok {
			return shellStr
		}
	}
	return "sh"
}

type Interpreter struct {
	Config      Config
	Extends     []string
	Statements  []tree.Statement
	Environment *env.Environment
}

func (i *Interpreter) Evaluate(statements []tree.Statement) (result []interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			if str, ok := r.(string); ok {
				err = fmt.Errorf(str)
			} else if e, ok := r.(error); ok {
				err = e
			} else {
				err = fmt.Errorf("unknown panic: %v", r)
			}
		}
	}()
	for _, statement := range statements {
		result = append(result, i.Accept(statement))
	}
	return
}

func (i *Interpreter) Accept(statement tree.Statement) interface{} {
	return statement.Accept(i)
}

func (i *Interpreter) VisitConfigStatement(stmt tree.ConfigStatement) interface{} {
	for _, config := range stmt.Items {
		i.Config[config.Name.Text] = i.Accept(config.Initialiser)
	}
	return nil
}

func (i *Interpreter) VisitVariableStatement(stmt tree.VariableStatement) interface{} {
	for _, variable := range stmt.Items {
		i.Environment.Define(variable.Name.Text, env.VTVar, variable.Initialiser)
	}
	return nil
}

func (i *Interpreter) VisitTargetStatement(stmt tree.TargetStatement) interface{} {
	i.Environment.Define(stmt.Name.Text, env.VTTarget, stmt.Body)
	return nil
}

func (i *Interpreter) VisitActionStatement(stmt tree.ActionStatement) interface{} {
	evaluated := make(map[string]interface{}, 0)
	for k := range i.Environment.GetAll(env.VTVar) {
		variable, _ := i.lookupVariable(k)
		evaluated[k] = variable
	}
	bytes := runShellCommand(stmt.Body.Text, evaluated, i.Config.getShell())
	fmt.Println("\033[0m" + stmt.Body.Text)
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
		targetBodyInt, err := i.Environment.Get(stmt.Name.Text, env.VTTarget)
		if err != nil {
			panic(i.error(err.Error()))
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

func (i *Interpreter) VisitExtendsStatement(stmt tree.ExtendsStatement) interface{} {
	for _, path := range stmt.Paths {
		evaluatedPath := path.Accept(i)
		if pathStr, isString := evaluatedPath.(string); isString {
			i.Extends = append(i.Extends, pathStr)
			// TODO: ingest file
		} else {
			panic(i.error(fmt.Sprintf("extends path %v is not a string", evaluatedPath)))
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

func (i *Interpreter) lookupVariable(name string) (interface{}, error) {
	variable, err := i.Environment.Get(name, env.VTVar)
	if err != nil {
		return nil, err
	}
	switch typedVal := variable.(type) {
	// if run statement evaluate its actions now
	case tree.RunStatement:
		var strBuilder strings.Builder
		for _, action := range typedVal.Body {
			if run, ok := action.(tree.ActionStatement); ok {
				stdout := runShellCommand(run.Body.Text, nil, i.Config.getShell())
				strBuilder.Write(stdout)
			}
		}
		return strBuilder.String(), nil
	case tree.Statement:
		return i.Accept(typedVal), nil
	default:
		return "", nil
	}
}

type RuntimeError struct {
	Message string
}

func (re *RuntimeError) Error() string {
	return re.Message
}

func (i *Interpreter) error(message string) *RuntimeError {
	err := &RuntimeError{
		Message: fmt.Sprintf("runtime error: %s\n", message),
	}
	return err
}

func runShellCommand(cmdString string, variables map[string]interface{}, shell string) []byte {
	if len(cmdString) == 0 {
		return []byte{}
	}
	cmd := exec.Command(shell, "-c", cmdString)
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
