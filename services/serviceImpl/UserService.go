package serviceImpl

import (
	"fmt"
	"time"
	"user_management_service/dto/request"
	"user_management_service/models"
	"user_management_service/repository"
	"user_management_service/services"
	"user_management_service/utils"
)

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) services.UserService {
	return &UserService{userRepo}
}

func (s *UserService) CreateUser(req *request.CreateUserRequest) (*models.User, error) {

	hashedPassword, err := utils.HashPassword(req.Password, 12)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.User{
		Email:           req.Email,
		PasswordHash:    hashedPassword,
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		Phone:           req.Phone,
		Username:        req.Username,
		IsActive:        true,
		IsEmailVerified: false,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		LastLogin:       nil,
	}

	err = s.userRepo.Create(user)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Clear password hash before returning
	user.PasswordHash = ""
	return user, nil
}
