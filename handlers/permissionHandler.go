package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"user_management_service/dto/request"
	"user_management_service/services"

	"github.com/gorilla/mux"
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

func (h *PermissionHandler) CreatePermission(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse request body
	var req request.CreatePermissionRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid JSON format"}`, http.StatusBadRequest)
		return
	}

	// Validate request
	if req.Name == "" || req.Resource == "" || req.Action == "" {
		http.Error(w, `{"error": "name, resource, and action are required"}`, http.StatusBadRequest)
		return
	}

	// Call service
	permission, err := h.permissionService.CreatePermission(&req)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":    "Permission created successfully",
		"permission": permission,
	})
}

func (h *PermissionHandler) UpdatePermission(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get permission ID from URL
	vars := mux.Vars(r)
	permissionID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, `{"error": "Invalid permission ID"}`, http.StatusBadRequest)
		return
	}

	// Parse request body
	var req request.UpdatePermissionRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid JSON format"}`, http.StatusBadRequest)
		return
	}

	// Validate request
	if req.Name == "" {
		http.Error(w, `{"error": "name is required"}`, http.StatusBadRequest)
		return
	}

	// Call service
	permission, err := h.permissionService.UpdatePermission(permissionID, &req)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":    "Permission updated successfully",
		"permission": permission,
	})
}
