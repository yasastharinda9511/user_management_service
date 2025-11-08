package request

type CreateRoleRequestDTO struct {
	RoleName      string `json:"role_name"`
	Description   string `json:"description"`
	PermissionIDs []int  `json:"permission_ids"`
}
