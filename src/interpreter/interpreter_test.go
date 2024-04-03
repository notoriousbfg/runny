package interpreter

import (
	"reflect"
	"runny/src/env"
	"runny/src/tree"
	"testing"
)

func TestInterpreter_Evaluate(t *testing.T) {
	type fields struct {
		Statements  []tree.Statement
		Environment *env.Environment
	}
	tests := []struct {
		name       string
		fields     fields
		wantResult []interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Interpreter{
				Environment: tt.fields.Environment,
			}
			if gotResult, _ := i.Evaluate(tt.fields.Statements); !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("Interpreter.Evaluate() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
