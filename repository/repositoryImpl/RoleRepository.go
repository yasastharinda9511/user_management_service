package repositoryImpl

import (
	"database/sql"
	"user_management_service/models"
	"user_management_service/repository"
)

type RolesRepository struct {
	db *sql.DB
}

func (r RolesRepository) GetUserRoles(userID int) ([]models.Role, error) {

	query := `
    SELECT r.id, r.name, r.description, r.created_at
    FROM user_roles ur
    JOIN roles r ON ur.role_id = r.id
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

func NewRoleRepository(db *sql.DB) repository.RoleRepository {
	return &RolesRepository{db: db}
}
