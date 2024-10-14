package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// Open the input file
	inputFile, err := os.Open("largefile.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer inputFile.Close()

	// Open the output file
	outputFile, err := os.Create("outputfile.txt")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer outputFile.Close()

	// Create a buffered reader and writer
	reader := bufio.NewReader(inputFile)
	writer := bufio.NewWriter(outputFile)

	// Buffer size for reading
	buffer := make([]byte, 4096) // 4KB buffer

	for {
		// Read into the buffer
		n, err := reader.Read(buffer)
		if err != nil {
			if err != os.EOF {
				fmt.Println("Error reading file:", err)
			}
			break
		}

		// Write the buffer content to the output file
		_, err = writer.Write(buffer[:n])
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	}

	// Flush the writer buffer to ensure all data is written
	err = writer.Flush()
	if err != nil {
		fmt.Println("Error flushing writer:", err)
		return
	}

	fmt.Println("Finished processing")
}
