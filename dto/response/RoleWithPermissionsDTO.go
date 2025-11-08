package response

import (
	"time"
	"user_management_service/models"
)

type RoleWithPermissionsDTO struct {
	ID          int                 `json:"id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	CreatedAt   time.Time           `json:"created_at"`
	Permissions []models.Permission `json:"permissions"`
}
