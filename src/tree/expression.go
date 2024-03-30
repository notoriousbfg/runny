package tree

type Expression interface {
	Accept(visitor ExpressionVisitor) interface{}
}

type ExpressionVisitor interface {
	// VisitVariableExpr(expr VariableExpression) interface{}
	VisitLiteralExpr(expr Literal) interface{}
}

type Literal struct {
	Value interface{}
}

func (l Literal) Accept(visitor ExpressionVisitor) interface{} {
	return visitor.VisitLiteralExpr(l)
}

// type VariableExpression struct {
// 	Name token.Token
// }

// func (v VariableExpression) Accept(visitor ExpressionVisitor) interface{} {
// 	return visitor.VisitVariableExpr(v)
// }
