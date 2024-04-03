package interpreter

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runny/src/env"
	"runny/src/lex"
	"runny/src/parser"
	"runny/src/token"
	"runny/src/tree"
	"strings"
)

func New(origin string) *Interpreter {
	return &Interpreter{
		Config:      make(map[string]interface{}, 0),
		Environment: env.NewEnvironment(nil),
		Origin:      origin,
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
	Statements  []tree.Statement
	Origin      string // the file path currently being read from
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
	i.Statements = statements
	for _, statement := range i.Statements {
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
	fmt.Println("\033[32m" + stmt.Body.Text + "\033[0m")
	err := runShellCommandAndPipeToStdout(stmt.Body.Text, evaluated, i.Config.getShell())
	if err != nil {
		panic(i.error(fmt.Sprintf("could not run command: %s", err)))
	}
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
		evaluatedPath = trimQuotes(evaluatedPath)
		if pathStr, isString := evaluatedPath.(string); isString {
			path := filepath.Join(filepath.Dir(i.Origin), pathStr)
			err := i.Extend(path)
			if err != nil {
				panic(i.error(err.Error()))
			}
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

func (i *Interpreter) Extend(file string) error {
	fileContents, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	lexer := lex.New()
	tokens, err := lexer.ReadInput(string(fileContents))
	if err != nil {
		return err
	}

	parser := parser.New()
	statements, err := parser.Parse(tokens)
	if err != nil {
		return err
	}

	_, err = i.Evaluate(statements)
	if err != nil {
		return err
	}

	return nil
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
				stdout := runShellCommandAndCapture(run.Body.Text, nil, i.Config.getShell())
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

func (i *Interpreter) error(message string) *RuntimeError {
	err := &RuntimeError{
		Message: fmt.Sprintf("runtime error: %s\n", message),
	}
	return err
}

type RuntimeError struct {
	Message string
}

func (re *RuntimeError) Error() string {
	return re.Message
}

func runShellCommandAndCapture(cmdString string, variables map[string]interface{}, shell string) []byte {
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

func runShellCommandAndPipeToStdout(cmdString string, variables map[string]interface{}, shell string) error {
	if len(cmdString) == 0 {
		// todo
	}
	cmd := exec.Command(shell, "-c", cmdString)
	cmd.Env = os.Environ()
	for name, value := range variables {
		cmd.Env = append(
			cmd.Env,
			fmt.Sprintf("%s=%s", name, trimQuotes(value)),
		)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func trimQuotes(input interface{}) interface{} {
	switch out := input.(type) {
	case string:
		return strings.Trim(out, "\"")
	}
	return input
}
