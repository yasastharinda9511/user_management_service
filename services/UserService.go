package services

import (
	"user_management_service/dto/request"
	"user_management_service/models"
)

type UserService interface {
	CreateUser(req *request.CreateUserRequest) (*models.User, error)
}
