package repository

import "user_management_service/models"

type RoleRepository interface {
	GetByID(id int) (*models.Role, error)
	GetByName(name string) (*models.Role, error)
	List() ([]models.Role, error)
	GetUserRoles(userID int) ([]models.Role, error)
	AssignRoleToUser(userID, roleID int) error
	RemoveRoleFromUser(userID, roleID int) error
}
