package types

type Student struct {
	Id    int    `json:"id" validate:"required"`
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
	Age   int    `json:"age" validate:"required,min=0"`
	Phone string `json:"phone" validate:"required"`
}

type response struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}
