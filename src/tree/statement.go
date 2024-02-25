package tree

import "runny/src/token"

type Statement interface {
	Accept(visitor StatementVisitor) interface{}
}

type StatementVisitor interface {
	VisitVariableStatement(stmt VariableStatement) interface{}
	VisitTargetStatement(stmt TargetStatement) interface{}
	VisitActionStatement(stmt ActionStatement) interface{}
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

type TargetStatement struct {
	Name token.Token
	Body []Statement
}

func (ts TargetStatement) Accept(visitor StatementVisitor) interface{} {
	return visitor.VisitTargetStatement(ts)
}

type ActionStatement struct {
	Body []token.Token
}

func (as ActionStatement) Accept(visitor StatementVisitor) interface{} {
	return visitor.VisitActionStatement(as)
}

type ExpressionStatement struct {
	Expression Expression
}

func (es ExpressionStatement) Accept(visitor StatementVisitor) interface{} {
	return visitor.VisitExpressionStatement(es)
}
