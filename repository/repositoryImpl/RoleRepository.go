package repositoryImpl

import (
	"database/sql"
	"fmt"
	"user_management_service/dto/response"
	"user_management_service/models"
	"user_management_service/repository"
)

type RolesRepository struct {
	db             *sql.DB
	permissionRepo repository.PermissionRepository
}

func (r RolesRepository) GetAll() ([]response.RoleWithPermissionsDTO, error) {
	query := `
		SELECT id, name, description, created_at
		FROM userManagement.roles
		ORDER BY name ASC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rolesWithPermissions []response.RoleWithPermissionsDTO

	for rows.Next() {
		var role models.Role
		if err := rows.Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt); err != nil {
			return nil, err
		}

		// Fetch permissions for this role
		permissions, err := r.permissionRepo.GetByRoleID(role.ID)
		if err != nil {
			return nil, err
		}

		roleWithPerms := response.RoleWithPermissionsDTO{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
			CreatedAt:   role.CreatedAt,
			Permissions: permissions,
		}

		rolesWithPermissions = append(rolesWithPermissions, roleWithPerms)
	}

	// check for iteration errors
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return rolesWithPermissions, nil
}

func (r RolesRepository) GetUserRoles(userID int) ([]models.Role, error) {

	query := `
    SELECT r.id, r.name, r.description, r.created_at
    FROM userManagement.user_roles ur
    JOIN userManagement.roles r ON ur.role_id = r.id
    WHERE ur.user_id = $1`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []models.Role

	for rows.Next() {
		var role models.Role
		if err := rows.Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	// check for iteration errors
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return roles, nil

}

func (r RolesRepository) GetByID(roleID int) (*models.Role, error) {
	query := `
		SELECT id, name, description, created_at
		FROM userManagement.roles
		WHERE id = $1`

	var role models.Role
	err := r.db.QueryRow(query, roleID).Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("role not found")
		}
		return nil, err
	}

	return &role, nil
}

func (r RolesRepository) Create(name, description string) (*models.Role, error) {
	query := `
		INSERT INTO userManagement.roles (name, description)
		VALUES ($1, $2)
		RETURNING id, name, description, created_at`

	var role models.Role
	err := r.db.QueryRow(query, name, description).Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}

	return &role, nil
}

func (r RolesRepository) Update(roleID int, name, description string) (*models.Role, error) {
	query := `
		UPDATE userManagement.roles
		SET name = $1, description = $2
		WHERE id = $3
		RETURNING id, name, description, created_at`

	var role models.Role
	err := r.db.QueryRow(query, name, description, roleID).Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("role not found")
		}
		return nil, fmt.Errorf("failed to update role: %w", err)
	}

	return &role, nil
}

func (r RolesRepository) AssignPermissionsToRole(roleID int, permissionIDs []int) error {
	// Insert permissions for this role
	for _, permissionID := range permissionIDs {
		query := `
			INSERT INTO userManagement.role_permissions (role_id, permission_id)
			VALUES ($1, $2)
			ON CONFLICT (role_id, permission_id) DO NOTHING`

		_, err := r.db.Exec(query, roleID, permissionID)
		if err != nil {
			return fmt.Errorf("failed to assign permission %d to role: %w", permissionID, err)
		}
	}

	return nil
}

func (r RolesRepository) RemoveAllPermissionsFromRole(roleID int) error {
	query := `DELETE FROM userManagement.role_permissions WHERE role_id = $1`

	_, err := r.db.Exec(query, roleID)
	if err != nil {
		return fmt.Errorf("failed to remove permissions from role: %w", err)
	}

	return nil
}

func NewRoleRepository(db *sql.DB, permissionRepo repository.PermissionRepository) repository.RoleRepository {
	return &RolesRepository{
		db:             db,
		permissionRepo: permissionRepo,
	}
}
