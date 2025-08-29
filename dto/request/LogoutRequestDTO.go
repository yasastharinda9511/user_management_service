package request

type LogoutRequestDTO struct {
	Token string `json:"token" validate:"required"`
}
