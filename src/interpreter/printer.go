package interpreter

import (
	"bufio"
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
}

type Printer struct {
	Statements []Statement
}

func (p *Printer) Print() {
	for _, statement := range p.Statements {
		p.printStatement(statement)
	}
}

func (p *Printer) printStatement(statement Statement) {
	go func() {
		scanner := bufio.NewScanner(statement.StdOut)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			fmt.Println(err.Error())
		}
	}()
	go func() {
		scanner := bufio.NewScanner(statement.StdErr)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			fmt.Println(err.Error())
		}
	}()
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
