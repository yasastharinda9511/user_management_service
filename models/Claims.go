package models

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	UserID      int          `json:"user_id"`
	Username    string       `json:"username"`
	Email       string       `json:"email"`
	TokenType   string       `json:"token_type"` // "access" or "refresh"
	Roles       []Role       `json:"roles"`
	Permissions []Permission `json:"permissions"`
	jwt.RegisteredClaims
}
