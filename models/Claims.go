package models

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	UserID   int
	Username string
	Email    string
	jwt.RegisteredClaims
}
