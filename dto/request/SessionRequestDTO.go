package request

import "time"

type CreateSessionRequest struct {
	UserID    int
	TokenHash string
	ExpiresAt time.Time
	UserAgent string
	IPAddress string
}
