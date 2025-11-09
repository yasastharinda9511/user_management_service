package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"user_management_service/dto/request"
	"user_management_service/services"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{userService}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// Parse request body
	var req request.CreateUserRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid JSON format"}`, http.StatusBadRequest)
		return
	}

	// Validate request
	if req.Email == "" || req.Password == "" {
		http.Error(w, `{"error": "Email and password are required"}`, http.StatusBadRequest)
		return
	}

	// Call service
	user, err := h.userService.CreateUser(&req)
	if err != nil {
		// You might want to handle different error types differently
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User created successfully",
		"user":    user,
	})
}

func (h *UserHandler) GetUserByUsername(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)          // get path variables
	username := vars["username"] // extract username

	if username == "" {
		http.Error(w, `{"error": "Username is required"}`, http.StatusBadRequest)
	}

	user, err := h.userService.GetUserByUsername(username)

	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User Retrieved successfully",
		"user":    user,
	})
}

func (h *UserHandler) GetUserByUserID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // get path variables
	idStr := vars["id"] // extract username

	fmt.Printf("user id: %s", idStr)

	id, err := strconv.Atoi(idStr) // convert to int
	if err != nil {
		http.Error(w, `{"error": "Invalid ID"}`, http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetUserByID(id)

	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User retrieved successfully",
		"user":    user,
	})
}

func (h *UserHandler) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)    // get path variables
	email := vars["email"] // extract username

	if email == "" {
		http.Error(w, `{"error": "Username is required"}`, http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetUserByEmail(email)

	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User Retrieved successfully",
		"user":    user,
	})
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	usersWithRoles, err := h.userService.GetAllUsers()

	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Users retrieved successfully",
		"users":   usersWithRoles,
		"count":   len(usersWithRoles),
	})
}

func (h *UserHandler) DeactivateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // get path variables
	idStr := vars["id"] // extract username

	fmt.Printf("user id: %s", idStr)

	id, err := strconv.Atoi(idStr) // convert to int
	if err != nil {
		http.Error(w, `{"error": "Invalid ID"}`, http.StatusBadRequest)
		return
	}

	err = h.userService.Deactivate(id)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User Deactivated successfully",
	})
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get user ID from URL
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error": "Invalid ID"}`, http.StatusBadRequest)
		return
	}

	// Parse request body
	var req request.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid JSON format"}`, http.StatusBadRequest)
		return
	}

	// Validate request
	if req.FirstName == "" || req.LastName == "" || req.Email == "" {
		http.Error(w, `{"error": "first_name, last_name, and email are required"}`, http.StatusBadRequest)
		return
	}

	// Call service
	user, err := h.userService.UpdateUser(id, &req)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User updated successfully",
		"user":    user,
	})
}

func (h *UserHandler) ToggleUserStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r) // get path variables
	idStr := vars["id"] // extract user ID

	id, err := strconv.Atoi(idStr) // convert to int
	if err != nil {
		http.Error(w, `{"error": "Invalid ID"}`, http.StatusBadRequest)
		return
	}

	newStatus, err := h.userService.ToggleUserStatus(id)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	statusText := "deactivated"
	if newStatus {
		statusText = "activated"
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":   "User status toggled successfully",
		"is_active": newStatus,
		"status":    statusText,
	})
}
