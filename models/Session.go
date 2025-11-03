package models

import (
	"time"
)

// Session represents a user session with access and refresh tokens
type Session struct {
	ID                    int        `json:"id" db:"id"`
	UserID                int        `json:"user_id" db:"user_id"`
	AccessTokenHash       string     `json:"-" db:"access_token_hash"`
	AccessTokenExpiresAt  time.Time  `json:"access_token_expires_at" db:"access_token_expires_at"`
	RefreshTokenHash      string     `json:"-" db:"refresh_token_hash"`
	RefreshTokenExpiresAt time.Time  `json:"refresh_token_expires_at" db:"refresh_token_expires_at"`
	CreatedAt             time.Time  `json:"created_at" db:"created_at"`
	LastRefreshedAt       *time.Time `json:"last_refreshed_at,omitempty" db:"last_refreshed_at"`
	IsRevoked             bool       `json:"is_revoked" db:"is_revoked"`
}
