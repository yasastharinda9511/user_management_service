package response

import (
	"time"
	models "user_management_service/models"
)

type LoginResponseDTO struct {
	User                  *models.User        `json:"user"`
	SessionID             int64               `json:"session_id"`
	Roles                 []models.Role       `json:"roles"`
	Permissions           []models.Permission `json:"permissions"`
	AccessToken           string              `json:"access_token"`
	AccessTokenExpiresAt  time.Time           `json:"access_token_expires_at"`
	RefreshToken          string              `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time           `json:"refresh_token_expires_at"`
}
