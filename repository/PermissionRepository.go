package repository

import "user_management_service/models"

type PermissionRepository interface {
	GetAll() ([]models.Permission, error)
	GetUserPermissions(userID int) ([]models.Permission, error)
	GetByRoleID(roleID int) ([]models.Permission, error)
	Create(name, resource, action, description string) (*models.Permission, error)
}
