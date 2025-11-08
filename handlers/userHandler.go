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

	users, err := h.userService.GetAllUsers()

	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Users retrieved successfully",
		"users":   users,
		"count":   len(users),
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
