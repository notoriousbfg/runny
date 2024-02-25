package interpreter

import (
	"runny/src/env"
	"runny/src/tree"
)

func New(statements []tree.Statement) *Interpreter {
	return &Interpreter{
		Statements: statements,
		Environment: &env.Environment{
			Enclosing: nil,
			Values:    make(map[string]interface{}),
		},
	}
}

type Interpreter struct {
	Statements  []tree.Statement
	Environment *env.Environment
}

func (i Interpreter) Evaluate() (result []interface{}) {
	// for _, statement := range i.Statements {
	// 	result = append(result, i.Execute(statement))
	// }
	return
}
