package services

import "user_management_service/dto/response"

type RoleService interface {
	GetAllRoles() ([]response.RoleWithPermissionsDTO, error)
}
