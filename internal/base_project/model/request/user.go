package request

type CreateUserBodyRequest struct {
	Name string `json:"name"`
	Age  int    `json:"age" validate:"gt=0"`
}
