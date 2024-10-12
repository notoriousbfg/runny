package tree

type Expression interface {
	Accept(visitor ExpressionVisitor) interface{}
}

type ExpressionVisitor interface {
	VisitLiteralExpr(expr Literal) interface{}
}

type Literal struct {
	Value interface{}
}

func (l Literal) Accept(visitor ExpressionVisitor) interface{} {
	return visitor.VisitLiteralExpr(l)
}
