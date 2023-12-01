package main

import (
	"bytes"
	"io"
	"os"
	"sync"
	"testing"
)

func TestShell_ExecuteCommand(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		// Add test cases here based on your specific requirements
		{"TestCommandExecution", "echo Hello, Go Shell!", "Hello, Go Shell!\n"},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Redirect os.Stdout to capture the output
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Create a new shell instance
			shell := &Shell{
				backgroundProcesses: sync.WaitGroup{},
			}

			// Execute the command
			shell.executeCommand(tt.input)

			// Close the write end of the pipe and restore os.Stdout
			w.Close()
			os.Stdout = oldStdout

			// Read the captured output from the pipe
			var capturedOutput bytes.Buffer
			io.Copy(&capturedOutput, r)

			// Compare the captured output with the expected output
			if got := capturedOutput.String(); got != tt.expected {
				t.Errorf("Expected output: %s, but got: %s", tt.expected, got)
			}
		})
	}
}

func TestShell_Run(t *testing.T) {
	// Create a pipe to capture user input
	r, w, _ := os.Pipe()
	os.Stdin = r
	defer func() {
		os.Stdin = os.Stdin
	}()

	// Redirect os.Stdout to capture the output
	oldStdout := os.Stdout
	rOut, wOut, _ := os.Pipe()
	os.Stdout = wOut
	defer func() {
		os.Stdout = oldStdout
	}()

	// Create a new shell instance
	shell := &Shell{
		backgroundProcesses: sync.WaitGroup{},
	}

	// Run the shell in a goroutine
	go shell.Run()

	// Write user input to the pipe
	io.WriteString(w, "echo Hello, Go Shell!\n")
	// Close the write end of the pipe
	w.Close()

	// Wait for the shell to complete
	shell.backgroundProcesses.Wait()

	// Close the write end of the output pipe and restore os.Stdout
	wOut.Close()
	os.Stdout = oldStdout

	// Read the captured output from the pipe
	var capturedOutput bytes.Buffer
	io.Copy(&capturedOutput, rOut)

	// Compare the captured output with the expected output
	expected := "Hello, Go Shell!\n"
	if got := capturedOutput.String(); got != expected {
		t.Errorf("Expected output: %s, but got: %s", expected, got)
	}
}
