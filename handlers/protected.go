package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Hilina-t/microservice-authenticator/middleware"
)

// ProtectedHandler handles protected resource requests
type ProtectedHandler struct{}

// NewProtectedHandler creates a new protected resource handler
func NewProtectedHandler() *ProtectedHandler {
	return &ProtectedHandler{}
}

// AdminOnly is an example endpoint that requires admin role
func (h *ProtectedHandler) AdminOnly(w http.ResponseWriter, r *http.Request) {
	user, _ := middleware.GetUserFromContext(r.Context())

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Welcome to admin area",
		"user":    user,
	})
}

// UserData is an example endpoint that requires user role
func (h *ProtectedHandler) UserData(w http.ResponseWriter, r *http.Request) {
	user, _ := middleware.GetUserFromContext(r.Context())

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User data access granted",
		"user":    user,
		"data": map[string]string{
			"resource": "user-data",
			"info":     "This is protected user data",
		},
	})
}

// ViewerData is an example endpoint that requires viewer role
func (h *ProtectedHandler) ViewerData(w http.ResponseWriter, r *http.Request) {
	user, _ := middleware.GetUserFromContext(r.Context())

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Viewer data access granted",
		"user":    user,
		"data": map[string]string{
			"resource": "viewer-data",
			"info":     "This is read-only data",
		},
	})
}
