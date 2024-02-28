package main

import (
	"os"
	"path/filepath"
	"runny/src/lexer"
)

func main() {
	path := os.Args[1]

	path, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}

	_, err = os.Stat(path)
	if err != nil {
		panic(err)
	}

	fileContents, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	lexer.TokenGenerator(string(fileContents))
}
