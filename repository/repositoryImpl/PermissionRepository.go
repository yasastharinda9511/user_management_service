package repositoryImpl

import (
	"database/sql"
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

func NewPermissionRepository(db *sql.DB) repository.PermissionRepository {
	return &PermissionRepository{db: db}
}
