package serviceImpl

import (
	"fmt"
	"user_management_service/dto/response"
	"user_management_service/repository"
	"user_management_service/services"
)

type RoleService struct {
	roleRepo repository.RoleRepository
}

func NewRoleService(roleRepo repository.RoleRepository) services.RoleService {
	return &RoleService{roleRepo}
}

func (s *RoleService) GetAllRoles() ([]response.RoleWithPermissionsDTO, error) {
	roles, err := s.roleRepo.GetAll()

	if err != nil {
		return nil, fmt.Errorf("failed to get all roles: %w", err)
	}

	return roles, nil
}
