package env

import (
	"fmt"
	"runny/src/tree"
)

func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		Values: Values{
			Variables: make(map[string]tree.Statement, 0),
			Targets:   make(map[string][]tree.Statement, 0),
			Runs:      make(map[string][]tree.Statement, 0),
		},
		Enclosing: enclosing,
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
}

func (e *Environment) DefineVariable(name string, value tree.Statement) {
	switch typedValue := value.(type) {
	case tree.ExpressionStatement:
		e.Values.Variables[name] = typedValue
		// if e.Enclosing != nil {
		// 	if _, ok := e.Enclosing.Values.Variables[name]; ok {
		// 		e.Enclosing.Values.Variables[name] = typedValue
		// 	}
		// }
	}
}

func (e *Environment) DefineTarget(name string, value []tree.Statement) {
	e.Values.Targets[name] = value
	if e.Enclosing != nil {
		if _, ok := e.Enclosing.Values.Targets[name]; ok {
			e.Enclosing.Values.Targets[name] = value
		}
	}
}

func (e *Environment) GetVariable(name string) (tree.Statement, error) {
	if _, ok := e.Values.Variables[name]; ok {
		return e.Values.Variables[name], nil
	}

	if e.Enclosing != nil {
		if _, ok := e.Enclosing.Values.Variables[name]; ok {
			return e.Enclosing.Values.Variables[name], nil
		}
	}

	return nil, fmt.Errorf("undefined variable '" + name + "'.")
}

func (e *Environment) GetTarget(name string) ([]tree.Statement, error) {
	if _, ok := e.Values.Targets[name]; ok {
		return e.Values.Targets[name], nil
	}

	if e.Enclosing != nil {
		fmt.Print(e.Enclosing.Values.Targets)
		if _, ok := e.Enclosing.Values.Targets[name]; ok {
			return e.Enclosing.Values.Targets[name], nil
		}
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
		if e.Enclosing != nil {
			for eKey, eVal := range e.Enclosing.Values.Variables {
				vars[eKey] = eVal
			}
		}
		for key, val := range e.Values.Variables {
			vars[key] = val
		}
		return vars
	}
	return nil
}
