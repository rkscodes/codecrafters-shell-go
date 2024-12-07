package main

import (
	"bufio"
	"fmt"
	"os"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

func main() {
	// Uncomment this block to pass the first stage

	// Wait for user input

	for {
		fmt.Fprint(os.Stdout, "$ ")
		var input, _ = bufio.NewReader(os.Stdin).ReadString('\n')
		fmt.Fprint(os.Stdout, input[:len(input)-1]+": command not found\n")
	}

}
