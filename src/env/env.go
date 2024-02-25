package env

import (
	"fmt"
)

func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		Values: Values{
			Variables: make(map[string]interface{}, 0),
			Targets:   make(map[string][]interface{}, 0),
			Runs:      make(map[string][]interface{}, 0),
		},
		Enclosing: enclosing,
	}
}

type ValueType int

const (
	VariableType ValueType = iota
	TargetType
	RunType
)

type Values struct {
	Variables map[string]interface{}
	Targets   map[string][]interface{}
	Runs      map[string][]interface{}
}

type Environment struct {
	Values    Values
	Enclosing *Environment
}

func (e *Environment) Define(name string, value interface{}, valueType ValueType) {
	switch valueType {
	case VariableType:
		e.Values.Variables[name] = value
		if e.Enclosing != nil {
			if _, ok := e.Enclosing.Values.Variables[name]; ok {
				e.Enclosing.Values.Variables[name] = value
			}
		}
	}
}

func (e *Environment) Get(name string, valueType ValueType) (interface{}, error) {
	switch valueType {
	case VariableType:
		if _, ok := e.Values.Variables[name]; ok {
			return e.Values.Variables[name], nil
		}

		if e.Enclosing != nil {
			if _, ok := e.Enclosing.Values.Variables[name]; ok {
				return e.Enclosing.Values.Variables[name], nil
			}
		}
	}

	return nil, fmt.Errorf("undefined variable '" + name + "'.")
}

func (e *Environment) GetAll(valueType ValueType) map[string]interface{} {
	switch valueType {
	case VariableType:
		return e.Values.Variables
	}
	return nil
}
