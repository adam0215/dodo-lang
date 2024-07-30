package main

import (
	"dodo-lang/repl"
	"fmt"
	"os"
	"os/user"
)

func main() {
	user, err := user.Current()

	if err != nil {
		panic(err)
	}

	// REPL
	fmt.Printf("\nHello %s! You are now running the Dodo programming language.\n", user.Username)
	fmt.Printf("Now, type some commands below :)\n\n")
	repl.Start(os.Stdin, os.Stdout)
}
