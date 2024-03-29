package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runny/src/interpreter"
	"runny/src/lexer"
	"runny/src/parser"
	"runny/src/resolver"
	"runny/src/token"
	"runny/src/tree"
	"strings"
)

type Runny struct {
	Lexer       *lexer.Lexer
	Parser      *parser.Parser
	Interpreter *interpreter.Interpreter
	Config      Config
}

func (r *Runny) Scan() ([]token.Token, error) {
	var err error
	fileContents, err := os.ReadFile(r.Config.File)
	if err != nil {
		return []token.Token{}, err
	}
	r.Lexer = lexer.New(string(fileContents))
	tokens, err := r.Lexer.ReadInput()
	if err != nil {
		return []token.Token{}, err
	}
	return tokens, nil
}

func (r *Runny) Parse(tokens []token.Token) ([]tree.Statement, error) {
	r.Parser = parser.New(tokens)
	statements, err := r.Parser.Parse()
	if err != nil {
		return []tree.Statement{}, nil
	}
	return statements, nil
}

func (r *Runny) Resolve(statements []tree.Statement) {
	resolver := resolver.NewResolver(r.Interpreter)
	resolver.ResolveStatements(statements)
}

func (r *Runny) Evaluate() {
	statements := r.Interpreter.Statements
	if r.Config.Target != "" {
		var foundTarget *tree.TargetStatement
		filteredStatements := make([]tree.Statement, 0)
		for _, statement := range statements {
			if variable, isVariable := statement.(tree.VariableStatement); isVariable {
				filteredStatements = append(filteredStatements, variable)
			}
			if target, isTarget := statement.(tree.TargetStatement); isTarget {
				if target.Name.Text == r.Config.Target {
					foundTarget = &target
					filteredStatements = append(filteredStatements, target)
				}
			}
		}
		if foundTarget == nil {
			fmt.Printf("target '%s' does not exist", r.Config.Target)
			return
		}
		filteredStatements = append(filteredStatements, tree.RunStatement{
			Name: foundTarget.Name,
		})
		statements = filteredStatements
	}
	r.Interpreter.Evaluate(statements)
}

type Config struct {
	Target string
	File   string
	Debug  bool
}

func main() {
	target, fileFlag := parseArgs()

	runny := Runny{
		Config: Config{
			Target: target,
			Debug:  true,
		},
	}

	file, err := configFile(fileFlag)
	if err != nil {
		fmt.Println("config error:", err)
		return
	}
	runny.Config.File = file

	tokens, err := runny.Scan()
	if err != nil {
		if runny.Config.Debug {
			fmt.Print(err, ", (tokens:", lexer.TokenNames(runny.Lexer.Tokens), ")")
		} else {
			fmt.Print(err)
		}
		return
	}

	// i think we can condense the scan & parse stages into one by using a channel
	statements, err := runny.Parse(tokens)
	if err != nil {
		fmt.Print(err)
		return
	}

	runny.Interpreter = interpreter.New(statements)
	runny.Resolve(statements)
	// if err != nil {
	// 	fmt.Print(err)
	// 	return
	// }

	runny.Evaluate()
}

func parseArgs() (string, string) {
	var target string
	var foundFlag bool
	var fileFlag string
	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "-f") {
			foundFlag = true
		} else if !foundFlag && target == "" {
			target = arg
		} else if foundFlag && fileFlag == "" {
			fileFlag = arg
			foundFlag = false
		}
	}

	if fileFlag == "" {
		fileFlag = "runny.rny"
	}
	return target, fileFlag
}

func configFile(flag string) (string, error) {
	path, err := filepath.Abs(flag)
	if err != nil {
		return "", err
	}
	_, err = os.Stat(path)
	if err != nil {
		return "", err
	}
	extension := filepath.Ext(path)
	if extension != ".rny" {
		return "", fmt.Errorf("rny file not found")
	}
	return path, nil
}
