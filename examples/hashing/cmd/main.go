package main

import (
	"fmt"

	bcrypt "github.com/ciazhar/go-start-small/pkg/hashing/bcrypt"
)

func main() {
	// Password to hash
	password := "mysecretpassword"

	// Hash the password
	hashedPassword, err := bcrypt.HashPassword(password)
	if err != nil {
		fmt.Println("Error hashing password:", err)
		return
	}
	fmt.Println("Hashed password:", hashedPassword)

	// Validate the password
	fmt.Println("Password validation result:", bcrypt.ValidatePassword(password, hashedPassword))
}
