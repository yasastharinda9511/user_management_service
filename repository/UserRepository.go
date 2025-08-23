package repository

import (
	"user_management_service/models"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Create(user *models.User) error
	GetByID(id int) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	//Update(user *models.User) error
	//Delete(id int) error
	//List(offset, limit int) ([]models.User, error)
	//Count() (int, error)
	UpdateLastLogin(userID int) error
	Deactivate(userID int) error
}
