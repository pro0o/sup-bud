package main

import (
	"fmt"
	"main/repl"
	"os"
)

func main() {
	fmt.Printf("hello patlu, this is paltu programming langtang! \n")
	fmt.Printf("type-in any cmds \n")
	repl.Start(os.Stdin, os.Stdout)
}
