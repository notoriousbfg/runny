package interpreter

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
)

// not a statement in a "tree" sense
// it just seemed the most appropriate word
type Statement struct {
	Cmd    *exec.Cmd
	StdOut io.ReadCloser
	StdErr io.ReadCloser
}

type Printer struct {
	Statements []Statement
}

func (p *Printer) Print() {
	for _, stmt := range p.Statements {
		content, err := io.ReadAll(stmt.StdOut)
		if err != nil {
			panic(p.error(err.Error()))
		}
		fmt.Print(string(content))
		stmt.StdOut.Close()
		if stmt.Cmd != nil {
			err := stmt.Cmd.Wait()
			if err != nil {
				panic(p.error(err.Error()))
			}
		}
	}
}

func (p *Printer) Push(stmt Statement) {
	p.Statements = append(p.Statements, stmt)
}

func (p *Printer) PushStr(str string) {
	p.Statements = append(p.Statements, statementStr(str))
}

func (i *Printer) error(message string) *RuntimeError {
	err := &RuntimeError{
		Message: fmt.Sprintf("runtime error: %s\n", message),
	}
	return err
}

func statementStr(str string) Statement {
	readCloser := io.NopCloser(strings.NewReader(str))
	return Statement{
		StdOut: readCloser,
	}
}
