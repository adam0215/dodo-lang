package main

import (
	"dodo-lang/repl"
	"flag"
	"fmt"
	"os"
	"os/user"
)

var filename string

func init() {
	flag.StringVar(&filename, "f", "", "Dodo file to run")
	flag.Parse()
}

func main() {
	user, err := user.Current()

	if err != nil {
		panic(err)
	}

	if filename != "" {
		// REPL / File mode
		repl.FileMode(os.Stdin, os.Stdout, filename)
	} else {
		// REPL / Interactive mode
		fmt.Printf("\nHello %s! You are now running the Dodo programming language.\n", user.Username)
		fmt.Printf("Now, type some commands below :)\n\n")
		repl.InteractiveMode(os.Stdin, os.Stdout)
	}
}
