package main

import (
	"fmt"
)

// BigStruct is a large struct that holds some data.
type BigStruct struct {
	ID    int
	Name  string
	Data  [1024]byte // Assume this is a large data field.
	Extra *ExtraData
}

// ExtraData represents additional data.
type ExtraData struct {
	Details string
}

func main() {
	// Create a slice of BigStruct with sensitive data.
	bigStructs := make([]BigStruct, 3)

	// Initialize the slice with some data.
	for i := range bigStructs {
		bigStructs[i] = BigStruct{
			ID:    i,
			Name:  fmt.Sprintf("Struct%d", i),
			Data:  [1024]byte{1, 2, 3, 4}, // Example data initialization.
			Extra: &ExtraData{Details: "Sensitive Details"},
		}
	}

	// Use the data for some operations.
	fmt.Printf("Using BigStruct slice: %+v\n", bigStructs)

	// When the data is no longer needed, zero out the slice.
	zeroBigStructs(bigStructs)

	// Set the slice to nil to avoid retaining the reference.
	bigStructs = nil

	// Continue with other operations...
}

// zeroBigStructs zeros out the fields of each BigStruct in the slice.
func zeroBigStructs(bs []BigStruct) {
	for i := range bs {
		bs[i] = BigStruct{
			ID:    0,
			Name:  "",
			Data:  [1024]byte{}, // Zero out the byte array.
			Extra: nil,          // Set pointer to nil.
		}
	}
}
