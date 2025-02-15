// internal/iam/handler.go
package iam

import (
	"clinicplus/internal/shared/utils"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type AuthHandler struct {
	service AuthService
}

func NewAuthHandler(service AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		log.Println("Error decoding request body:", err)
		utils.SendJSONResponse(w, http.StatusBadRequest, nil, "Invalid request body", nil)
		return
	}

	// Call the service to handle login
	token, username, expirationTime, err := h.service.Login(loginRequest.Username, loginRequest.Password)
	if err != nil {
		log.Printf("Login failed for user: %s, error: %v\n", loginRequest.Username, err)
		utils.SendJSONResponse(w, http.StatusUnauthorized, nil, "Invalid username or password", nil)
		return
	}

	// Prepare login response
	loginResponse := struct {
		Token    string `json:"token"`
		Username string `json:"username"`
	}{
		Token:    token,
		Username: username,
	}

	// Send standardized response
	utils.SendJSONResponse(w, http.StatusOK, loginResponse, nil, map[string]interface{}{
		"expires_at": expirationTime.Format(time.RFC3339),
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Extract the token from the Authorization header
	authHeader := r.Header.Get("Authorization")

	// Call the service to handle logout
	if err := h.service.Logout(authHeader); err != nil {
		log.Printf("Logout failed for token: %s, error: %v\n", authHeader, err)
		utils.SendJSONResponse(w, http.StatusInternalServerError, nil, "Failed to process logout", nil)
		return
	}

	// Prepare logout response
	logoutResponse := map[string]string{
		"message": "Successfully logged out",
	}

	// Send standardized response
	utils.SendJSONResponse(w, http.StatusOK, logoutResponse, nil, map[string]interface{}{
		"timestamp": time.Now().Format(time.RFC3339),
	})
}
