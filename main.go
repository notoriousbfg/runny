package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runny/src/interpreter"
	"runny/src/lex"
	"runny/src/parser"
	"strings"
)

type Runny struct {
	Config Config
}

func (r *Runny) Run() {
	fileContents, err := os.ReadFile(r.Config.File)
	if err != nil {
		fmt.Println("error reading file:", err)
	}

	lexer := lex.New()
	tokens, err := lexer.ReadInput(string(fileContents))
	if err != nil {
		if r.Config.Debug {
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

	interpreter := interpreter.New(r.Config.File, !r.Config.Testing)
	if r.Config.Target != "" {
		var err error
		statements, err = interpreter.FilterStatementsByTarget(r.Config.Target, statements)
		if err != nil {
			fmt.Print(err)
			return
		}
	}

	_, err = interpreter.Evaluate(statements)
	if err != nil {
		fmt.Print(err)
		return
	}
}

type Config struct {
	Target  string
	File    string
	Debug   bool
	Testing bool
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

	runny.Run()
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
