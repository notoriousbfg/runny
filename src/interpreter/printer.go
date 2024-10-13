package interpreter

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
)

// not the same as a "tree" statement
// it just seemed the most appropriate word
type Statement struct {
	Cmd    *exec.Cmd
	StdOut io.ReadCloser
	StdErr io.ReadCloser
	Before func() []Statement
	After  func() []Statement
}

type Printer struct {
	Statements []Statement
}

func (p *Printer) Print() {
	for _, statement := range p.Statements {
		if statement.Before != nil {
			for _, before := range statement.Before() {
				p.printStatement(before)
			}
		}
		p.printStatement(statement)
		if statement.After != nil {
			for _, after := range statement.After() {
				p.printStatement(after)
			}
		}
	}
}

func (p *Printer) printStatement(statement Statement) {
	content, err := io.ReadAll(statement.StdOut)
	if err != nil {
		panic(p.error(err.Error()))
	}
	fmt.Print(string(content))
	statement.StdOut.Close()
	if statement.Cmd != nil {
		err := statement.Cmd.Wait()
		if err != nil {
			panic(p.error(err.Error()))
		}
	}
}

func (p *Printer) Push(statement Statement) {
	p.Statements = append(p.Statements, statement)
}

func (p *Printer) PushStr(str string) {
	p.Statements = append(p.Statements, StrStatement(str))
}

func (i *Printer) error(message string) *RuntimeError {
	err := &RuntimeError{
		Message: fmt.Sprintf("runtime error: %s\n", message),
	}
	return err
}

func StrStatement(str string) Statement {
	readCloser := io.NopCloser(strings.NewReader(str))
	return Statement{
		StdOut: readCloser,
	}
}
