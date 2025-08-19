package models

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID              int        `json:"id" db:"id"`
	Username        string     `json:"username" db:"username"`
	Email           string     `json:"email" db:"email"`
	PasswordHash    string     `json:"-" db:"password_hash"`
	FirstName       string     `json:"first_name" db:"first_name"`
	LastName        string     `json:"last_name" db:"last_name"`
	Phone           string     `json:"phone" db:"phone"`
	IsActive        bool       `json:"is_active" db:"is_active"`
	IsEmailVerified bool       `json:"is_email_verified" db:"is_email_verified"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
	LastLogin       *time.Time `json:"last_login" db:"last_login"`
}
