package env

import (
	"fmt"
	"runny/src/token"
)

func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		Values:    make(map[string]interface{}),
		Enclosing: enclosing,
	}
}

type Environment struct {
	Values    map[string]interface{}
	Enclosing *Environment
}

func (e *Environment) Define(name string, value interface{}) {
	e.Values[name] = value

	// only modify global if exists
	if e.Enclosing != nil {
		if _, ok := e.Enclosing.Values[name]; ok {
			e.Enclosing.Values[name] = value
		}
	}
}

func (e *Environment) Get(token token.Token) (interface{}, error) {
	if _, ok := e.Values[token.Text]; ok {
		return e.Values[token.Text], nil
	}

	if e.Enclosing != nil {
		if _, ok := e.Enclosing.Values[token.Text]; ok {
			return e.Enclosing.Values[token.Text], nil
		}
	}

	return nil, fmt.Errorf("undefined variable '" + token.Text + "'.")
}
