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

// set a key/value within a scope
func (s Scope) put(key string, val bool) {
	s[key] = val
}

// has returns whether a variable with the given
// name is declared and defined in this scope
func (s Scope) has(name string) (declared, defined bool) {
	v, ok := s[name]
	if !ok {
		return false, false
	}
	return true, v
}

type ScopeStack struct {
	scopes []Scope
}

// append to stack
func (s *ScopeStack) push(scope Scope) {
	s.scopes = append(s.scopes, scope)
}

// remove scope from top of stack
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

// get last scope in stack
func (s *ScopeStack) peek() Scope {
	if len(s.scopes) == 0 {
		return nil
	}
	return s.scopes[len(s.scopes)-1]
}

// get last scope in stack
func (s *ScopeStack) get(index int) Scope {
	if len(s.scopes) == 0 {
		return nil
	}
	return s.scopes[index]
}

func NewResolver(interpreter *interpreter.Interpreter) *Resolver {
	return &Resolver{
		Interpreter: interpreter,
		Scopes:      NewScopeStack(),
	}
}

type Resolver struct {
	Interpreter *interpreter.Interpreter
	Scopes      *ScopeStack
}

func (r *Resolver) ResolveStatements(statements []tree.Statement) {
	for _, statement := range statements {
		r.resolveStatement(statement)
	}
}

func (r *Resolver) VisitTargetStatement(statement tree.TargetStatement) interface{} {
	r.declare(statement.Name)
	r.define(statement.Name)

	r.beginScope()
	r.ResolveStatements(statement.Body)
	r.endScope()
	return nil
}

func (r *Resolver) VisitVariableStatement(statement tree.VariableStatement) interface{} {
	for _, item := range statement.Items {
		r.declare(item.Name)
		r.resolveStatement(item.Initialiser)
		r.define(item.Name)
	}
	return nil
}

func (r *Resolver) VisitRunStatement(statement tree.RunStatement) interface{} {
	r.beginScope()
	r.ResolveVariables() // ?
	r.ResolveStatements(statement.Body)
	r.endScope()
	return nil
}

func (r *Resolver) VisitActionStatement(statement tree.ActionStatement) interface{} {
	return nil
}

func (r *Resolver) VisitExpressionStatement(statement tree.ExpressionStatement) interface{} {
	// r.resolveExpression(statement.Expression)
	return nil
}

func (r *Resolver) VisitLiteralExpr(expr tree.Literal) interface{} {
	return nil
}

func (r *Resolver) ResolveVariables() map[string]interface{} {
	for i := r.Scopes.size() - 1; i >= 0; i-- {
		s := r.Scopes.get(i)
		if _, defined := s.has(); defined {
			// 	depth := len(r.scopes) - 1 - i
			// 	r.interpreter.Resolve(expr, depth)
			// 	s.use(name.Lexeme)
			// 	return
		}
	}

	return map[string]interface{}{}
}

// func (r *Resolver) resolveExpression(expression tree.Expression) {
// 	expression.Accept(r)
// }

func (r *Resolver) resolveStatement(statement tree.Statement) {
	statement.Accept(r)
}

// func (r *Resolver) resolveLocal() {
// 	statement.Accept(r)
// }

func (r *Resolver) beginScope() {
	r.Scopes.push(make(Scope))
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
	r.Scopes.peek().put(name.Text, true)
}
