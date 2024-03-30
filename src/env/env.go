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

// func (e *Environment) Assign(name string, value interface{}) {
// 	if _, ok := e.Values[name]; ok {
// 		e.Define(name, value)
// 	}
// 	if e.Enclosing != nil {
// 		e.Enclosing.Assign(name, value)
// 	}
// }

// func (e *Environment) AssignAt(distance int, name string, value interface{}) {
// 	e.ancestor(distance).Define(name, value)
// }

func (e *Environment) Get(name string) (interface{}, error) {
	if val, ok := e.Values[name]; ok {
		return val, nil
	}

	if e.Enclosing != nil {
		return e.Enclosing.Get(name)
	}

	return nil, fmt.Errorf("undefined variable '" + name + "'.")
}

func (e *Environment) GetAt(distance int, name string) interface{} {
	return e.ancestor(distance).Values[name]
}

func (e *Environment) printValues() {
	for _, val := range e.Values {
		fmt.Println(val)
	}
}

func (e *Environment) ancestor(distance int) *Environment {
	environment := e
	for i := 0; i < distance; i++ {
		environment = environment.Enclosing
	}
	return environment
}
