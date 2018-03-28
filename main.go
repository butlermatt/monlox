package main

import (
	"os/user"
	"fmt"
	"os"
	"github.com/butlermatt/monlox/repl"
)

func main() {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello %s! This is the Monlox programming language!\n", usr.Username)
	fmt.Printf("Feel free to type in commands\n")
	repl.Start(os.Stdin, os.Stdout)
	fmt.Printf("Good Byte!")
}
