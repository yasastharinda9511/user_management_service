package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"user_management_service/utils"
)

type Config struct {
	Port           string
	DatabaseURL    string
	JWTSecret      string
	TokenDuration  int // in hours
	BCryptCost     int
	Environment    string
	AllowedOrigins []string
}

func Load() (*Config, error) {
	cfg := &Config{
		Port:           getEnv("PORT", "8080"),
		JWTSecret:      getEnv("JWT_SECRET", utils.GenerateSecureJWTSecret()),
		TokenDuration:  getEnvAsInt("TOKEN_DURATION", 24),
		BCryptCost:     getEnvAsInt("BCRYPT_COST", 12),
		Environment:    getEnv("ENVIRONMENT", "development"),
		AllowedOrigins: getEnvAsSlice("ALLOWED_ORIGINS", []string{"*"}),
	}

	// Build database URL
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "3306")
	dbUser := getEnv("DB_USER", "root")
	dbPassword := getEnv("DB_PASSWORD", "password")
	dbName := getEnv("DB_NAME", "user_management")

	cfg.DatabaseURL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	// Validate required fields
	fmt.Println(cfg.JWTSecret)
	if cfg.JWTSecret == "your-super-secret-jwt-key-change-this-in-production" && cfg.Environment == "production" {
		return nil, fmt.Errorf("JWT_SECRET must be set in production environment")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return fallback
}

func getEnvAsSlice(key string, fallback []string) []string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	// Handle single wildcard
	if value == "*" {
		return []string{"*"}
	}

	// Split by comma and clean up
	origins := strings.Split(value, ",")
	var result []string
	for _, origin := range origins {
		trimmed := strings.TrimSpace(origin)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	if len(result) == 0 {
		return fallback
	}

	return result
}
