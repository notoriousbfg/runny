package env

import (
	"fmt"
)

func NewEnvironment(enclosing *Environment) *Environment {
	depth := 0
	if enclosing != nil {
		depth = enclosing.Depth + 1
	}
	return &Environment{
		Values:    make(map[string]interface{}),
		Enclosing: enclosing,
		Depth:     depth,
	}
}

type Environment struct {
	Values    map[string]interface{}
	Enclosing *Environment
	Depth     int
}

func (e *Environment) Define(name string, value interface{}) {
	if len(name) > 0 {
		e.Values[name] = value
	}
}

func (e *Environment) Get(name string) interface{} {
	if val, ok := e.Values[name]; ok {
		return val
	}

	if e.Enclosing != nil {
		return e.Enclosing.Get(name)
	}

	return ""
}

func (e *Environment) printValues() {
	for _, val := range e.Values {
		fmt.Println(val)
	}
}
