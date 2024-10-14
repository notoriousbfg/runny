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
	"sort"
	"strings"
)

func New(origin string, printOutput bool) *Interpreter {
	return &Interpreter{
		Config:      make(map[string]interface{}, 0),
		Environment: env.NewEnvironment(nil),
		Origin:      origin,
		Printer:     &Printer{},
		PrintOutput: printOutput,
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
	PrintOutput bool
	Printer     *Printer
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
	if i.PrintOutput {
		i.Printer.Print()
	}
	return
}

func (i *Interpreter) FilterStatementsByTarget(targetStr string, statements []tree.Statement) ([]tree.Statement, error) {
	var foundTarget *tree.TargetStatement
	filteredStatements := make([]tree.Statement, 0)
	for _, statement := range statements {
		if _, isRun := statement.(tree.RunStatement); isRun {
			continue
		}
		if target, isTarget := statement.(tree.TargetStatement); isTarget {
			if target.Name.Text == targetStr {
				foundTarget = &target
			}
		}
		filteredStatements = append(filteredStatements, statement)
	}
	if foundTarget == nil {
		return nil, fmt.Errorf("target '%s' does not exist", targetStr)
	}
	filteredStatements = append(filteredStatements, tree.RunStatement{
		Name: foundTarget.Name,
	})
	return filteredStatements, nil
}

func (i *Interpreter) Accept(statement tree.Statement) interface{} {
	return statement.Accept(i)
}

func (i *Interpreter) VisitConfigStatement(statement tree.ConfigStatement) interface{} {
	for _, config := range statement.Items {
		i.Config[config.Name.Text] = i.Accept(config.Initialiser)
	}
	return nil
}

func (i *Interpreter) VisitVariableStatement(statement tree.VariableStatement) interface{} {
	for _, variable := range statement.Items {
		i.Environment.Define(variable.Name.Text, env.VTVar, variable.Initialiser)
	}
	return nil
}

func (i *Interpreter) VisitTargetStatement(statement tree.TargetStatement) interface{} {
	// presort body by order
	statements := statement.Body
	sort.SliceStable(statements, func(i, j int) bool {
		return orderValue(statements[i]) < orderValue(statements[j])
	})
	i.Environment.Define(statement.Name.Text, env.VTTarget, statements)
	return nil
}

func orderValue(statement tree.Statement) int {
	switch statementTyped := statement.(type) {
	case tree.RunStatement:
		switch statementTyped.Stage {
		case tree.BEFORE:
			return 1
		case tree.DURING:
			return 2
		case tree.AFTER:
			return 3
		}
	}
	return 2
}

const (
	foreColour = "\033[32m"
	aftColour  = "\033[0m"
)

func (i *Interpreter) VisitActionStatement(statement tree.ActionStatement) interface{} {
	evaluated := make(map[string]interface{}, 0)
	for k := range i.Environment.GetAll(env.VTVar) {
		variable, _ := i.lookupVariable(k)
		evaluated[k] = variable
	}

	// print command string (highlighted)
	i.Printer.PushStr(fmt.Sprintf("%s%s%s\n", foreColour, relativeDedent(statement.Body.Text), aftColour))

	cmd := createCommand(statement.Body.Text, evaluated, i.Config.getShell())

	// creates a pipe to stdout that can be scanned by printer instance
	cmdOut, err := cmd.StdoutPipe()
	if err != nil {
		panic(i.error(fmt.Sprintf("could not create command pipe: %s", err.Error())))
	}

	cmdErr, err := cmd.StderrPipe()
	if err != nil {
		panic(i.error(fmt.Sprintf("could not create command pipe: %s", err.Error())))
	}

	i.Printer.Push(Statement{
		Cmd:    cmd, // cmd included here so printer can wait
		StdOut: cmdOut,
		StdErr: cmdErr,
	})

	if err := cmd.Start(); err != nil {
		panic(i.error(fmt.Sprintf("could not run command: %s", err.Error())))
	}

	return nil
}

func relativeDedent(inputString string) string {
	lines := strings.Split(inputString, "\n")
	if len(lines) > 1 {
		lowestPositiveIndent := 0
		for _, line := range lines {
			whitespace := countLeadingSpaces(line)
			if lowestPositiveIndent <= 0 {
				lowestPositiveIndent = whitespace
				continue
			}
			if whitespace < lowestPositiveIndent && whitespace != 0 {
				lowestPositiveIndent = whitespace
			}
		}
		for i, line := range lines {
			if i == 0 {
				continue
			}
			lines[i] = line[lowestPositiveIndent:]
		}
		return strings.Join(lines, "\n")
	}
	return inputString
}

func countLeadingSpaces(line string) int {
	return len(line) - len(strings.TrimLeft(line, " "))
}

func (i *Interpreter) VisitRunStatement(statement tree.RunStatement) interface{} {
	startEnvironment := i.Environment
	i.Environment = env.NewEnvironment(i.Environment)
	defer func() {
		i.Environment = startEnvironment
	}()

	body := statement.Body

	if statement.Name != (token.Token{}) {
		targetBodyInt, err := i.Environment.Get(statement.Name.Text, env.VTTarget)
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

func (i *Interpreter) VisitDescribeStatement(statement tree.DescribeStatement) interface{} {
	for _, line := range statement.Lines {
		i.Printer.PushStr(fmt.Sprintf("> %s\n", trimQuotes(line.Value)))
	}
	return nil
}

func (i *Interpreter) VisitExtendsStatement(statement tree.ExtendsStatement) interface{} {
	for _, path := range statement.Paths {
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

func (i *Interpreter) VisitExpressionStatement(statement tree.ExpressionStatement) interface{} {
	return statement.Expression.Accept(i)
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
				cmd := createCommand(run.Body.Text, nil, i.Config.getShell())
				stdOutStdErr, _ := cmd.CombinedOutput()
				// opinionated: always trim trailing newline
				// var runs are mostly variables inserted into something else
				// where it's not helpful to have a trailing newline
				trimmedOutput := strings.TrimRight(string(stdOutStdErr), "\n")
				strBuilder.WriteString(trimmedOutput)
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

func createCommand(cmdString string, variables map[string]interface{}, shell string) *exec.Cmd {
	cmd := exec.Command(shell, "-c", cmdString)
	cmd.Env = os.Environ()
	for name, value := range variables {
		cmd.Env = append(
			cmd.Env,
			fmt.Sprintf("%s=%s", name, trimQuotes(value)),
		)
	}
	return cmd
}

func trimQuotes(input interface{}) interface{} {
	switch out := input.(type) {
	case string:
		return strings.Trim(out, "\"")
	}
	return input
}
