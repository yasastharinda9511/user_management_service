package response

import (
	"time"
	models "user_management_service/models"
)

type LoginResponseDTO struct {
	User        *models.User        `json:"user"`
	SessionID   int64               `json:"session_id"`
	Roles       []models.Role       `json:"roles"`
	Permissions []models.Permission `json:"permissions"`
	Token       string              `json:"token"`
	ExpiresAt   time.Time           `json:"expires_at"`
}
