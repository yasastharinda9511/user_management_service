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

func (s *UserService) CreateUser(req *request.CreateUserRequestDTO) (*models.User, error) {

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

func (s *UserService) GetUserByUsername(username string) (*models.User, error) {
	user, err := s.userRepo.GetByUsername(username)

	if err != nil {
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}

	return user, nil
}

func (s *UserService) GetUserByID(id int) (*models.User, error) {
	user, err := s.userRepo.GetByID(id)

	if err != nil {
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}

	return user, nil
}

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	user, err := s.userRepo.GetByEmail(email)

	if err != nil {
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}

	return user, nil
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	users, err := s.userRepo.GetAll()

	if err != nil {
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}

	// Clear password hashes before returning
	for i := range users {
		users[i].PasswordHash = ""
	}

	return users, nil
}

func (s *UserService) Deactivate(userID int) error {
	err := s.userRepo.Deactivate(userID)

	if err != nil {
		return fmt.Errorf("failed to get user by username: %w", err)
	}

	return nil
}

func (s *UserService) UpdateUser(userID int, req *request.UpdateUserRequest) (*models.User, error) {
	// Validate input
	if req.FirstName == "" || req.LastName == "" || req.Email == "" {
		return nil, fmt.Errorf("first_name, last_name, and email are required")
	}

	// Check if user exists and get current data
	existingUser, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Determine is_active value (use provided value or keep existing)
	isActive := existingUser.IsActive
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	// Update user information
	user, err := s.userRepo.Update(userID, req.FirstName, req.LastName, req.Phone, req.Email, isActive)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Handle role assignment if role_id is provided
	if req.RoleID != nil {
		// Remove all existing roles
		err = s.userRepo.RemoveAllRolesFromUser(userID)
		if err != nil {
			return nil, fmt.Errorf("failed to remove existing roles: %w", err)
		}

		// Assign new role
		err = s.userRepo.AssignRoleToUser(userID, *req.RoleID)
		if err != nil {
			return nil, fmt.Errorf("failed to assign role to user: %w", err)
		}
	}

	// Clear password hash before returning
	user.PasswordHash = ""
	return user, nil
}

func (s *UserService) ToggleUserStatus(userID int) (bool, error) {
	newStatus, err := s.userRepo.ToggleStatus(userID)

	if err != nil {
		return false, fmt.Errorf("failed to toggle user status: %w", err)
	}

	return newStatus, nil
}
