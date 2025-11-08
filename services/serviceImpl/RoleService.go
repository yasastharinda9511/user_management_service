package serviceImpl

import (
	"fmt"
	"user_management_service/dto/request"
	"user_management_service/dto/response"
	"user_management_service/repository"
	"user_management_service/services"
)

type RoleService struct {
	roleRepo       repository.RoleRepository
	permissionRepo repository.PermissionRepository
}

func NewRoleService(roleRepo repository.RoleRepository, permissionRepo repository.PermissionRepository) services.RoleService {
	return &RoleService{
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
	}
}

func (s *RoleService) GetAllRoles() ([]response.RoleWithPermissionsDTO, error) {
	roles, err := s.roleRepo.GetAll()

	if err != nil {
		return nil, fmt.Errorf("failed to get all roles: %w", err)
	}

	return roles, nil
}

func (s *RoleService) CreateRole(req *request.CreateRoleRequestDTO) (*response.RoleWithPermissionsDTO, error) {
	// Validate input
	if req.RoleName == "" {
		return nil, fmt.Errorf("role_name is required")
	}

	// Create role
	role, err := s.roleRepo.Create(req.RoleName, req.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}

	// Assign permissions if provided
	if len(req.PermissionIDs) > 0 {
		err = s.roleRepo.AssignPermissionsToRole(role.ID, req.PermissionIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to assign permissions to role: %w", err)
		}
	}

	// Fetch permissions for the response
	permissions, err := s.permissionRepo.GetByRoleID(role.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get permissions for role: %w", err)
	}

	roleWithPermissions := &response.RoleWithPermissionsDTO{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		CreatedAt:   role.CreatedAt,
		Permissions: permissions,
	}

	return roleWithPermissions, nil
}

func (s *RoleService) UpdateRole(roleID int, req *request.UpdateRoleRequestDTO) (*response.RoleWithPermissionsDTO, error) {
	// Validate input
	if req.RoleName == "" {
		return nil, fmt.Errorf("role_name is required")
	}

	// Check if role exists
	_, err := s.roleRepo.GetByID(roleID)
	if err != nil {
		return nil, fmt.Errorf("role not found: %w", err)
	}

	// Update role
	role, err := s.roleRepo.Update(roleID, req.RoleName, req.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to update role: %w", err)
	}

	// Remove all existing permissions
	err = s.roleRepo.RemoveAllPermissionsFromRole(roleID)
	if err != nil {
		return nil, fmt.Errorf("failed to remove existing permissions: %w", err)
	}

	// Assign new permissions if provided
	if len(req.PermissionIDs) > 0 {
		err = s.roleRepo.AssignPermissionsToRole(roleID, req.PermissionIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to assign permissions to role: %w", err)
		}
	}

	// Fetch permissions for the response
	permissions, err := s.permissionRepo.GetByRoleID(roleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get permissions for role: %w", err)
	}

	roleWithPermissions := &response.RoleWithPermissionsDTO{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		CreatedAt:   role.CreatedAt,
		Permissions: permissions,
	}

	return roleWithPermissions, nil
}
