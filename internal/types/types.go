package types

type Student struct {
	Id    int    `json:"id" validate:"omitempty"`
	Name  string `json:"name" validate:"required,min=2,max=100"`
	Email string `json:"email" validate:"required,email"`
	Age   int    `json:"age" validate:"required,min=0,max=150"`
}