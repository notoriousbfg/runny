package parser

import "runny/src/token"

func New(tokens []token.Token) *Parser {
	return &Parser{
		Tokens:  tokens,
		Current: 0,
	}
}

type Parser struct {
	Tokens  []token.Token
	Current int
}

// func (p *Parser) Parse() ([]tree.Stmt, error) {
// 	statements := make([]tree.Stmt, 0)
// 	for !p.isAtEnd() {
// 		statements = append(statements, p.Declaration())
// 	}
// 	return statements
// }
