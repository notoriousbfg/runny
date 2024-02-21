package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Runny struct {
	Config Config
	Body   RunnyFileBody
}

func (r *Runny) Load() error {
	fileContents, err := os.ReadFile(r.Config.File)
	if err != nil {
		return err
	}

	structure := RunnyFileBody{}
	err = yaml.Unmarshal(fileContents, &structure)
	if err != nil {
		return err
	}

	r.Body = structure
	return nil
}

func (r *Runny) Run() error {
	// var output [][]byte
	for targetName, runVars := range r.Body.Run {
		if targetCmd, exists := r.Body.Targets[targetName]; exists {
			// switch typedCmd := cmd.(type) {
			// case map[string]string:
			// 	for key, val := range typedCmd {
			// 	}
			// }
			// exec command here
			cmd := exec.Command("sh", "-c", targetCmd)

			if vars, ok := runVars.(map[interface{}]interface{}); ok {
				cmd.Env = os.Environ()
				for name, val := range vars {
					nameStr, isStr := name.(string)
					if !isStr {
						return fmt.Errorf("the variable '%s' could not be parsed", name)
					}
					valStr, isStr := val.(string)
					if !isStr {
						return fmt.Errorf("the variable '%s' could not be parsed", name)
					}
					cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", nameStr, valStr))
				}
			}
			fmt.Printf("running command: %s", targetCmd)
			stdOutStdErr, err := cmd.CombinedOutput()
			if err != nil {
				return err
			}
			fmt.Println(string(stdOutStdErr))

		} else {
			return fmt.Errorf("target '%s' not found", targetName)
		}
	}
	return nil
}

type Config struct {
	File string
}

type RunnyFileBody struct {
	Vars    map[string]string      `yaml:"vars"`
	Targets map[string]string      `yaml:"targets"`
	Run     map[string]interface{} `yaml:"run"`
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

	err = runny.Load()
	if err != nil {
		fmt.Println("parse error:", err)
		return
	}

	err = runny.Run()
	if err != nil {
		fmt.Println("runtime error:", err)
		return
	}
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
