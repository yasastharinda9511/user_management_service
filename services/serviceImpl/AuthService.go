package serviceImpl

import (
	"fmt"
	"strconv"
	"time"
	"user_management_service/dto/request"
	"user_management_service/dto/response"
	"user_management_service/models"
	"user_management_service/repository"
	"user_management_service/services"
	"user_management_service/utils"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	userRepo             repository.UserRepository
	sessionRepo          repository.SessionRepository
	rolesRepo            repository.RoleRepository
	permissionRepo       repository.PermissionRepository
	jwtSecret            string
	accessTokenDuration  int // in minutes
	refreshTokenDuration int // in days
	bcryptCost           int
}

func NewAuthService(userRepo repository.UserRepository, sessionRepo repository.SessionRepository, userRolesRepo repository.RoleRepository, permissionRepo repository.PermissionRepository, jwtSecret string, accessTokenDuration int, refreshTokenDuration int, bcryptCost int) services.AuthService {
	return &AuthService{
		userRepo:             userRepo,
		sessionRepo:          sessionRepo,
		rolesRepo:            userRolesRepo,
		permissionRepo:       permissionRepo,
		jwtSecret:            jwtSecret,
		accessTokenDuration:  accessTokenDuration,
		refreshTokenDuration: refreshTokenDuration,
		bcryptCost:           bcryptCost,
	}
}

func (a AuthService) Register(req request.CreateUserRequestDTO) (*models.User, error) {
	if existingUser, _ := a.userRepo.GetByUsername(req.Username); existingUser != nil {
		return nil, fmt.Errorf("username already exists")
	}

	// Check if email already exists
	if existingUser, _ := a.userRepo.GetByEmail(req.Email); existingUser != nil {
		return nil, fmt.Errorf("email already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password, a.bcryptCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user model
	user := &models.User{
		Username:        req.Username,
		Email:           req.Email,
		PasswordHash:    hashedPassword,
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		Phone:           req.Phone,
		IsActive:        true,
		IsEmailVerified: false,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Save user
	err = a.userRepo.Create(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Clear password hash before returning
	user.PasswordHash = ""
	return user, nil
}

func (a AuthService) Login(req request.LoginRequestDTO) (*response.LoginResponseDTO, error) {

	user, err := a.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, fmt.Errorf("account is deactivated")
	}

	// Verify password
	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Update last login
	if err := a.userRepo.UpdateLastLogin(user.ID); err != nil {
		// Log but don't fail the login
		fmt.Printf("Warning: failed to update last login for user %d: %v\n", user.ID, err)
	}

	// Fetch roles and permissions before generating tokens
	roles, err := a.rolesRepo.GetUserRoles(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get roles for user %d: %w", user.ID, err)
	}

	permissions, err := a.permissionRepo.GetUserPermissions(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get permissions for user %d: %w", user.ID, err)
	}

	// Generate access token with roles and permissions
	accessToken, accessExpiresAt, err := a.generateToken(user, "access", roles, permissions)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token with roles and permissions
	refreshToken, refreshExpiresAt, err := a.generateToken(user, "refresh", roles, permissions)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Create session with both tokens
	session := &models.Session{
		UserID:                user.ID,
		AccessTokenHash:       utils.HashSHA256(accessToken),
		AccessTokenExpiresAt:  accessExpiresAt,
		RefreshTokenHash:      utils.HashSHA256(refreshToken),
		RefreshTokenExpiresAt: refreshExpiresAt,
		CreatedAt:             time.Now(),
		IsRevoked:             false,
	}
	sessionID, err := a.sessionRepo.Create(session)

	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	go func() {
		if err := a.sessionRepo.CleanupExpired(user.ID); err != nil {
			fmt.Printf("Warning: failed to cleanup expired sessions for user %d: %v\n", user.ID, err)
		}
	}()

	loginResponse := response.LoginResponseDTO{
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessExpiresAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshExpiresAt,
		User:                  user,
		SessionID:             sessionID,
		Roles:                 roles,
		Permissions:           permissions,
	}

	return &loginResponse, nil

}

func (a AuthService) Logout(req request.LogoutRequestDTO) error {
	// Extract token from request (usually from Authorization header)
	tokenString := req.Token
	if tokenString == "" {
		return fmt.Errorf("token is required")
	}

	// Parse and validate the JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Make sure the signing method is what you expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(a.jwtSecret), nil // Your JWT secret key
	})

	if err != nil {
		return fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return fmt.Errorf("invalid token claims")
	}

	// Verify it's an access token (not a refresh token)
	tokenType, ok := claims["token_type"].(string)
	if !ok || tokenType != "access" {
		return fmt.Errorf("logout requires access token")
	}

	// Extract user ID from token claims
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return fmt.Errorf("invalid user ID in token")
	}

	intUserID := int(userID)
	// Hash the token to match what's stored in the database
	tokenHash := utils.HashSHA256(tokenString)
	// Find and revoke the session
	session, err := a.sessionRepo.GetByTokenHash(tokenHash)
	if err != nil {
		return fmt.Errorf("session not found")
	}

	// Verify the session belongs to the user
	if session.UserID != intUserID {
		return fmt.Errorf("unauthorized")
	}

	// Revoke the session
	if err := a.sessionRepo.RevokeSession(session.ID); err != nil {
		return fmt.Errorf("failed to revoke session: %w", err)
	}

	return nil
}

// RefreshToken exchanges a valid refresh token for a new access token
func (a AuthService) RefreshToken(refreshToken string) (*response.RefreshTokenResponseDTO, error) {
	// Parse and validate the refresh token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(a.jwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	// Verify it's a refresh token
	tokenType, ok := claims["token_type"].(string)
	if !ok || tokenType != "refresh" {
		return nil, fmt.Errorf("token is not a refresh token")
	}

	// Extract user ID
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid user ID in token")
	}

	// Hash the refresh token to look up in database
	tokenHash := utils.HashSHA256(refreshToken)

	// Find the session by refresh token hash
	session, err := a.sessionRepo.GetByRefreshTokenHash(tokenHash)
	if err != nil {
		return nil, fmt.Errorf("invalid or expired refresh token")
	}

	// Verify the session belongs to the user
	if session.UserID != int(userID) {
		return nil, fmt.Errorf("unauthorized")
	}

	// Get the user
	user, err := a.userRepo.GetByID(int(userID))
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Check if user is still active
	if !user.IsActive {
		return nil, fmt.Errorf("account is deactivated")
	}

	// Fetch roles and permissions for the new access token
	roles, err := a.rolesRepo.GetUserRoles(int(userID))
	if err != nil {
		return nil, fmt.Errorf("failed to get roles for user %d: %w", int(userID), err)
	}

	permissions, err := a.permissionRepo.GetUserPermissions(int(userID))
	if err != nil {
		return nil, fmt.Errorf("failed to get permissions for user %d: %w", int(userID), err)
	}

	// Generate a new access token with roles and permissions
	newAccessToken, newAccessExpiresAt, err := a.generateToken(user, "access", roles, permissions)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new access token: %w", err)
	}

	// Update the session with the new access token
	newAccessTokenHash := utils.HashSHA256(newAccessToken)
	err = a.sessionRepo.UpdateAccessToken(session.ID, newAccessTokenHash, newAccessExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("failed to update session: %w", err)
	}

	// Return the new access token
	refreshResponse := &response.RefreshTokenResponseDTO{
		AccessToken:          newAccessToken,
		AccessTokenExpiresAt: newAccessExpiresAt,
	}

	return refreshResponse, nil
}

// generateToken generates a JWT token (access or refresh)
func (a AuthService) generateToken(user *models.User, tokenType string, roles []models.Role, permissions []models.Permission) (string, time.Time, error) {
	var expirationTime time.Time

	if tokenType == "access" {
		expirationTime = time.Now().Add(time.Duration(a.accessTokenDuration) * time.Minute)
	} else if tokenType == "refresh" {
		expirationTime = time.Now().Add(time.Duration(a.refreshTokenDuration) * 24 * time.Hour)
	} else {
		return "", time.Time{}, fmt.Errorf("invalid token type: %s", tokenType)
	}

	// Format role as string (user has single role)
	var roleName string
	if len(roles) > 0 {
		roleName = roles[0].Name
	}

	// Format permissions as []string in "resource.action" format
	permissionStrings := make([]string, 0, len(permissions))
	for _, perm := range permissions {
		permissionStrings = append(permissionStrings, fmt.Sprintf("%s.%s", perm.Resource, perm.Action))
	}

	claims := &models.Claims{
		UserID:      user.ID,
		Username:    user.Username,
		Email:       user.Email,
		TokenType:   tokenType,
		Role:        roleName,
		Permissions: permissionStrings,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   strconv.Itoa(user.ID),
			Issuer:    "user-management-service",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(a.jwtSecret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expirationTime, nil
}

func (a AuthService) Introspect(tokenString string) (*response.IntrospectResponse, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Make sure the signing method is what you expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(a.jwtSecret), nil // Your JWT secret key
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	// Verify it's an access token (not a refresh token)
	tokenType, ok := claims["token_type"].(string)
	if !ok || tokenType != "access" {
		return nil, fmt.Errorf("token is not an access token")
	}

	userID, ok := claims["user_id"].(float64)

	if !ok {
		return nil, fmt.Errorf("invalid user ID in token")
	}

	tokenHash := utils.HashSHA256(tokenString)

	session, err := a.sessionRepo.GetByTokenHash(tokenHash)
	if err != nil || session == nil {
		return nil, fmt.Errorf("session not found")
	}

	user, err := a.userRepo.GetByID(int(userID))
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	introspectResponse := response.IntrospectResponse{
		Active: true,
		User:   user,
	}

	return &introspectResponse, nil
}
