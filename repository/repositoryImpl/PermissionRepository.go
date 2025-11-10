package repositoryImpl

import (
	"database/sql"
	"fmt"
	"user_management_service/models"
	"user_management_service/repository"
)

type PermissionRepository struct {
	db *sql.DB
}

func (p PermissionRepository) GetAll() ([]models.Permission, error) {
	query := `
        SELECT id, name, resource, action, description, created_at
        FROM userManagement.permissions
        ORDER BY resource, action, name`

	rows, err := p.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var permissions []models.Permission
	for rows.Next() {
		var perm models.Permission
		if err := rows.Scan(&perm.ID, &perm.Name, &perm.Resource, &perm.Action, &perm.Description, &perm.CreatedAt); err != nil {
			return nil, err
		}
		permissions = append(permissions, perm)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}

func (p PermissionRepository) GetUserPermissions(userID int) ([]models.Permission, error) {
	query := `
        SELECT DISTINCT p.id, p.name, p.resource, p.action, p.description, p.created_at
        FROM userManagement.user_roles ur
        JOIN userManagement.roles r ON ur.role_id = r.id
        JOIN userManagement.role_permissions rp ON r.id = rp.role_id
        JOIN userManagement.permissions p ON rp.permission_id = p.id
        WHERE ur.user_id = $1
        ORDER BY p.resource, p.action, p.name`

	rows, err := p.db.Query(query, userID)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var permissions []models.Permission
	for rows.Next() {
		var perm models.Permission
		if err := rows.Scan(&perm.ID, &perm.Name, &perm.Resource, &perm.Action, &perm.Description, &perm.CreatedAt); err != nil {
			return nil, err
		}
		permissions = append(permissions, perm)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}

func (p PermissionRepository) GetByRoleID(roleID int) ([]models.Permission, error) {
	query := `
        SELECT p.id, p.name, p.resource, p.action, p.description, p.created_at
        FROM userManagement.role_permissions rp
        JOIN userManagement.permissions p ON rp.permission_id = p.id
        WHERE rp.role_id = $1
        ORDER BY p.resource, p.action, p.name`

	rows, err := p.db.Query(query, roleID)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var permissions []models.Permission
	for rows.Next() {
		var perm models.Permission
		if err := rows.Scan(&perm.ID, &perm.Name, &perm.Resource, &perm.Action, &perm.Description, &perm.CreatedAt); err != nil {
			return nil, err
		}
		permissions = append(permissions, perm)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}

func (p PermissionRepository) Create(name, resource, action, description string) (*models.Permission, error) {
	query := `
        INSERT INTO userManagement.permissions (name, resource, action, description)
        VALUES ($1, $2, $3, $4)
        RETURNING id, name, resource, action, description, created_at`

	var permission models.Permission
	err := p.db.QueryRow(query, name, resource, action, description).Scan(
		&permission.ID, &permission.Name, &permission.Resource, &permission.Action,
		&permission.Description, &permission.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create permission: %w", err)
	}

	return &permission, nil
}

func (p PermissionRepository) Update(permissionID int, name, description string) (*models.Permission, error) {
	query := `
        UPDATE userManagement.permissions
        SET name = $1, description = $2
        WHERE id = $3
        RETURNING id, name, resource, action, description, created_at`

	var permission models.Permission
	err := p.db.QueryRow(query, name, description, permissionID).Scan(
		&permission.ID, &permission.Name, &permission.Resource, &permission.Action,
		&permission.Description, &permission.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("permission not found")
		}
		return nil, fmt.Errorf("failed to update permission: %w", err)
	}

	return &permission, nil
}

func (p PermissionRepository) HasRoleAssociations(permissionID int) (bool, error) {
	query := `
        SELECT EXISTS(
            SELECT 1 FROM userManagement.role_permissions
            WHERE permission_id = $1
        )`

	var exists bool
	err := p.db.QueryRow(query, permissionID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check role associations: %w", err)
	}

	return exists, nil
}

func (p PermissionRepository) Delete(permissionID int) error {
	query := `DELETE FROM userManagement.permissions WHERE id = $1`

	result, err := p.db.Exec(query, permissionID)
	if err != nil {
		return fmt.Errorf("failed to delete permission: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("permission not found")
	}

	return nil
}

func NewPermissionRepository(db *sql.DB) repository.PermissionRepository {
	return &PermissionRepository{db: db}
}
