package tree

import (
	"runny/src/token"
)

type Statement interface {
	Accept(visitor StatementVisitor) interface{}
}

type StatementVisitor interface {
	VisitConfigStatement(statement ConfigStatement) interface{}
	VisitVariableStatement(statement VariableStatement) interface{}
	VisitTargetStatement(statement TargetStatement) interface{}
	VisitActionStatement(statement ActionStatement) interface{}
	VisitRunStatement(statement RunStatement) interface{}
	VisitDescribeStatement(statement DescribeStatement) interface{}
	VisitExtendsStatement(statement ExtendsStatement) interface{}
	VisitExpressionStatement(statement ExpressionStatement) interface{}
}

type ConfigStatement struct {
	Items []Config
}

type Config struct {
	Name        token.Token
	Initialiser Statement
}

func (c ConfigStatement) Accept(visitor StatementVisitor) interface{} {
	return visitor.VisitConfigStatement(c)
}

type VariableStatement struct {
	Items  []Variable
	Parent Statement
	Order  ExecutionOrder
}

type Variable struct {
	Name        token.Token
	Initialiser Statement
}

func (vs VariableStatement) Accept(visitor StatementVisitor) interface{} {
	return visitor.VisitVariableStatement(vs)
}

type TargetStatement struct {
	Name   token.Token
	Body   []Statement
	Parent Statement
}

func (ts TargetStatement) Accept(visitor StatementVisitor) interface{} {
	return visitor.VisitTargetStatement(ts)
}

type ActionStatement struct {
	Body   token.Token
	Parent Statement
}

func (as ActionStatement) Accept(visitor StatementVisitor) interface{} {
	return visitor.VisitActionStatement(as)
}

type ExecutionOrder int

const (
	DURING ExecutionOrder = iota
	BEFORE
	AFTER
)

type RunStatement struct {
	Name   token.Token
	Body   []Statement
	Order  ExecutionOrder
	Parent Statement
}

func (rs RunStatement) Accept(visitor StatementVisitor) interface{} {
	return visitor.VisitRunStatement(rs)
}

type DescribeStatement struct {
	Name   token.Token
	Lines  []Literal
	Parent Statement
}

func (ds DescribeStatement) Accept(visitor StatementVisitor) interface{} {
	return visitor.VisitDescribeStatement(ds)
}

type ExtendsStatement struct {
	Paths []Expression
}

func (es ExtendsStatement) Accept(visitor StatementVisitor) interface{} {
	return visitor.VisitExtendsStatement(es)
}

type ExpressionStatement struct {
	Expression Expression
}

func (es ExpressionStatement) Accept(visitor StatementVisitor) interface{} {
	return visitor.VisitExpressionStatement(es)
}
