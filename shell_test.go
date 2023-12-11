/******************************************************************************
 * fifo_test.go
 * Author:
 * Usage:    `go test`  or  `go test -v`
 * Description:
 *    An incomplete unit testing suite for fifo.go. You are welcome to change
 *    anything in this file however you would like. You are strongly encouraged
 *    to create additional tests for your implementation, as the ones provided
 *    here are extremely basic, and intended only to demonstrate how to test
 *    your program.
 ******************************************************************************/

 package 
 import (
	 "fmt"
	 "testing"
 )
 
 /******************************************************************************/
 /*                                Constants                                   */
 /******************************************************************************/
 // Constants can go here
 
 /******************************************************************************/
 /*                                  Tests                                     */
 /******************************************************************************/
 
 func TestFIFO(t *testing.T) {
	 // capacity := 16
	 // capacity2 := 1024
	 // fifo := NewFifo(capacity)
	 // fifo2 := NewFifo(capacity2)
	 // fifo3 := NewFifo(capacity2)
	 // checkCapacity(t, fifo, capacity)
	 // checkCapacity(t, fifo2, capacity2)
 
	 // for i := 0; i < 4; i++ {
	 // 	key := fmt.Sprintf("key%d", i)
	 // 	val := []byte(key)
	 // 	ok := fifo.Set(key, val)
	 // 	fmt.Printf("\n The key is %v and the value is %v \n", key, val)
	 // 	fmt.Printf("\n The current is %v and the size is %v \n", fifo.current, fifo.limit)
	 // 	if !ok {
	 // 		t.Errorf("Failed to add binding with key: %s", key)
	 // 		t.FailNow()
	 // 	}
 
	 // 	res, _ := fifo.Get(key)
	 // 	if !bytesEqual(res, val) {
	 // 		t.Errorf("Wrong value %s for binding with key: %s", res, key)
	 // 		t.FailNow()
	 // 	}
	 // }
 
	 // ok := fifo2.Set("key", []byte("old"))
	 // if !ok {
	 // 	t.Errorf("problem with setting the key")
	 // }
	 // val, exists := fifo2.Get("key")
	 // if exists {
	 // 	fmt.Printf("the key is %v \n", val)
	 // } else {
	 // 	t.Errorf("Failed to get key")
	 // }
 
	 // fmt.Printf("max storage is %d\n", fifo2.MaxStorage())
	 // fmt.Printf("remaining storage is %d\n", fifo2.RemainingStorage())
 
	 // ok = fifo2.Set("key", []byte("nw"))
	 // if !ok {
	 // 	t.Errorf("problem with setting the key")
	 // }
	 // val, exists = fifo2.Get("key")
	 // if exists {
	 // 	fmt.Printf("the value is %v \n", val)
	 // } else {
	 // 	t.Errorf("Failed to get key")
	 // }
 
	 // fmt.Printf("max storage is %d\n", fifo2.MaxStorage())
	 // fmt.Printf("remaining storage is %d\n", fifo2.RemainingStorage())
 
	 // ok = fifo3.Set("____1", []byte("____1"))
	 // val, _ = fifo3.Get("____1")
	 // val, _ = fifo3.Get("miss")
	 // fmt.Println(fifo3.stats.Hits)
	 // fmt.Println(fifo3.stats.Misses)
 
	 // // next test
	 // fmt.Println("new tests")
	 // fifo4 := NewFifo(100)
	 // fifo4.Set("____0", []byte("____0"))
	 // fifo4.Set("____1", []byte("____1"))
	 // fifo4.Set("____2", []byte("____2"))
	 // fifo4.Set("____3", []byte("____3"))
	 // fifo4.Set("____4", []byte("____4"))
	 // fifo4.Set("____5", []byte("____5"))
	 // fifo4.Set("____6", []byte("____6"))
	 // fifo4.Set("____7", []byte("____7"))
	 // fifo4.Set("____8", []byte("____8"))
	 // fifo4.Set("____9", []byte("____9"))
	 // fifo4.Set("____10", []byte("____a"))
	 // fmt.Printf("remaining storage = %v\n", fifo4.RemainingStorage())
	 // fifo4.Get("____0")
	 // fmt.Printf("misses are %v\n", fifo4.stats.Misses)
	 // fmt.Printf("hits are %v\n", fifo4.stats.Hits)
	 // // fmt.Println(fifo.order.Back().Value.([]byte))
	 // // fifo4.Get("____0")
	 // // fifo4.Get("____1")
 
	 // fifo5 := NewFifo(10)
	 // ok = fifo5.Set("😂", []byte("🙈"))
	 // fmt.Println(fifo5.Get("😂"))
	 // fmt.Printf("remaining is %v\n", fifo5.RemainingStorage())
	 // ok = fifo5.Set("12", []byte("12"))
 
	 // if !ok {
	 // 	t.Errorf("Failed to be okay")
	 // }
 
	 fifo6 := NewFifo(30)
	 fmt.Println(fifo6.Len())
	 fifo6.Set("____0", []byte("____0"))
	 fmt.Printf("remainstorage is %v\n", fifo6.RemainingStorage()) // 20
	 fifo6.Set("____1", []byte("____1"))
	 fmt.Printf("remainign is %v\n", fifo6.RemainingStorage()) // 10
	 fifo6.Set("____2", []byte("____2"))
	 fmt.Printf("remaining is %v\n", fifo6.RemainingStorage()) // 0
	 fifo6.Set("____3", []byte("____3"))
	 fmt.Printf("remaining is%v\n", fifo6.RemainingStorage())
	 fifo6.Set("____4", []byte("____4"))
	 fmt.Printf("length is %v", fifo6.Len())
 
 }


 // main_test.go

package main

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"
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
