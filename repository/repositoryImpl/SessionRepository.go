package repositoryImpl

import (
	"database/sql"
	"fmt"

	"user_management_service/models"
)

type sessionRepository struct {
	db *sql.DB
}

// NewSessionRepository creates a new session repository instance
func NewSessionRepository(db *sql.DB) *sessionRepository {
	return &sessionRepository{db: db}
}

// Create creates a new session
func (r *sessionRepository) Create(session *models.Session) error {
	query := `
		INSERT INTO user_sessions (user_id, token_hash, expires_at, user_agent, ip_address) 
		VALUES (?, ?, ?, ?, ?)`

	result, err := r.db.Exec(query, session.UserID, session.TokenHash, session.ExpiresAt,
		session.UserAgent, session.IPAddress)
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get session ID: %w", err)
	}

	session.ID = int(id)
	return nil
}

// GetByTokenHash retrieves a session by token hash
func (r *sessionRepository) GetByTokenHash(tokenHash string) (*models.Session, error) {
	query := `
		SELECT id, user_id, token_hash, expires_at, created_at, is_revoked, user_agent, ip_address
		FROM user_sessions 
		WHERE token_hash = ? AND expires_at > NOW()`

	session := &models.Session{}
	err := r.db.QueryRow(query, tokenHash).Scan(
		&session.ID, &session.UserID, &session.TokenHash, &session.ExpiresAt,
		&session.CreatedAt, &session.IsRevoked, &session.UserAgent, &session.IPAddress,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("session not found or expired")
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	return session, nil
}

// GetByUserID retrieves all sessions for a user
func (r *sessionRepository) GetByUserID(userID int) ([]models.Session, error) {
	query := `
		SELECT id, user_id, token_hash, expires_at, created_at, is_revoked, user_agent, ip_address
		FROM user_sessions 
		WHERE user_id = ? AND expires_at > NOW()
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get sessions: %w", err)
	}
	defer rows.Close()

	var sessions []models.Session
	for rows.Next() {
		var session models.Session
		err := rows.Scan(&session.ID, &session.UserID, &session.TokenHash,
			&session.ExpiresAt, &session.CreatedAt, &session.IsRevoked,
			&session.UserAgent, &session.IPAddress)
		if err != nil {
			return nil, fmt.Errorf("failed to scan session: %w", err)
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
}

// Revoke revokes a specific session by token hash
func (r *sessionRepository) Revoke(tokenHash string) error {
	query := `UPDATE user_sessions SET is_revoked = TRUE WHERE token_hash = ?`
	result, err := r.db.Exec(query, tokenHash)
	if err != nil {
		return fmt.Errorf("failed to revoke session: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("session not found")
	}

	return nil
}

//
