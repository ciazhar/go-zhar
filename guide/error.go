package main

import (
	"errors"
	"fmt"
)

// package foo

func Open() error {
	return errors.New("could not open")
}

// package bar

func main() {
	if err := Open(); err != nil {

		if errors.Is(err, errors.New("could not open")) {
			// Handle the error.
			fmt.Println("error ini")
		}

		// Can't handle the error.
		panic("unknown error")
	}

	sameArr := []int{1, 2, 3} // Use ... instead of 3

	a := 111_111111_111.0

}
