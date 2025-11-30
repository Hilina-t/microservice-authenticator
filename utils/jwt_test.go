package utils

import (
	"testing"
	"time"

	"github.com/Hilina-t/microservice-authenticator/models"
)

func TestGenerateAndValidateJWT(t *testing.T) {
	secret := "test-secret-key"
	user := &models.User{
		ID:       "123",
		Email:    "test@example.com",
		Name:     "Test User",
		Provider: "google",
		Roles:    []string{"user"},
		Created:  time.Now(),
	}

	// Generate JWT
	token, err := GenerateJWT(user, secret, 24)
	if err != nil {
		t.Fatalf("Failed to generate JWT: %v", err)
	}

	if token == "" {
		t.Fatal("Generated token is empty")
	}

	// Validate JWT
	claims, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("Failed to validate JWT: %v", err)
	}

	// Verify claims
	if claims.UserID != user.ID {
		t.Errorf("UserID = %s, want %s", claims.UserID, user.ID)
	}
	if claims.Email != user.Email {
		t.Errorf("Email = %s, want %s", claims.Email, user.Email)
	}
	if claims.Name != user.Name {
		t.Errorf("Name = %s, want %s", claims.Name, user.Name)
	}
	if claims.Provider != user.Provider {
		t.Errorf("Provider = %s, want %s", claims.Provider, user.Provider)
	}
	if len(claims.Roles) != len(user.Roles) {
		t.Errorf("Roles length = %d, want %d", len(claims.Roles), len(user.Roles))
	}
}

func TestValidateJWT_InvalidToken(t *testing.T) {
	secret := "test-secret-key"

	tests := []struct {
		name  string
		token string
	}{
		{"Empty token", ""},
		{"Invalid format", "not-a-jwt-token"},
		{"Malformed JWT", "header.payload"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ValidateJWT(tt.token, secret)
			if err == nil {
				t.Error("Expected error for invalid token, got nil")
			}
		})
	}
}

func TestValidateJWT_WrongSecret(t *testing.T) {
	user := &models.User{
		ID:       "123",
		Email:    "test@example.com",
		Name:     "Test User",
		Provider: "google",
		Roles:    []string{"user"},
		Created:  time.Now(),
	}

	// Generate with one secret
	token, err := GenerateJWT(user, "secret1", 24)
	if err != nil {
		t.Fatalf("Failed to generate JWT: %v", err)
	}

	// Try to validate with different secret
	_, err = ValidateJWT(token, "secret2")
	if err == nil {
		t.Error("Expected error when validating with wrong secret, got nil")
	}
}

func TestGenerateJWT_ExpirationTime(t *testing.T) {
	secret := "test-secret-key"
	user := &models.User{
		ID:       "123",
		Email:    "test@example.com",
		Name:     "Test User",
		Provider: "google",
		Roles:    []string{"user"},
		Created:  time.Now(),
	}

	expirationHours := 1
	token, err := GenerateJWT(user, secret, expirationHours)
	if err != nil {
		t.Fatalf("Failed to generate JWT: %v", err)
	}

	claims, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("Failed to validate JWT: %v", err)
	}

	// Check expiration is approximately correct (within 1 minute)
	expectedExpiration := time.Now().Add(time.Hour * time.Duration(expirationHours))
	actualExpiration := claims.ExpiresAt.Time

	diff := actualExpiration.Sub(expectedExpiration)
	if diff < -time.Minute || diff > time.Minute {
		t.Errorf("Token expiration time differs by %v, expected within 1 minute", diff)
	}
}
