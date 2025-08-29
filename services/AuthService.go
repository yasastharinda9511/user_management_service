package services

import (
	"user_management_service/dto/request"
	"user_management_service/dto/response"
	"user_management_service/models"
)

type AuthService interface {
	Register(req request.CreateUserRequestDTO) (*models.User, error)
	Login(req request.LoginRequestDTO) (*response.LoginResponseDTO, error)
	Logout(req request.LogoutRequestDTO) error
	Introspect(token string) (*response.IntrospectResponse, error)
}
