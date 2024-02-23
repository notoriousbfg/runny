package tree

import "runny/src/token"

type Statement interface {
	Accept(visitor StatementVisitor) interface{}
}

type StatementVisitor interface {
	VisitVariableStatement(stmt VariableStatement) interface{}
	VisitExpressionStatement(stmt ExpressionStatement) interface{}
}

type VariableStatement struct {
	Items []Variable
}

type Variable struct {
	Name        token.Token
	Initialiser Expression
}

func (vs VariableStatement) Accept(visitor StatementVisitor) interface{} {
	return visitor.VisitVariableStatement(vs)
}

type ExpressionStatement struct {
	Expression Expression
}

func (es ExpressionStatement) Accept(visitor StatementVisitor) interface{} {
	return visitor.VisitExpressionStatement(es)
}
