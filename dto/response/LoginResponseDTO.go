package response

import (
	"time"
	models "user_management_service/models"
)

type LoginResponseDTO struct {
	Token     string       `json:"token"`
	User      *models.User `json:"user"`
	ExpiresAt time.Time    `json:"expires_at"`
}
