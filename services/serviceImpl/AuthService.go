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
	userRepo       repository.UserRepository
	sessionRepo    repository.SessionRepository
	rolesRepo      repository.RoleRepository
	permissionRepo repository.PermissionRepository
	jwtSecret      string
	tokenDuration  int
	bcryptCost     int
}

func NewAuthService(userRepo repository.UserRepository, sessionRepo repository.SessionRepository, userRolesRepo repository.RoleRepository, permissionRepo repository.PermissionRepository, jwtSecret string, tokenDuration int, bcryptCost int) services.AuthService {
	return &AuthService{userRepo: userRepo, sessionRepo: sessionRepo, rolesRepo: userRolesRepo, permissionRepo: permissionRepo, jwtSecret: jwtSecret, tokenDuration: tokenDuration, bcryptCost: bcryptCost}
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

	tokenString, expiresAt, err := a.generateJWT(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	session := &models.Session{
		UserID:    user.ID,
		TokenHash: utils.HashSHA256(tokenString),
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
		IsRevoked: false,
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

	roles, err := a.rolesRepo.GetUserRoles(user.ID)

	if err != nil {
		return nil, fmt.Errorf("failed to get roles for user %d: %w", user.ID, err)
	}

	permissions, err := a.permissionRepo.GetUserPermissions(user.ID)

	if err != nil {
		return nil, fmt.Errorf("failed to get permissions for user %d: %w", user.ID, err)
	}

	loginResponse := response.LoginResponseDTO{
		Token:       tokenString,
		User:        user,
		ExpiresAt:   expiresAt,
		SessionID:   sessionID,
		Roles:       roles,
		Permissions: permissions,
	}

	return &loginResponse, nil

}

func (a AuthService) generateJWT(user *models.User) (string, time.Time, error) {

	expirationTime := time.Now().Add(time.Duration(a.tokenDuration) * time.Hour)

	claims := &models.Claims{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
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
