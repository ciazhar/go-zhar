package main

import (
	"fmt"
	"unsafe"
)

type Data struct {
	value int
}

func main() {
	// Creating nested pointers
	data := &Data{value: 42}
	ptr1 := &data
	ptr2 := &ptr1
	ptr3 := &ptr2

	// Accessing value using nested pointers
	fmt.Printf("Value via nested pointers: %d\n", (***ptr3).value)

	// Directly assigning the last pointer to reduce nesting
	simplifiedPtr := **ptr3

	// Accessing value using simplified pointer
	fmt.Printf("Value via simplified pointer: %d\n", (*simplifiedPtr).value)

	// Showing memory address to confirm they point to the same object
	fmt.Printf("Address of original data: %v\n", unsafe.Pointer(data))
	fmt.Printf("Address via nested pointers: %v\n", unsafe.Pointer(***ptr3))
	fmt.Printf("Address via simplified pointer: %v\n", unsafe.Pointer(simplifiedPtr))
}
