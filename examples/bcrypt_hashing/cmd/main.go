package main

import (
	"context"

	bcrypt "github.com/ciazhar/go-start-small/pkg/hashing/bcrypt"
	"github.com/ciazhar/go-start-small/pkg/logger"
)

func main() {
	// Password to hash
	password := "mysecretpassword"

	// Hash the password
	hashedPassword, err := bcrypt.HashPassword(password)
	if err != nil {
		logger.LogFatal(context.Background(), err, "could not hash password", map[string]interface{}{"password": password})
		return
	}
	logger.LogInfo(context.Background(), "hashed password", map[string]interface{}{"hashedPassword": hashedPassword})

	// Validate the password
	valid := bcrypt.ValidatePassword(password, hashedPassword)
	if !valid {
		logger.LogFatal(context.Background(), err, "could not validate password", map[string]interface{}{"password": password})
		return
	} else {
		logger.LogInfo(context.Background(), "password is valid", map[string]interface{}{"password": password})
	}
}
