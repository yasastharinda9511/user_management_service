package services

import (
	"user_management_service/dto/request"
	"user_management_service/models"
)

type PermissionService interface {
	GetAllPermissions() ([]models.Permission, error)
	CreatePermission(req *request.CreatePermissionRequestDTO) (*models.Permission, error)
	UpdatePermission(permissionID int, req *request.UpdatePermissionRequestDTO) (*models.Permission, error)
	DeletePermission(permissionID int) error
}
