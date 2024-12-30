package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

type commandFunction func(InputPrompt)

type InputPrompt struct {
	command string
	args    []string
}

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

func echo(inputPrompt InputPrompt) {
	args := inputPrompt.args
	joinedArgs := strings.Join(args, " ")
	fmt.Fprint(os.Stdout, joinedArgs+"\n")
}

func exit(inputPrompt InputPrompt) {
	args := inputPrompt.args
	firstArg := "1"
	if len(args) > 0 {
		firstArg = args[0]
	}
	if len(args) > 0 && firstArg == "0" {
		os.Exit(0)
	}
}

func type_(inputPrompt InputPrompt) {
	// command := inputPrompt.command
	args := inputPrompt.args
	if len(args) > 1 {
		log.Fatal("more than 1 args")
		return
	}
	operandCommand := args[0]
	if _, ok := shellBuiltIn[operandCommand]; ok {
		fmt.Fprint(os.Stdout, operandCommand+" is a shell builtin\n")
		return
	}

	if commandPath, err := findCommand(operandCommand); err == nil {
		fmt.Fprintf(os.Stdout, "%s\n", commandPath)
		return
	}

	//If found nowhere print this by default
	fmt.Fprint(os.Stdout, operandCommand+": not found\n")
}

func pwd(inputPrompt InputPrompt) {
	args := inputPrompt.args

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

func cd(inputPrompt InputPrompt) {
	args := inputPrompt.args
	if len(args) > 1 {
		fmt.Fprint(os.Stdout, "pwd: too many arguments\n")
		return
	}
	path := args[0]
	if path == "~" || path == "" {
		err := os.Chdir(os.Getenv("HOME"))
		if err != nil {
			fmt.Fprintf(os.Stdout, "$HOME variable not set\n")
		}
		return
	}

	err := os.Chdir(path)
	if err != nil {
		fmt.Fprintf(os.Stdout, "cd: %s: No such file or directory\n", path)
	}
}

func findCommand(command string) (string, error) {
	pathEnv := os.Getenv("PATH")
	// Split PATH into directories
	paths := strings.Split(pathEnv, ":")

	// Search each directory for the command
	for _, path := range paths {
		commandPath := path + "/" + command
		if _, err := os.Stat(commandPath); err == nil {
			return commandPath, nil
		}
	}
	return "", fmt.Errorf("command not found")
}

// TODO: let's use []string for arguments passing to better implement the quotes additions
func main() {
	// Uncomment this block to pass the first stage

	// Wait for user input

	for {
		fmt.Fprint(os.Stdout, "$ ")
		var input, error = bufio.NewReader(os.Stdin).ReadString('\n')
		if error != nil {
			fmt.Fprint(os.Stderr, "Failed to read from stdin")
			return
		}
		// Seprate frist part of the comman and rest
		caa, err := cmdAndArgs(input)
		if err != nil {
			log.Fatal(err)
		}
		execute_command(caa)
	}

}

func execute_command(caa InputPrompt) {
	// fmt.Fprint(os.Stdout, command)
	command := caa.command
	args := caa.args

	if function, ok := shellBuiltIn[command]; ok {
		function(caa)
		return
	}

	if commandPath, err := findCommand(command); err == nil {
		var cmd *exec.Cmd
		if len(args) > 0 {
			cmd = exec.Command(commandPath, args...)
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

// Not satisfied with this func seems bit long and complex
func cmdAndArgs(input string) (inputPrompt InputPrompt, err error) {
	parts := strings.SplitN(input, " ", 2)
	command := strings.TrimSpace(parts[0])

	//there is no args
	if len(parts) < 2 {
		return InputPrompt{command: command}, nil
	}

	sinput := strings.Split(parts[1], "")
	args := []string{}
	remains := ""
	for i := 0; i < len(sinput); i++ {
		if sinput[i] == "'" {
			s := ""

			for i < len(sinput) {
				i++
				if i >= len(sinput) {
					i--
					break
				}
				if sinput[i] == "'" {
					args = append(args, s)
					break
				}
				s = s + sinput[i]
				// throw and error here if ' is not matched

			}
			if sinput[i] != "'" {
				return InputPrompt{}, errors.New("' not found")
			}
		} else if sinput[i] == " " {
			remains_split := strings.Split(remains, " ")
			for _, ele := range remains_split {
				ele = strings.TrimSpace(ele)
				if ele != "" {
					args = append(args, ele)
				}
			}
			remains = ""
		} else {
			remains += sinput[i]
		}
	}

	remains_split := strings.Split(remains, " ")
	for _, ele := range remains_split {
		ele = strings.TrimSpace(ele)
		if ele != "" {
			args = append(args, ele)
		}
	}

	ip := InputPrompt{
		command: command,
		args:    args,
	}
	return ip, nil
}
