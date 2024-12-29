package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

type commandFunction func(string)

var shellBuiltIn map[string]commandFunction

func init() {
	shellBuiltIn = map[string]commandFunction{
		"echo": echo,
		"exit": exit,
		"type": type_,
		"pwd":  pwd,
		"cd":   cd,
	}
}

func echo(args string) {
	fmt.Fprint(os.Stdout, args+"\n")
}

func exit(args string) {
	if len(args) > 0 && args == "0" {
		os.Exit(0)
	}
}

func type_(args string) {
	if _, ok := shellBuiltIn[args]; ok {
		fmt.Fprint(os.Stdout, args+" is a shell builtin\n")
		return
	}

	if commandPath, err := findCommand(args); err == nil {
		fmt.Fprintf(os.Stdout, "%s\n", commandPath)
		return
	}

	//If found nowhere print this by default
	fmt.Fprint(os.Stdout, args+": not found\n")
}

func pwd(args string) {
	if len(args) >= 1 {
		fmt.Fprint(os.Stdout, "pwd: too many arguments\n")
		return
	}

	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stdout, "failed to get current working dir")
		return
	}
	fmt.Fprint(os.Stdout, dir, "\n")
}

func cd(args string) {
	args = strings.TrimSpace(args)
	if args == "~" || args == "" {
		err := os.Chdir(os.Getenv("HOME"))
		if err != nil {
			fmt.Fprintf(os.Stdout, "$HOME variable not set\n")
		}
		return
	}

	err := os.Chdir(args)
	if err != nil {
		fmt.Fprintf(os.Stdout, "cd: %s: No such file or directory\n", args)
	}
}

func findCommand(args string) (string, error) {
	pathEnv := os.Getenv("PATH")
	// Split PATH into directories
	paths := strings.Split(pathEnv, ":")

	// Search each directory for the command
	for _, path := range paths {
		commandPath := path + "/" + args
		if _, err := os.Stat(commandPath); err == nil {
			return commandPath, nil
		}
	}
	return "", fmt.Errorf("command not found")
}

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

	if function, ok := shellBuiltIn[command]; ok {
		function(args)
		return
	}

	if commandPath, err := findCommand(command); err == nil {
		var cmd *exec.Cmd
		if len(args) > 0 {
			cmd = exec.Command(commandPath, args)
		} else {
			cmd = exec.Command(commandPath)
		}
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
		return
	}

	fmt.Fprint(os.Stdout, command+": command not found\n")
}
