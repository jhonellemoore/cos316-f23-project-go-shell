package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestPiping(t *testing.T) {
	// Create a simple pipeline: ls -l | grep "output"
	pipeline := "ls | cat"

	// Split the pipeline into individual commands
	commands := strings.Split(pipeline, "|")

	// Process each command in the pipeline
	var cmd *exec.Cmd
	var lastCmd *exec.Cmd

	for _, command := range commands {
		args := strings.Fields(strings.TrimSpace(command))

		// Create a command
		cmd = exec.Command(args[0], args[1:]...)

		// If not the first command, set the input of this command to the output of the previous command
		if lastCmd != nil {
			cmd.Stdin, _ = lastCmd.StdoutPipe()
		}

		// Set the output of this command as the input for the next command
		cmd.Stdout = os.Stdout

		// Start the command
		if err := cmd.Start(); err != nil {
			t.Fatalf("Error starting command: %v", err)
		}

		// Wait for the command to finish
		if err := cmd.Wait(); err != nil {
			t.Fatalf("Error waiting for command: %v", err)
		}

		// Update the last command
		lastCmd = cmd
	}

	// Check the exit status of the last command
	if cmd.ProcessState.ExitCode() != 0 {
		t.Fatalf("Last command failed with exit code %d", cmd.ProcessState.ExitCode())
	}

	fmt.Println("Pipeline execution successful.")
}
