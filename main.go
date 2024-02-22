package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

type Runny struct {
	Config Config
}

type Config struct {
	File string
}

func main() {
	var fileFlag string
	flag.StringVar(&fileFlag, "f", "runny.yml", "config file location")
	flag.Parse()

	runny := Runny{
		Config: Config{},
	}

	file, err := configFile(fileFlag)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	runny.Config.File = file

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
	if extension != ".yml" && extension != ".yaml" {
		return "", fmt.Errorf("the config file is not a yaml file")
	}
	return path, nil
}
