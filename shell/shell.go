package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/chzyer/readline"
)

// Shell structure
type Shell struct {
	backgroundProcesses sync.WaitGroup
	currentDirectory    string
}

// Run the shell
func (s *Shell) Run() {
	// Set up signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTSTP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGHUP)

	// Set up a channel to notify when a command is completed
	cmdCompleteChan := make(chan bool)

	// Handle signals asynchronously
	go func() {
		for {
			sig := <-sigChan
			switch sig {
			case syscall.SIGINT:
				fmt.Println("\nReceived SIGINT. Press Enter to continue.")
			case syscall.SIGTSTP:
				fmt.Println("\nReceived SIGTSTP. Press Enter to continue.")
			case syscall.SIGQUIT:
				fmt.Println("\nReceived SIGQUIT.")
				// Handle SIGQUIT (e.g., cleanup, graceful exit)
			case syscall.SIGTERM:
				fmt.Println("\nReceived SIGTERM.")
				// Handle SIGTERM (e.g., cleanup, graceful exit)
			case syscall.SIGHUP:
				fmt.Println("\nReceived SIGHUP.")
				// Handle SIGHUP (e.g., configuration reload)
			}
		}
	}()

	rl, err := readline.NewEx(&readline.Config{
		Prompt:       ">> ",
		AutoComplete: getAutoComplete(),
	})

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
		fmt.Println(input)
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
		input = strings.ReplaceAll(input, "'", "")
		go func() {
			s.executeCommand(input)
			// Notify that the command is completed
			cmdCompleteChan <- true
		}()

		// Wait for the command to complete or for a signal to be received
		select {
		case <-cmdCompleteChan:
			// Continue to the next iteration
		case sig := <-sigChan:
			// Handle signals received while waiting for the command to complete
			switch sig {
			case syscall.SIGINT:
				fmt.Println("\nReceived SIGINT during command execution.")
			case syscall.SIGTSTP:
				fmt.Println("\nReceived SIGTSTP during command execution.")
			case syscall.SIGQUIT:
				fmt.Println("\nReceived SIGQUIT during command execution.")
			case syscall.SIGTERM:
				fmt.Println("\nReceived SIGTERM during command execution.")
			case syscall.SIGHUP:
				fmt.Println("\nReceived SIGHUP during command execution.")
			}
			// Wait for the command to complete before continuing
			<-cmdCompleteChan
		}
	}
}

// Start a new shell
func (s *Shell) startNewShell() {

	cmd := exec.Command(os.Args[0]) // Restart the current program

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error starting new shell:", err)
	}
}

func (s *Shell) executeCommand(input string) {
	// Trim newline character from input
	input = strings.TrimSuffix(input, "\n")

	// Check if the last argument is "&" for background execution
	background := false
	if strings.HasSuffix(input, "&") {
		background = true
		input = strings.TrimSuffix(input, "&") // Remove "&" from input
	}

	if len(input) == 0 {
		return // Skip execution and continue with the next iteration of the loop
	}

	// Check for piping
	var pipeIndex int
	args := strings.Fields(input)
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

	var cmd *exec.Cmd
	if background {
		cmd = exec.Command("bash", "-c", input)
	} else {
		cmd = exec.Command("bash", "-c", input)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	if background {
		// If background execution, increment the counter and start a new goroutine
		s.backgroundProcesses.Add(1)
		go func(cmd *exec.Cmd) { // Capture cmd variable
			defer s.backgroundProcesses.Done()
			err := cmd.Run()
			if err != nil {
				fmt.Println("Error:", err)
			}
		}(cmd)
		if cmd.Process != nil {
			fmt.Println("Background process started:", cmd.Process.Pid)
		} else {
			fmt.Println("Background process started.")
		}
	} else {
		// If foreground execution, wait for the command to complete
		err := cmd.Run()
		if err != nil {
			fmt.Println("Error:", err)
		}
	}
}

func getAutoComplete() readline.AutoCompleter {
	completer := readline.NewPrefixCompleter(
		readline.PcItem("newshell"),
		readline.PcItem("exit"),
		readline.PcItem("cd"),
		readline.PcItem("your_custom_command_here"),
		readline.PcItemDynamic(getDynamicCompleteFunc()), // Add this line for dynamic completion
	)
	return completer
}

func getDynamicCompleteFunc() readline.DynamicCompleteFunc {
	return func(prefix string) []string {
		files, err := os.ReadDir(".")
		if err != nil {
			fmt.Println("Error reading directory:", err)
			return nil
		}

		var completions []string
		for _, file := range files {
			if strings.HasPrefix(file.Name(), prefix) {
				completions = append(completions, file.Name())
			}
		}
		return completions
	}
}

// Run piped commands
func (s *Shell) runPipedCommands(args []string, pipeIndex int) {
	// Join the first part of the command
	cmd1 := strings.Join(args[:pipeIndex], " ")

	// Join the second part of the command
	cmd2 := strings.Join(args[pipeIndex+1:], " ")

	// Use bash to interpret the commands correctly
	cmd := exec.Command("bash", "-c", cmd1+" | "+cmd2)

	// Redirect output and error streams
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	// Start the command and wait for it to complete
	if err := cmd.Run(); err != nil {
		fmt.Println("Error running piped commands:", err)
		return
	}
}

// // Run piped commands
// func (s *Shell) runPipedCommands(args []string, pipeIndex int) {
// 	// Split the command into two parts
// 	cmd1 := exec.Command(args[0], args[1:pipeIndex]...)
// 	cmd2 := exec.Command(args[pipeIndex+1], args[pipeIndex+2:]...)

// 	// Create a pipe to connect the output of cmd1 to the input of cmd2
// 	pipeReader, pipeWriter := io.Pipe()
// 	cmd1.Stdout = pipeWriter
// 	cmd2.Stdin = pipeReader

// 	// Redirect output and error streams
// 	cmd1.Stderr = os.Stderr
// 	cmd2.Stdout = os.Stdout
// 	cmd2.Stderr = os.Stderr

// 	// Use a WaitGroup to wait for both commands to complete
// 	var wg sync.WaitGroup
// 	wg.Add(2)

// 	// Start cmd1
// 	go func() {
// 		defer wg.Done()
// 		if err := cmd1.Run(); err != nil {
// 			fmt.Println("Error running command 1:", err)
// 		}
// 		// Close the writer end of the pipe after cmd1 completes
// 		pipeWriter.Close()
// 	}()

// 	// Start cmd2
// 	go func() {
// 		defer wg.Done()
// 		if err := cmd2.Run(); err != nil {
// 			fmt.Println("Error running command 2:", err)
// 		}
// 	}()

// 	// Wait for both commands to complete
// 	wg.Wait()
// }

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
