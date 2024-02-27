package tree

import (
	"runny/src/token"
)

type Statement interface {
	Accept(visitor StatementVisitor) interface{}
}

type StatementVisitor interface {
	VisitVariableStatement(stmt VariableStatement) interface{}
	VisitTargetStatement(stmt TargetStatement) interface{}
	VisitActionStatement(stmt ActionStatement) interface{}
	VisitRunStatement(stmt RunStatement) interface{}
	VisitExpressionStatement(stmt ExpressionStatement) interface{}
}

type VariableStatement struct {
	Items []Variable
}

type Variable struct {
	Name        token.Token
	Initialiser Statement
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
	Body token.Token
}

func (as ActionStatement) Accept(visitor StatementVisitor) interface{} {
	return visitor.VisitActionStatement(as)
}

// func (as ActionStatement) String() string {
// 	// var builder strings.Builder
// 	// for _, t := range as.Body {
// 	// 	if t.Type == token.NEWLINE {
// 	// 		builder.WriteString("\n")
// 	// 	} else {
// 	// 		builder.WriteString(fmt.Sprintf(" %s", t.Text))
// 	// 	}
// 	// }
// 	// return builder.String()
// 	return as.Body.Text
// }

type RunStatement struct {
	Name token.Token
	Body []Statement
}

func (rs RunStatement) Accept(visitor StatementVisitor) interface{} {
	return visitor.VisitRunStatement(rs)
}

type ExpressionStatement struct {
	Expression Expression
}

func (es ExpressionStatement) Accept(visitor StatementVisitor) interface{} {
	return visitor.VisitExpressionStatement(es)
}
