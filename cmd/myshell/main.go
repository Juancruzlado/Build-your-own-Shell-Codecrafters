package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	// Initialize the reader once
	reader := bufio.NewReader(os.Stdin)

	// Print the initial prompt
	fmt.Print("$ ")

	for {
		// Wait for user input
		input, _ := reader.ReadString('\n')

		// Remove the trailing newline character from the input
		input = strings.TrimSpace(input)

		// Split the input into commands
		commands := strings.Fields(input)

		// Ensure there is at least one command
		if len(commands) == 0 {
			// Print the prompt again
			fmt.Print("$ ")
			continue
		}

		// Check for the exit command
		if commands[0] == "exit" && len(commands) > 1 && commands[1] == "0" {
			os.Exit(0)
		} else if commands[0] == "echo" {
			fmt.Printf("%s\n", strings.Join(commands[1:], " "))
		} else if commands[0] == "type" && len(commands) == 2 {
			// Handle the type command
			switch commands[1] {
			case "exit", "echo", "type":
				fmt.Printf("%s is a shell builtin\n", commands[1])
			default:
				// Search for the command in the PATH directories
				env := os.Getenv("PATH")
				paths := strings.Split(env, ":")
				found := false
				for _, path := range paths {
					exec := path + "/" + commands[1]
					if _, err := os.Stat(exec); err == nil {
						fmt.Printf("%s is %s\n", commands[1], exec)
						found = true
						break
					}
				}
				if !found {
					fmt.Printf("bash: type: %s: not found\n", commands[1])
				}
			}
		} else if commands[0] == "cd" {
			// Handle the cd command
			if len(commands) < 2 {
				fmt.Println("cd: missing argument")
			} else {
				path := commands[1]
				if path == "~" {
					// Replace ~ with the value of the HOME environment variable
					path = os.Getenv("HOME")
				}
				err := os.Chdir(path)
				if err != nil {
					fmt.Printf("cd: %s: No such file or directory\n", path)
				}
			}
		} else {
			// Attempt to execute the command
			cmd := exec.Command(commands[0], commands[1:]...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				fmt.Printf("%s: command not found\n", commands[0])
			}
		}

		// Print the prompt again
		fmt.Print("$ ")
	}
}
