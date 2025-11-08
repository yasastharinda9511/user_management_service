package repository

import (
	"user_management_service/dto/response"
	"user_management_service/models"
)

type RoleRepository interface {
	GetAll() ([]response.RoleWithPermissionsDTO, error)
	GetByID(roleID int) (*models.Role, error)
	Create(name, description string) (*models.Role, error)
	Update(roleID int, name, description string) (*models.Role, error)
	AssignPermissionsToRole(roleID int, permissionIDs []int) error
	RemoveAllPermissionsFromRole(roleID int) error
	//GetByName(name string) (*models.Role, error)
	//List() ([]models.Role, error)
	GetUserRoles(userID int) ([]models.Role, error)
	//AssignRoleToUser(userID, roleID int) error
	//RemoveRoleFromUser(userID, roleID int) error
}
