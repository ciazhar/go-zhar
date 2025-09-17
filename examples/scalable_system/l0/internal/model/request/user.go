package request

type CreateUserBodyRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	FullName string `json:"full_name" validate:"required"`
}

type UpdateUserBodyRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	FullName string `json:"full_name" validate:"required"`
}

type UpsertUserBodyRequest struct {
	Id       string `json:"id" validate:"uuid"`
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	FullName string `json:"full_name" validate:"required"`
}

type UserPathParam struct {
	ID string `params:"id" validate:"required,uuid"`
}

type GetUsersQueryParam struct {
	Name  string `query:"name"` // optional
	Age   int    `query:"age"`  // optional
	Page  int    `query:"page" validate:"min=1"`
	Size  int    `query:"size" validate:"min=1,max=100"`
	Sort  string `query:"sort" validate:"oneof=name age"`
	Order string `query:"order" validate:"oneof=asc desc"`
}

type UserEmailQueryParam struct {
	Email string `query:"email" validate:"required,email"`
}
