package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runny/src/interpreter"
	"runny/src/lex"
	"runny/src/parser"
	"runny/src/tree"
	"strings"
)

type Runny struct {
	Config Config
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
			Debug:  os.Getenv("DEBUG") == "true",
		},
	}

	file, err := configFile(fileFlag)
	if err != nil {
		fmt.Println("config error:", err)
		return
	}
	runny.Config.File = file

	fileContents, err := os.ReadFile(file)
	if err != nil {
		fmt.Println("error reading file:", err)
	}
	lexer := lex.New()
	tokens, err := lexer.ReadInput(string(fileContents))
	if err != nil {
		if runny.Config.Debug {
			fmt.Print(err, ", (tokens:", lex.TokenNames(lexer.Tokens), ")")
		} else {
			fmt.Print(err)
		}
		return
	}

	// i think we can condense the scan & parse stages into one by using a channel
	parser := parser.New()
	statements, err := parser.Parse(tokens)
	if err != nil {
		fmt.Print(err)
		return
	}

	interpreter := interpreter.New(file)
	if runny.Config.Target != "" {
		var foundTarget *tree.TargetStatement
		filteredStatements := make([]tree.Statement, 0)
		for _, statement := range statements {
			if _, isRun := statement.(tree.RunStatement); isRun {
				continue
			}
			if target, isTarget := statement.(tree.TargetStatement); isTarget {
				if target.Name.Text == runny.Config.Target {
					foundTarget = &target
				}
			}
			filteredStatements = append(filteredStatements, statement)
		}
		if foundTarget == nil {
			fmt.Printf("target '%s' does not exist", runny.Config.Target)
			return
		}
		filteredStatements = append(filteredStatements, tree.RunStatement{
			Name: foundTarget.Name,
		})
		statements = filteredStatements
	}
	_, err = interpreter.Evaluate(statements)
	if err != nil {
		fmt.Print(err)
		return
	}
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
