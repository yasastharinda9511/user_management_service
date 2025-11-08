package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"user_management_service/dto/request"
	"user_management_service/services"

	"github.com/gorilla/mux"
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

func (h *RoleHandler) CreateRole(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Parse request body
	var req request.CreateRoleRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid JSON format"}`, http.StatusBadRequest)
		return
	}

	// Validate request
	if req.RoleName == "" {
		http.Error(w, `{"error": "role_name is required"}`, http.StatusBadRequest)
		return
	}

	// Call service
	role, err := h.roleService.CreateRole(&req)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Role created successfully",
		"role":    role,
	})
}

func (h *RoleHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get role ID from URL
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error": "Invalid ID"}`, http.StatusBadRequest)
		return
	}

	// Parse request body
	var req request.UpdateRoleRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid JSON format"}`, http.StatusBadRequest)
		return
	}

	// Validate request
	if req.RoleName == "" {
		http.Error(w, `{"error": "role_name is required"}`, http.StatusBadRequest)
		return
	}

	// Call service
	role, err := h.roleService.UpdateRole(id, &req)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Role updated successfully",
		"role":    role,
	})
}
