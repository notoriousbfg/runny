package env

import (
	"fmt"
	"runny/src/tree"
)

func NewEnvironment(enclosing *Environment) *Environment {
	depth := 0
	if enclosing != nil {
		depth = enclosing.Depth + 1
	}
	return &Environment{
		Values: Values{
			Variables: make(map[string]tree.Statement, 0),
			Targets:   make(map[string][]tree.Statement, 0),
			Runs:      make(map[string][]tree.Statement, 0),
		},
		Enclosing: enclosing,
		Depth:     depth,
	}
}

type Values struct {
	Variables map[string]tree.Statement
	Targets   map[string][]tree.Statement
	Runs      map[string][]tree.Statement
}

type Environment struct {
	Values    Values
	Enclosing *Environment
	Depth     int
}

func (e *Environment) DefineVariable(name string, value tree.Statement) {
	if len(name) > 0 {
		e.Values.Variables[name] = value
	}
}

func (e *Environment) DefineTarget(name string, value []tree.Statement) {
	if len(name) > 0 {
		e.Values.Targets[name] = value
	}
}

func (e *Environment) GetVariable(name string) (tree.Statement, error) {
	if val, ok := e.Values.Variables[name]; ok {
		return val, nil
	}

	if e.Enclosing != nil {
		return e.Enclosing.GetVariable(name)
	}

	return nil, fmt.Errorf("undefined variable '" + name + "'.")
}

func (e *Environment) GetTarget(name string) ([]tree.Statement, error) {
	if val, ok := e.Values.Targets[name]; ok {
		return val, nil
	}

	if e.Enclosing != nil {
		return e.Enclosing.GetTarget(name)
	}

	return nil, fmt.Errorf("undefined target '" + name + "'.")
}

type ValueType int

const (
	VariableType ValueType = iota
	TargetType
	RunType
)

func (e *Environment) GetAll(valueType ValueType) map[string]tree.Statement {
	switch valueType {
	case VariableType:
		vars := make(map[string]tree.Statement, 0)
		for key := range e.Values.Variables {
			val, _ := e.GetVariable(key)
			vars[key] = val
		}
		return vars
	}
	return nil
}
