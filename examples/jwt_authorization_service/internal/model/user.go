package model

// User structure
// @Description User structure
type User struct {
	ID       string `json:"id"`
	Username string `json:"username" validate:"required,min=3,max=30"`
	Password string `json:"password" validate:"required,min=6,max=100"`
}
