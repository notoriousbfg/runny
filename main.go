package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runny/src/interpreter"
	"runny/src/lexer"
	"runny/src/parser"
)

type Runny struct {
	Lexer       *lexer.Lexer
	Parser      *parser.Parser
	Interpreter *interpreter.Interpreter
	Config      Config
}

func (r *Runny) Scan() error {
	var err error
	fileContents, err := os.ReadFile(r.Config.File)
	if err != nil {
		return err
	}
	r.Lexer, err = lexer.New(string(fileContents))
	if err != nil {
		return err
	}
	return nil
}

func (r *Runny) Parse() error {
	r.Parser = parser.New(r.Lexer.Tokens)
	err := r.Parser.Parse()
	if err != nil {
		return err
	}
	return nil
}

func (r *Runny) Evaluate() {
	r.Interpreter = interpreter.New(r.Parser.Statements)
	r.Interpreter.Evaluate()
}

type Config struct {
	Debug bool
	File  string
}

func main() {
	var fileFlag string
	flag.StringVar(&fileFlag, "f", "config.rny", "config file location")
	flag.Parse()

	runny := Runny{
		Config: Config{
			Debug: true,
		},
	}

	file, err := configFile(fileFlag)
	if err != nil {
		fmt.Println("config error:", err)
		return
	}
	runny.Config.File = file

	if err := runny.Scan(); err != nil {
		if runny.Config.Debug {
			fmt.Print("scan error: ", err, ", (tokens:", lexer.TokenNames(runny.Lexer.Tokens), ")")
		} else {
			fmt.Print("scan error: ", err)
		}
		return
	}

	// i think we can condense the scan & parse stages into one by using a channel
	if err := runny.Parse(); err != nil {
		fmt.Print("parse error: ", err)
		return
	}

	runny.Evaluate()
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
