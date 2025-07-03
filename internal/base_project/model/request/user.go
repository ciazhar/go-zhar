package request

type CreateUserBodyRequest struct {
	Name string `json:"name" validate:"required"`
	Age  int    `json:"age" validate:"required,gt=0"`
}

type UpdateUserBodyRequest struct {
	Name string `json:"name" validate:"required"`
	Age  int    `json:"age" validate:"required,gt=0"`
}

type UserPathParam struct {
	ID string `params:"id" validate:"required,uuid"`
}

type GetUsersQueryParam struct {
	Name string `query:"name"` // optional
	Age  int    `query:"age"`  // optional
	Page int    `query:"page" validate:"min=1"`
	Size int    `query:"size" validate:"min=1,max=100"`
}
