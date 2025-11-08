package services

import (
	"user_management_service/dto/request"
	"user_management_service/models"
)

type UserService interface {
	CreateUser(req *request.CreateUserRequestDTO) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	GetUserByID(id int) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetAllUsers() ([]models.User, error)
	Deactivate(userID int) error
	ToggleUserStatus(userID int) (bool, error)
}
