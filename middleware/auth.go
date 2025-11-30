package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Hilina-t/microservice-authenticator/config"
	"github.com/Hilina-t/microservice-authenticator/models"
	"github.com/Hilina-t/microservice-authenticator/utils"
)

type contextKey string

const UserContextKey contextKey = "user"

// AuthMiddleware validates JWT tokens
func AuthMiddleware(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			// Check for Bearer token
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]

			// Validate token
			claims, err := utils.ValidateJWT(tokenString, cfg.JWTSecret)
			if err != nil {
				http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
				return
			}

			// Create user from claims
			user := &models.User{
				ID:       claims.UserID,
				Email:    claims.Email,
				Name:     claims.Name,
				Roles:    claims.Roles,
				Provider: claims.Provider,
			}

			// Add user to context
			ctx := context.WithValue(r.Context(), UserContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserFromContext retrieves the user from the request context
func GetUserFromContext(ctx context.Context) (*models.User, bool) {
	user, ok := ctx.Value(UserContextKey).(*models.User)
	return user, ok
}
