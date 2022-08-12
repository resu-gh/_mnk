package main

import (
	"fmt"
	"mnk/src/repl"
	"os"
	"os/user"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

    fmt.Printf("hello %s! this is the mnk programming languange!\n", user.Username)
    fmt.Printf("feel free to type in commands\n")
    repl.Start(os.Stdin, os.Stdout)
}
