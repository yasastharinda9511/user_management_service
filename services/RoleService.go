package services

import (
	"user_management_service/dto/request"
	"user_management_service/dto/response"
)

type RoleService interface {
	GetAllRoles() ([]response.RoleWithPermissionsDTO, error)
	CreateRole(req *request.CreateRoleRequestDTO) (*response.RoleWithPermissionsDTO, error)
	UpdateRole(roleID int, req *request.UpdateRoleRequestDTO) (*response.RoleWithPermissionsDTO, error)
}
