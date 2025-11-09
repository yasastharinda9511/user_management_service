package serviceImpl

import (
	"fmt"
	"user_management_service/dto/request"
	"user_management_service/models"
	"user_management_service/repository"
	"user_management_service/services"
)

type PermissionService struct {
	permissionRepo repository.PermissionRepository
}

func NewPermissionService(permissionRepo repository.PermissionRepository) services.PermissionService {
	return &PermissionService{permissionRepo}
}

func (s *PermissionService) GetAllPermissions() ([]models.Permission, error) {
	permissions, err := s.permissionRepo.GetAll()

	if err != nil {
		return nil, fmt.Errorf("failed to get all permissions: %w", err)
	}

	return permissions, nil
}

func (s *PermissionService) CreatePermission(req *request.CreatePermissionRequestDTO) (*models.Permission, error) {
	// Validate input
	if req.Name == "" || req.Resource == "" || req.Action == "" {
		return nil, fmt.Errorf("name, resource, and action are required")
	}

	// Create permission
	permission, err := s.permissionRepo.Create(req.Name, req.Resource, req.Action, req.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to create permission: %w", err)
	}

	return permission, nil
}
