package request

type UpdateUserRequest struct {
	FirstName string `json:"first_name" validate:"required,min=1,max=50"`
	LastName  string `json:"last_name" validate:"required,min=1,max=50"`
	Phone     string `json:"phone" validate:"max=20"`
	Email     string `json:"email" validate:"required,email"`
}
