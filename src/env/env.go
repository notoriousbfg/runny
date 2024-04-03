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
		Values:    NewValues(),
		Enclosing: enclosing,
		Depth:     depth,
	}
}

type ValueType int

func (vt ValueType) String() string {
	switch vt {
	case VTVar:
		return "variable"
	case VTTarget:
		return "target"
	}
	return "unknown"
}

const (
	VTUnknown ValueType = iota
	VTVar
	VTTarget
)

type Values struct {
	Vars    map[string]interface{}
	Targets map[string]interface{}
}

func NewValues() Values {
	return Values{
		Vars:    make(map[string]interface{}),
		Targets: map[string]interface{}{},
	}
}

type Environment struct {
	Values    Values
	Enclosing *Environment
	Depth     int
}

func (e *Environment) Define(name string, valueType ValueType, value interface{}) {
	if len(name) > 0 {
		switch valueType {
		case VTVar:
			e.Values.Vars[name] = value
		case VTTarget:
			e.Values.Targets[name] = value
		}
	}
}

func (e *Environment) Get(name string, valueType ValueType) (interface{}, error) {
	switch valueType {
	case VTVar:
		if val, ok := e.Values.Vars[name]; ok {
			return val, nil
		}

		if e.Enclosing != nil {
			return e.Enclosing.Get(name, VTVar)
		}
	case VTTarget:
		if val, ok := e.Values.Targets[name]; ok {
			return val, nil
		}

		if e.Enclosing != nil {
			return e.Enclosing.Get(name, VTTarget)
		}
	}
	return nil, fmt.Errorf("undefined %s '%s'", valueType, name)
}

func (e *Environment) GetAll(valueType ValueType) map[string]interface{} {
	switch valueType {
	case VTVar:
		vars := e.Values.Vars
		if e.Enclosing != nil {
			for k, v := range e.Enclosing.GetAll(VTVar) {
				// local variable is preferred to global
				if _, exists := vars[k]; !exists {
					vars[k] = v
				}
			}
		}
		return vars
	case VTTarget:
		targets := e.Values.Targets
		if e.Enclosing != nil {
			for k, v := range e.Enclosing.GetAll(VTTarget) {
				// local variable is preferred to global
				if _, exists := targets[k]; !exists {
					targets[k] = v
				}
			}
		}
		return targets
	}
	return make(map[string]interface{}, 0)
}
