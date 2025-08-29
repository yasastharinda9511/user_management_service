package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"user_management_service/dto/request"
	"user_management_service/dto/response"
	"user_management_service/services"
)

type AuthHandler struct {
	auth services.AuthService
}

func NewAuthHandler(auth services.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// Parse request body
	var req request.LoginRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid JSON format"}`, http.StatusBadRequest)
		return
	}

	if req.Password == "" || req.Email == "" {
		http.Error(w, `{"error": "Wrong username or Password"}`, http.StatusBadRequest)
		return
	}

	auth, err := h.auth.Login(req)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User created successfully",
		"auth":    auth,
	})

}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var req request.LogoutRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error": "Invalid JSON format"}`, http.StatusBadRequest)
	}

	if req.Token == "" {
		http.Error(w, `{"error": "Wrong token"}`, http.StatusBadRequest)
		return
	}

	err := h.auth.Logout(req)

	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User logged out successfully",
	})

}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
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
	user, err := h.auth.Register(req)
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

func (h *AuthHandler) Introspect(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	authHeader := r.Header.Get("Authorization")

	// Check if Authorization header is present
	if authHeader == "" {
		// Return inactive token response
		res := response.IntrospectResponse{
			Active: false,
			User:   nil,
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)
		return
	}

	// Extract the token (remove "Bearer " prefix)
	token := strings.TrimPrefix(authHeader, "Bearer ")
	token = strings.TrimSpace(token)

	// Check if token is empty after removing Bearer prefix
	if token == "" || token == authHeader {
		// Token doesn't have Bearer prefix or is empty
		res := response.IntrospectResponse{
			Active: false,
			User:   nil,
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)
		return
	}

	// Call your introspection service
	introspectResponse, err := h.auth.Introspect(token)
	if err != nil {
		res := response.IntrospectResponse{
			Active: false,
			User:   nil,
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)
		return
	}

	// Return the response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(introspectResponse)

}
