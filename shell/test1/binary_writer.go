package main

import (
	"encoding/binary"
	"fmt"
	"os"
)

func main() {
	file, err := os.Create("sample.bin")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Create some data to write to the file
	data := []uint16{42, 123, 789, 456}

	// Write binary data to the file
	for _, value := range data {
		err := binary.Write(file, binary.LittleEndian, value)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	}
	fmt.Println("Binary file created successfully.")
}