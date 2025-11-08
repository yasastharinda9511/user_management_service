package services

import "user_management_service/models"

type PermissionService interface {
	GetAllPermissions() ([]models.Permission, error)
}
