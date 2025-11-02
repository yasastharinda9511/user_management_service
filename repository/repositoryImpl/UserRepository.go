package repositoryImpl

import (
	"database/sql"
	"fmt"
	"user_management_service/models"
	"user_management_service/repository"
)

type userRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository instance
func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &userRepository{db: db}
}

// Create creates a new user
func (r *userRepository) Create(user *models.User) error {
	query := `
		INSERT INTO userManagement.users (username, email, password_hash, first_name, last_name, phone)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`

	fmt.Printf("\"New user ID: %s", user.Username)
	var id int
	err := r.db.QueryRow(query,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.FirstName,
		user.LastName,
		user.Phone,
	).Scan(&id)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	fmt.Printf("New user ID: %d\n", id)
	user.ID = id
	return nil
}

// GetByID retrieves a user by ID
func (r *userRepository) GetByID(id int) (*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, first_name, last_name,
		       phone, is_active, is_email_verified, created_at, updated_at, last_login
		FROM userManagement.users WHERE id = $1`

	user := &models.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.FirstName, &user.LastName, &user.Phone, &user.IsActive,
		&user.IsEmailVerified, &user.CreatedAt, &user.UpdatedAt, &user.LastLogin,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetByUsername retrieves a user by username
func (r *userRepository) GetByUsername(username string) (*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, first_name, last_name,
		       phone, is_active, is_email_verified, created_at, updated_at, last_login
		FROM userManagement.users WHERE username = $1`

	user := &models.User{}
	err := r.db.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.FirstName, &user.LastName, &user.Phone, &user.IsActive,
		&user.IsEmailVerified, &user.CreatedAt, &user.UpdatedAt, &user.LastLogin,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// GetByEmail retrieves a user by email
func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, first_name, last_name,
		       phone, is_active, is_email_verified, created_at, updated_at, last_login
		FROM userManagement.users WHERE email = $1`

	user := &models.User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.FirstName, &user.LastName, &user.Phone, &user.IsActive,
		&user.IsEmailVerified, &user.CreatedAt, &user.UpdatedAt, &user.LastLogin,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// Update updates a user
//func (r *userRepository) Update(user *models.User) error {
//	query := `
//		UPDATE users
//		SET first_name = ?, last_name = ?, phone = ?, email = ?, updated_at = NOW()
//		WHERE id = ?`
//
//	result, err := r.db.Exec(query, user.FirstName, user.LastName, user.Phone, user.Email, user.ID)
//	if err != nil {
//		return fmt.Errorf("failed to update user: %w", err)
//	}
//
//	rowsAffected, err := result.RowsAffected()
//	if err != nil {
//		return fmt.Errorf("failed to check affected rows: %w", err)
//	}
//
//	if rowsAffected == 0 {
//		return fmt.Errorf("user not found")
//	}
//
//	return nil
//}

// Delete deletes a user (soft delete by setting is_active to false)
//func (r *userRepository) Delete(id int) error {
//	return r.Deactivate(id)
//}

// List retrieves users with pagination
//
//	func (r *userRepository) List(offset, limit int) ([]models.User, error) {
//		query := `
//			SELECT id, username, email, first_name, last_name, phone,
//			       is_active, is_email_verified, created_at, updated_at, last_login
//			FROM users
//			WHERE is_active = true
//			ORDER BY created_at DESC
//			LIMIT ? OFFSET ?`
//
//		rows, err := r.db.Query(query, limit, offset)
//		if err != nil {
//			return nil, fmt.Errorf("failed to list users: %w", err)
//		}
//		defer rows.Close()
//
//		var users []models.User
//		for rows.Next() {
//			var user models.User
//			err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.FirstName,
//				&user.LastName, &user.Phone, &user.IsActive, &user.IsEmailVerified,
//				&user.CreatedAt, &user.UpdatedAt, &user.LastLogin)
//			if err != nil {
//				return nil, fmt.Errorf("failed to scan user: %w", err)
//			}
//			users = append(users, user)
//		}
//
//		return users, nil
//	}
//
// // Count returns the total number of active users
//
//	func (r *userRepository) Count() (int, error) {
//		var count int
//		query := `SELECT COUNT(*) FROM users WHERE is_active = true`
//		err := r.db.QueryRow(query).Scan(&count)
//		if err != nil {
//			return 0, fmt.Errorf("failed to count users: %w", err)
//		}
//		return count, nil
//	}
//
// UpdateLastLogin updates the last login timestamp for a user

func (r *userRepository) UpdateLastLogin(userID int) error {
	query := `UPDATE userManagement.users SET last_login = NOW() WHERE id = $1`
	_, err := r.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}
	return nil
}

// Deactivate deactivates a user account
func (r *userRepository) Deactivate(userID int) error {
	query := `UPDATE userManagement.users SET is_active = FALSE, updated_at = NOW() WHERE id = $1`
	result, err := r.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to deactivate user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}
