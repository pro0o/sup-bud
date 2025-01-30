package main

import (
	"fmt"
	"olaf/repl"
	"os"
)

const ascii_image = ` ▒█████   ██▓    ▄▄▄        █████▒
▒██▒  ██▒▓██▒   ▒████▄    ▓██   ▒ 
▒██░  ██▒▒██░   ▒██  ▀█▄  ▒████ ░ 
▒██   ██░▒██░   ░██▄▄▄▄██ ░▓█▒  ░ 
░ ████▓▒░░██████▒▓█   ▓██▒░▒█░    
░ ▒░▒░▒░ ░ ▒░▓  ░▒▒   ▓▒█░ ▒ ░    
  ░ ▒ ▒░ ░ ░ ▒  ░ ▒   ▒▒ ░ ░      
░ ░ ░ ▒    ░ ░    ░   ▒    ░ ░    
    ░ ░      ░  ░     ░  ░        
                                  `

func main() {

	fmt.Printf(ascii_image)
	fmt.Printf("\n")
	// fmt.Printf("Type-in any cmds: \n")
	fmt.Printf("Welcome to olaf! A toy language in golang. \n")
	repl.Start(os.Stdin, os.Stdout)
}
