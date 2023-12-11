package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/chzyer/readline"
)

// Shell structure
type Shell struct {
	backgroundProcesses sync.WaitGroup
	currentDirectory    string
}

// Run the shell
func (s *Shell) Run() {
	rl, err := readline.New("> ")
	if err != nil {
		fmt.Println("Error creating readline instance:", err)
		return
	}
	defer rl.Close()

	for {
		line, err := rl.Readline()
		if err != nil {
			fmt.Println("Error reading input:", err)
			break
		}

		input := strings.TrimSpace(line)
		if input == "exit" {
			// Wait for background processes to complete before exiting
			s.backgroundProcesses.Wait()
			break
		}
		if input == "newshell" {
			// Wait for background processes to complete before starting a new shell
			s.backgroundProcesses.Wait()
			s.startNewShell()
			break
		}

		// Execute command synchronously to ensure proper input handling
		s.executeCommand(input)
	}
}

// Start a new shell
func (s *Shell) startNewShell() {
	// You can use os/exec to start a new shell process
	cmd := exec.Command(os.Args[0]) // Restart the current program

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error starting new shell:", err)
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

	// Check for piping
	var pipeIndex int
	for i, arg := range args {
		if arg == "|" {
			pipeIndex = i
			break
		}
	}

	if pipeIndex > 0 && pipeIndex < len(args)-1 {
		s.runPipedCommands(args, pipeIndex)
		return
	}

	// Check for output redirection
	var outputRedirectionIndex int
	for i, arg := range args {
		if arg == ">" {
			outputRedirectionIndex = i
			break
		}
	}

	if outputRedirectionIndex > 0 && outputRedirectionIndex < len(args)-1 {
		s.runCommandWithOutputRedirection(args, outputRedirectionIndex)
		return
	}

	if args[0] == "cd" {
		s.changeDirectory(args[1:])
		return
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

// Run piped commands
func (s *Shell) runPipedCommands(args []string, pipeIndex int) {
	cmd1 := exec.Command(args[0], args[1:pipeIndex]...)
	cmd2 := exec.Command(args[pipeIndex+1], args[pipeIndex+2:]...)

	// Create a pipe to connect the output of cmd1 to the input of cmd2
	pipeReader, pipeWriter := io.Pipe()
	cmd1.Stdout = pipeWriter
	cmd2.Stdin = pipeReader

	// Redirect output and error streams
	cmd1.Stderr = os.Stderr
	cmd2.Stdout = os.Stdout
	cmd2.Stderr = os.Stderr

	// Use a WaitGroup to wait for both commands to complete
	var wg sync.WaitGroup
	wg.Add(2)

	// Start cmd1
	go func() {
		defer wg.Done()
		if err := cmd1.Start(); err != nil {
			fmt.Println("Error starting command:", err)
			return
		}
		// Close the writer end of the pipe after starting cmd1
		pipeWriter.Close()
	}()

	// Start cmd2
	go func() {
		defer wg.Done()
		if err := cmd2.Run(); err != nil {
			fmt.Println("Error running command:", err)
		}
	}()

	// Wait for both commands to complete
	wg.Wait()
}

// Run command with output redirection
func (s *Shell) runCommandWithOutputRedirection(args []string, outputRedirectionIndex int) {
	cmd := exec.Command(args[0], args[1:outputRedirectionIndex]...)

	// Open a file for writing (create or overwrite)
	outputFile, err := os.Create(args[outputRedirectionIndex+1])
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer outputFile.Close()

	// Redirect output and error streams to the file
	cmd.Stdout = outputFile
	cmd.Stderr = os.Stderr

	// Start the command and wait for it to complete
	if err := cmd.Run(); err != nil {
		fmt.Println("Error running command:", err)
		return
	}
}

// Change current directory
func (s *Shell) changeDirectory(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: cd <directory>")
		return
	}

	dir := args[0]
	err := os.Chdir(dir)
	if err != nil {
		fmt.Println("Error changing directory:", err)
		return
	}

	s.currentDirectory, err = os.Getwd()
	if err != nil {
		fmt.Println("Error getting working directory:", err)
		return
	}

	fmt.Println("Changed directory to:", s.currentDirectory)
}

func main() {
	shell := &Shell{}
	shell.Run()
}
