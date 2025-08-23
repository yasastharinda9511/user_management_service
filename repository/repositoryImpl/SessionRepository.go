package repositoryImpl

import (
	"database/sql"
	"time"
	"user_management_service/repository"

	"user_management_service/models"
)

type SessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) repository.SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) Create(session *models.Session) (int64, error) {

	query := `
        INSERT INTO user_sessions (user_id, token_hash, expires_at, created_at, is_revoked)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id`

	var sessionID int64
	err := r.db.QueryRow(query,
		session.UserID,
		session.TokenHash,
		session.ExpiresAt,
		time.Now(),
		session.IsRevoked,
	).Scan(&sessionID)

	return sessionID, err
}

func (r *SessionRepository) CleanupExpired(userID int) error {
	query := `DELETE FROM user_sessions WHERE user_id = $1 AND expires_at < NOW()`
	_, err := r.db.Exec(query, userID)
	return err
}
