package repository

import "user_management_service/models"

type SessionRepository interface {
	Create(session *models.Session) (int64, error)
	GetByTokenHash(tokenHash string) (*models.Session, error)
	//GetByUserID(userID int) ([]models.Session, error)
	RevokeSession(sessionID int) error
	RevokeAllUserSessions(userID int) error
	CleanupExpired(userID int) error
	IsSessionValid(tokenHash string) bool
	//IsTokenRevoked(tokenHash string) (bool, error)
}
