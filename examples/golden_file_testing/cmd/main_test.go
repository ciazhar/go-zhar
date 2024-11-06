package main

import (
	"os"
	"testing"
)

func TestGenerateOutput(t *testing.T) {
	// Input
	input := "example"

	// Call the function
	output := GenerateOutput(input)

	// Read the golden file
	golden, err := os.ReadFile("../testdata/data.golden")
	if err != nil {
		t.Fatalf("unable to read golden file: %v", err)
	}

	// Compare output with golden file content
	if string(golden) != output {
		t.Errorf("output does not match golden file:\nExpected: %s\nGot: %s", golden, output)
	}
}
