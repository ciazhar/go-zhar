package model

// User structure
type User struct {
	ID       string
	Username string
	Password string // This should be a hashed password
}
