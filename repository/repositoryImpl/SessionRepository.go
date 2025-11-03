package repositoryImpl

import (
	"database/sql"
	"fmt"
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
        INSERT INTO userManagement.user_sessions (
            user_id, access_token_hash, access_token_expires_at,
            refresh_token_hash, refresh_token_expires_at,
            created_at, is_revoked
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id`

	var sessionID int64
	err := r.db.QueryRow(query,
		session.UserID,
		session.AccessTokenHash,
		session.AccessTokenExpiresAt,
		session.RefreshTokenHash,
		session.RefreshTokenExpiresAt,
		time.Now(),
		session.IsRevoked,
	).Scan(&sessionID)

	return sessionID, err
}

func (r *SessionRepository) CleanupExpired(userID int) error {
	query := `
        DELETE FROM userManagement.user_sessions
        WHERE user_id = $1
        AND refresh_token_expires_at < NOW()`
	_, err := r.db.Exec(query, userID)
	return err
}

// GetByTokenHash retrieves session by access token hash
func (r *SessionRepository) GetByTokenHash(tokenHash string) (*models.Session, error) {
	query := `
        SELECT id, user_id, access_token_hash, access_token_expires_at,
               refresh_token_hash, refresh_token_expires_at,
               created_at, last_refreshed_at, is_revoked
        FROM userManagement.user_sessions
        WHERE access_token_hash = $1 AND is_revoked = false AND access_token_expires_at > $2
    `

	var session models.Session
	err := r.db.QueryRow(query, tokenHash, time.Now()).Scan(
		&session.ID,
		&session.UserID,
		&session.AccessTokenHash,
		&session.AccessTokenExpiresAt,
		&session.RefreshTokenHash,
		&session.RefreshTokenExpiresAt,
		&session.CreatedAt,
		&session.LastRefreshedAt,
		&session.IsRevoked,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("session not found")
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	return &session, nil
}

func (r *SessionRepository) RevokeSession(sessionID int) error {
	query := `
        UPDATE userManagement.user_sessions
        SET is_revoked = true
        WHERE id = $1
    `

	result, err := r.db.Exec(query, sessionID)
	if err != nil {
		return fmt.Errorf("failed to revoke session: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("session not found")
	}

	return nil
}

func (r *SessionRepository) RevokeAllUserSessions(userID int) error {
	query := `
        UPDATE userManagement.user_sessions
        SET is_revoked = true
        WHERE user_id = $1 AND is_revoked = false
    `

	_, err := r.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to revoke all user sessions: %w", err)
	}

	return nil
}

func (r *SessionRepository) IsSessionValid(tokenHash string) bool {
	query := `
        SELECT COUNT(*)
        FROM userManagement.user_sessions
        WHERE access_token_hash = $1 AND is_revoked = false AND access_token_expires_at > $2
    `

	var count int
	err := r.db.QueryRow(query, tokenHash, time.Now()).Scan(&count)
	if err != nil {
		return false
	}

	return count > 0
}

// GetByRefreshTokenHash retrieves session by refresh token hash
func (r *SessionRepository) GetByRefreshTokenHash(tokenHash string) (*models.Session, error) {
	query := `
        SELECT id, user_id, access_token_hash, access_token_expires_at,
               refresh_token_hash, refresh_token_expires_at,
               created_at, last_refreshed_at, is_revoked
        FROM userManagement.user_sessions
        WHERE refresh_token_hash = $1 AND is_revoked = false AND refresh_token_expires_at > $2
    `

	var session models.Session
	err := r.db.QueryRow(query, tokenHash, time.Now()).Scan(
		&session.ID,
		&session.UserID,
		&session.AccessTokenHash,
		&session.AccessTokenExpiresAt,
		&session.RefreshTokenHash,
		&session.RefreshTokenExpiresAt,
		&session.CreatedAt,
		&session.LastRefreshedAt,
		&session.IsRevoked,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("session not found or refresh token expired")
		}
		return nil, fmt.Errorf("failed to get session by refresh token: %w", err)
	}

	return &session, nil
}

// UpdateAccessToken updates the access token and expiration for a session (used during refresh)
func (r *SessionRepository) UpdateAccessToken(sessionID int, accessTokenHash string, expiresAt time.Time) error {
	query := `
        UPDATE userManagement.user_sessions
        SET access_token_hash = $1, access_token_expires_at = $2, last_refreshed_at = $3
        WHERE id = $4
    `

	result, err := r.db.Exec(query, accessTokenHash, expiresAt, time.Now(), sessionID)
	if err != nil {
		return fmt.Errorf("failed to update access token: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("session not found")
	}

	return nil
}
