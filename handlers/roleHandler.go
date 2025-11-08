package handlers

import (
	"encoding/json"
	"net/http"
	"user_management_service/services"
)

type RoleHandler struct {
	roleService services.RoleService
}

func NewRoleHandler(roleService services.RoleService) *RoleHandler {
	return &RoleHandler{roleService}
}

func (h *RoleHandler) GetAllRoles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	roles, err := h.roleService.GetAllRoles()

	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Roles retrieved successfully",
		"roles":   roles,
		"count":   len(roles),
	})
}
