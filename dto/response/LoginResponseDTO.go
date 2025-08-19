package response

import (
	models "user_management_service/models"
)

type LoginResponse struct {
	Token     string      `json:"token"`
	User      models.User `json:"user"`
	ExpiresAt int64       `json:"expires_at"`
}
