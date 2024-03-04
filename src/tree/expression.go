package tree

import "runny/src/token"

type Expression interface {
	Accept(visitor ExpressionVisitor) interface{}
}

type ExpressionVisitor interface {
	VisitLiteralExpression(expression Literal) interface{}
	VisitVariableExpression(expression VariableExpression) interface{}
}

type VariableExpression struct {
	Name token.Token
}

func (v VariableExpression) Accept(visitor ExpressionVisitor) interface{} {
	return visitor.VisitVariableExpression(v)
}

type Literal struct {
	Value interface{}
}

func (l Literal) Accept(visitor ExpressionVisitor) interface{} {
	return visitor.VisitLiteralExpression(l)
}
