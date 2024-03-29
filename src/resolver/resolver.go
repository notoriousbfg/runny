package resolver

import (
	"runny/src/interpreter"
	"runny/src/token"
	"runny/src/tree"
)

func NewScopeStack() *ScopeStack {
	return &ScopeStack{}
}

type Scope map[string]bool

func (s Scope) put(key string, val bool) {
	s[key] = val
}

type ScopeStack struct {
	scopes []Scope
}

func (s *ScopeStack) push(scope Scope) {
	s.scopes = append(s.scopes, scope)
}

func (s *ScopeStack) pop() Scope {
	if len(s.scopes) == 0 {
		return nil
	}
	lastIndex := len(s.scopes) - 1
	topScope := s.scopes[lastIndex]
	s.scopes = s.scopes[:lastIndex]
	return topScope
}

func (s *ScopeStack) isEmpty() bool {
	return len(s.scopes) == 0
}

func (s *ScopeStack) size() int {
	return len(s.scopes)
}

func (s *ScopeStack) peek() Scope {
	if len(s.scopes) == 0 {
		return nil
	}
	return s.scopes[len(s.scopes)-1]
}

func NewResolver(interpreter interpreter.Interpreter) *Resolver {
	return &Resolver{
		Interpreter: interpreter,
		Scopes:      NewScopeStack(),
	}
}

type Resolver struct {
	Interpreter interpreter.Interpreter
	Scopes      *ScopeStack
}

func (r *Resolver) VisitTargetStatement(statement tree.TargetStatement) interface{} {
	return nil
}

func (r *Resolver) VisitVariableStatement(statement tree.VariableStatement) interface{} {
	return nil
}

func (r *Resolver) VisitRunStatement(statement tree.RunStatement) interface{} {
	return nil
}

func (r *Resolver) VisitActionStatement(statement tree.ActionStatement) interface{} {
	return nil
}

func (r *Resolver) VisitExpressionStatement(statement tree.ExpressionStatement) interface{} {
	return nil
}

func (r *Resolver) VisitLiteralExpr(expr tree.Literal) interface{} {
	return expr.Value
}

func (r *Resolver) resolveStatements(statements []tree.Statement) {
	for _, statement := range statements {
		r.resolveStatement(statement)
	}
}

func (r *Resolver) resolveExpression(expression tree.Expression) {
	expression.Accept(r)
}

func (r *Resolver) resolveStatement(statement tree.Statement) {
	statement.Accept(r)
}

func (r *Resolver) beginScope() {
	r.Scopes.push(Scope{})
}

func (r *Resolver) endScope() {
	r.Scopes.pop()
}

func (r *Resolver) declare(name token.Token) {
	if r.Scopes.isEmpty() {
		return
	}
	r.Scopes.peek().put(name.Text, false)
}

func (r *Resolver) define(name token.Token) {
	if r.Scopes.isEmpty() {
		return
	}
	scope := r.Scopes.peek()
	scope.put(name.Text, false)
}
