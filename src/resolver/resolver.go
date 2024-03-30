package resolver

import (
	"runny/src/interpreter"
	"runny/src/token"
	"runny/src/tree"
)

func NewScopeStack() *ScopeStack {
	scopes := make([]Scope, 0)
	scopes = append(scopes, Scope{}) // start with a single scope

	return &ScopeStack{
		scopes: scopes,
	}
}

type ScopeVariable struct {
	token   token.Token
	defined bool
}

type Scope map[string]ScopeVariable

// set a key/value within a scope
func (s Scope) put(token token.Token, defined bool) {
	s[token.Text] = ScopeVariable{
		token:   token,
		defined: defined,
	}
}

// has returns whether a variable with the given
// name is declared and defined in this scope
func (s Scope) has(name string) (declared, defined bool) {
	v, ok := s[name]
	if !ok {
		return false, false
	}
	return true, v.defined
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
	r.ResolveStatements(statement.Body)
	r.endScope()
	return nil
}

func (r *Resolver) VisitActionStatement(statement tree.ActionStatement) interface{} {
	return nil
}

func (r *Resolver) VisitExpressionStatement(statement tree.ExpressionStatement) interface{} {
	r.resolveExpression(statement.Expression)
	return nil
}

func (r *Resolver) VisitLiteralExpr(expr tree.Literal) interface{} {
	return nil
}

func (r *Resolver) VisitVariableExpr(expr tree.VariableExpression) interface{} {
	r.resolveLocal(expr, expr.Name)
	return nil
}

func (r *Resolver) resolveExpression(expression tree.Expression) {
	expression.Accept(r)
}

func (r *Resolver) resolveStatement(statement tree.Statement) {
	statement.Accept(r)
}

// "We start at the innermost scope and work outwards, looking in each map for a
// matching name. If we find the variable, we resolve it, passing in the number of
// scopes between the current innermost scope and the scope where the variable was
// found. So, if the variable was found in the current scope, we pass in 0. If it’s
// in the immediately enclosing scope, 1"
func (r *Resolver) resolveLocal(expr tree.Expression, name token.Token) {
	for i := r.Scopes.size() - 1; i >= 0; i-- {
		s := r.Scopes.get(i)
		if _, defined := s.has(name.Text); defined {
			depth := r.Scopes.size() - 1 - i
			r.Interpreter.Resolve(expr, depth)
			// s.use(name.Lexeme)
			return
		}
	}
}

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
	r.Scopes.peek().put(name, false)
}

func (r *Resolver) define(name token.Token) {
	if r.Scopes.isEmpty() {
		return
	}
	r.Scopes.peek().put(name, true)
}
