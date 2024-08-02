package main

import (
	"dodo-lang/repl"
	"flag"
	"fmt"
	"os"
	"os/user"
)

var filename string
var verbose bool

func init() {
	flag.StringVar(&filename, "f", "", "Dodo file to run")
	flag.BoolVar(&verbose, "v", false, "Verbose mode")
	flag.Parse()
}

func main() {
	user, err := user.Current()

	if err != nil {
		panic(err)
	}

	if filename != "" {
		// REPL / File mode
		repl.FileMode(os.Stdin, os.Stdout, filename, verbose)
	} else {
		// REPL / Interactive mode
		fmt.Printf("\nHello %s! You are now running the Dodo programming language.\n", user.Username)
		fmt.Printf("Now, type some commands below :)\n\n")
		repl.InteractiveMode(os.Stdin, os.Stdout, verbose)
	}
}
