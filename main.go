package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runny/src/lexer"
)

type Runny struct {
	Lexer  lexer.Lexer
	Config Config
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
	return nil
}

type Config struct {
	File string
}

func main() {
	var fileFlag string
	flag.StringVar(&fileFlag, "f", "config.rny", "config file location")
	flag.Parse()

	runny := Runny{
		Config: Config{},
	}

	file, err := configFile(fileFlag)
	if err != nil {
		fmt.Println("config error:", err)
		return
	}
	runny.Config.File = file

	if err := runny.Scan(); err != nil {
		fmt.Println("scan error:", err, ", (tokens:", lexer.TokenTypeNames(runny.Lexer.TokenTypes()), ")")
		return
	}

	fmt.Print(runny.Lexer.Tokens)
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
