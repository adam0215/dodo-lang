package repl

import (
	"dodo-lang/cli"
	"dodo-lang/evaluator"
	"dodo-lang/lexer"
	"dodo-lang/object"
	"dodo-lang/parser"
	"fmt"
	"io"
	"os"
)

type ExecMode = string

const PROMPT = "\n>> "

func InteractiveMode(in io.Reader, out io.Writer, verbose bool) {
	cli := cli.New(verbose)
	cli.Init()
}

func FileMode(in io.Reader, out io.Writer, filename string, verbose bool) {
	content, err := os.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	l := lexer.New(string(content))
	p := parser.New(l)

	program := p.ParseProgram()
	env := object.NewEnvironment()

	if len(p.Errors()) != 0 && verbose {
		p.PrintParserErrors(out)
	}

	switch evaluated := evaluator.Eval(program, env).(type) {
	case *object.Error:
		if verbose && evaluated != nil {
			io.WriteString(out, fmt.Sprintf("%s", evaluated.Inspect()))
		}
	}
}
