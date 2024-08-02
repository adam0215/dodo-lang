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
	"os/exec"
)

type ExecMode = string

const PROMPT = "\n>> "

func InteractiveMode(in io.Reader, out io.Writer, verbose bool) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	// disable input buffering
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	// do not display entered characters on the screen
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()

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

		if len(p.Errors()) != 0 && verbose {
			printParserErrors(out, p.Errors())
			continue
		}

		switch evaluated := evaluator.Eval(program, env).(type) {
		case *object.Error:
			if verbose && evaluated != nil {
				io.WriteString(out, fmt.Sprintf("%s", evaluated.Inspect()))
			}
		default:
			if evaluated != nil {
				io.WriteString(out, fmt.Sprintf("%s", evaluated.Inspect()))
			}
		}
	}
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
		printParserErrors(out, p.Errors())
	}

	switch evaluated := evaluator.Eval(program, env).(type) {
	case *object.Error:
		if verbose && evaluated != nil {
			io.WriteString(out, fmt.Sprintf("%s", evaluated.Inspect()))
		}
	default:
		if evaluated != nil {
			io.WriteString(out, fmt.Sprintf("%s", evaluated.Inspect()))
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
