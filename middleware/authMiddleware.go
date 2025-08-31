package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"user_management_service/models"
	"user_management_service/services"

	"github.com/gorilla/mux"
)

type contextKey string

const (
	UserIDKey      contextKey = "user_id"
	UsernameKey    contextKey = "username"
	EmailKey       contextKey = "email"
	RolesKey       contextKey = "roles"
	PermissionsKey contextKey = "permissions"
	SessionIDKey   contextKey = "session_id"
)

type AuthMiddleware struct {
	auth services.AuthService
}

func NewAuthMiddleware(auth services.AuthService) *AuthMiddleware {
	return &AuthMiddleware{auth: auth}
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			m.unauthorizedResponse(w, "Authorization header is required")
			return
		}

		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			m.unauthorizedResponse(w, "Authorization header must start with 'Bearer '")
			return
		}

		token := strings.TrimPrefix(authHeader, bearerPrefix)
		if token == "" {
			m.unauthorizedResponse(w, "Token is required")
			return
		}

		introspectResponse, err := m.auth.Introspect(token)
		if err != nil {
			m.internalErrorResponse(w, err.Error())
			return
		}

		if !introspectResponse.Active {
			m.unauthorizedResponse(w, "Token validation failed")
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, UserIDKey, introspectResponse.User.ID)
		ctx = context.WithValue(ctx, UsernameKey, introspectResponse.User.Username)
		ctx = context.WithValue(ctx, EmailKey, introspectResponse.User.Email)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *AuthMiddleware) RequirePermission(permission string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// First run authentication
			m.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Get permissions from context
				permissions, ok := r.Context().Value(PermissionsKey).([]models.Permission)
				if !ok {
					m.forbiddenResponse(w, "Unable to verify permissions")
					return
				}

				// Check if user has the required permission
				hasPermission := false
				for _, perm := range permissions {
					if perm.Name == permission {
						hasPermission = true
						break
					}
				}

				if !hasPermission {
					m.forbiddenResponse(w, fmt.Sprintf("Permission '%s' is required", permission))
					return
				}

				// Continue with the request
				next.ServeHTTP(w, r)
			})).ServeHTTP(w, r)
		})
	}
}

// RequireRole middleware to check for specific roles
func (m *AuthMiddleware) RequireRole(role string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// First run authentication
			m.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Get roles from context
				roles, ok := r.Context().Value(RolesKey).([]models.Role)
				if !ok {
					m.forbiddenResponse(w, "Unable to verify roles")
					return
				}

				// Check if user has the required role
				hasRole := false
				for _, userRole := range roles {
					if userRole.Name == role {
						hasRole = true
						break
					}
				}

				if !hasRole {
					m.forbiddenResponse(w, fmt.Sprintf("Role '%s' is required", role))
					return
				}

				// Continue with the request
				next.ServeHTTP(w, r)
			})).ServeHTTP(w, r)
		})
	}
}

// Helper functions to extract user information from context
func GetUserIDFromContext(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value(UserIDKey).(int)
	return userID, ok
}

func GetUsernameFromContext(ctx context.Context) (string, bool) {
	username, ok := ctx.Value(UsernameKey).(string)
	return username, ok
}

func GetEmailFromContext(ctx context.Context) (string, bool) {
	email, ok := ctx.Value(EmailKey).(string)
	return email, ok
}

func GetRolesFromContext(ctx context.Context) ([]models.Role, bool) {
	roles, ok := ctx.Value(RolesKey).([]models.Role)
	return roles, ok
}

func GetPermissionsFromContext(ctx context.Context) ([]models.Permission, bool) {
	permissions, ok := ctx.Value(PermissionsKey).([]models.Permission)
	return permissions, ok
}

func GetSessionIDFromContext(ctx context.Context) (int, bool) {
	sessionID, ok := ctx.Value(SessionIDKey).(int)
	return sessionID, ok
}

func (m *AuthMiddleware) unauthorizedResponse(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":   "Unauthorized",
		"message": message,
		"status":  http.StatusUnauthorized,
	})
}

func (m *AuthMiddleware) forbiddenResponse(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":   "Forbidden",
		"message": message,
		"status":  http.StatusForbidden,
	})
}
func (m *AuthMiddleware) internalErrorResponse(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":   "Internal Server Error",
		"message": message,
		"status":  http.StatusInternalServerError,
	})
}
