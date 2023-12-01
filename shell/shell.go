package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
)

// Shell structure
type Shell struct {
	backgroundProcesses sync.WaitGroup
}

// Run the shell
func (s *Shell) Run() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print(">>> ")
		scanner.Scan()
		input := scanner.Text()

		if input == "exit" {
			// Wait for background processes to complete before exiting
			s.backgroundProcesses.Wait()
			break
		}

		go s.executeCommand(input)
	}
}

// Execute a command
func (s *Shell) executeCommand(input string) {
	// Split input into command and arguments
	args := strings.Fields(input)

	// Check if the last argument is "&" for background execution
	background := false
	if len(args) > 0 && args[len(args)-1] == "&" {
		background = true
		args = args[:len(args)-1] // Remove "&" from arguments
	}

	if len(args) == 0 {
		return // Skip execution and continue with the next iteration of the loop
	}
	cmd := exec.Command(args[0], args[1:]...)

	// Redirect input, output, and error streams
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if background {
		// If background execution, increment the counter and start a new goroutine
		s.backgroundProcesses.Add(1)
		go func() {
			defer s.backgroundProcesses.Done()
			err := cmd.Run()
			if err != nil {
				fmt.Println("Error:", err)
			}
		}()
		fmt.Println("Background process started:", cmd.Process.Pid)
	} else {
		// If foreground execution, wait for the command to complete
		err := cmd.Run()
		if err != nil {
			fmt.Println("Error:", err)
		}
	}
}

func main() {
	shell := &Shell{}
	shell.Run()
}
