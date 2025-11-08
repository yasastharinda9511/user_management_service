package handlers

import (
	"encoding/json"
	"net/http"
	"user_management_service/services"
)

type PermissionHandler struct {
	permissionService services.PermissionService
}

func NewPermissionHandler(permissionService services.PermissionService) *PermissionHandler {
	return &PermissionHandler{permissionService}
}

func (h *PermissionHandler) GetAllPermissions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	permissions, err := h.permissionService.GetAllPermissions()

	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":     "Permissions retrieved successfully",
		"permissions": permissions,
		"count":       len(permissions),
	})
}
