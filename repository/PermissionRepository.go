package repository

import "user_management_service/models"

type PermissionRepository interface {
	GetUserPermissions(userID int) ([]models.Permission, error)
}
