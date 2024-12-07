package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

func main() {
	// Uncomment this block to pass the first stage

	// Wait for user input

	for {
		fmt.Fprint(os.Stdout, "$ ")
		var input, error = bufio.NewReader(os.Stdin).ReadString('\n')
		if error != nil {
			fmt.Fprint(os.Stderr, "Failed to read from stdin")
		}
		// Seprate frist part of the comman and rest
		parts := strings.SplitN(input, " ", 2)
		command := parts[0]
		args := ""
		if len(parts) > 1 {
			args = parts[1]
		}
		command = strings.TrimSuffix(strings.TrimSpace(command), "\n")
		args = strings.TrimSuffix(strings.TrimSpace(args), "\n")
		execute_command(command, args)
	}

}

func execute_command(command string, args string) {
	// fmt.Fprint(os.Stdout, command)
	switch command {
	case "echo":
		fmt.Fprint(os.Stdout, args+"\n")
		break

	case "exit":
		if len(args) > 0 && args == "0" {
			os.Exit(0)
		}
		break

	default:
		fmt.Fprint(os.Stdout, command+": command not found\n")

	}

}
