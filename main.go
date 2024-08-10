package main

import (
	"dodo-lang/repl"
	"flag"
	"os"
)

var filename string
var verbose bool

func init() {
	flag.StringVar(&filename, "f", "", "Dodo file to run")
	flag.BoolVar(&verbose, "v", false, "Verbose mode")
	flag.Parse()
}

func main() {
	if filename != "" {
		// REPL / File mode
		repl.FileMode(os.Stdin, os.Stdout, filename, verbose)
	} else {
		// REPL / Interactive mode
		repl.InteractiveMode(os.Stdin, os.Stdout, verbose)
	}
}
