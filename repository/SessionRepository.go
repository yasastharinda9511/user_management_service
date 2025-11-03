package repository

import (
	"time"
	"user_management_service/models"
)

type SessionRepository interface {
	Create(session *models.Session) (int64, error)
	GetByTokenHash(tokenHash string) (*models.Session, error)
	GetByRefreshTokenHash(tokenHash string) (*models.Session, error)
	UpdateAccessToken(sessionID int, accessTokenHash string, expiresAt time.Time) error
	RevokeSession(sessionID int) error
	RevokeAllUserSessions(userID int) error
	CleanupExpired(userID int) error
	IsSessionValid(tokenHash string) bool
}
