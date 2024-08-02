package repl

import (
	"bufio"
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

func InteractiveMode(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()

		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)

		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func FileMode(in io.Reader, out io.Writer, filename string) {
	content, err := os.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	l := lexer.New(string(content))
	p := parser.New(l)

	program := p.ParseProgram()
	env := object.NewEnvironment()

	if len(p.Errors()) != 0 {
		printParserErrors(out, p.Errors())
	}

	evaluated := evaluator.Eval(program, env)

	if evaluated != nil {
		io.WriteString(out, evaluated.Inspect())
		io.WriteString(out, "\n")
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
