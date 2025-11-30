package models

import (
	"testing"
	"time"
)

func TestUser_HasRole(t *testing.T) {
	user := &User{
		ID:      "123",
		Email:   "test@example.com",
		Name:    "Test User",
		Roles:   []string{"admin", "user"},
		Created: time.Now(),
	}

	tests := []struct {
		name     string
		role     string
		expected bool
	}{
		{"Has admin role", "admin", true},
		{"Has user role", "user", true},
		{"Does not have viewer role", "viewer", false},
		{"Does not have non-existent role", "superuser", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := user.HasRole(tt.role)
			if result != tt.expected {
				t.Errorf("HasRole(%s) = %v, want %v", tt.role, result, tt.expected)
			}
		})
	}
}

func TestUser_HasPermission(t *testing.T) {
	tests := []struct {
		name     string
		user     *User
		resource string
		action   string
		expected bool
	}{
		{
			name: "Admin has all permissions",
			user: &User{
				ID:    "1",
				Roles: []string{"admin"},
			},
			resource: "any-resource",
			action:   "any-action",
			expected: true,
		},
		{
			name: "User can read profile",
			user: &User{
				ID:    "2",
				Roles: []string{"user"},
			},
			resource: "profile",
			action:   "read",
			expected: true,
		},
		{
			name: "User can create data",
			user: &User{
				ID:    "3",
				Roles: []string{"user"},
			},
			resource: "data",
			action:   "create",
			expected: true,
		},
		{
			name: "User cannot delete data",
			user: &User{
				ID:    "4",
				Roles: []string{"user"},
			},
			resource: "data",
			action:   "delete",
			expected: false,
		},
		{
			name: "Viewer can read data",
			user: &User{
				ID:    "5",
				Roles: []string{"viewer"},
			},
			resource: "data",
			action:   "read",
			expected: true,
		},
		{
			name: "Viewer cannot create data",
			user: &User{
				ID:    "6",
				Roles: []string{"viewer"},
			},
			resource: "data",
			action:   "create",
			expected: false,
		},
		{
			name: "User with no roles has no permissions",
			user: &User{
				ID:    "7",
				Roles: []string{},
			},
			resource: "data",
			action:   "read",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.user.HasPermission(tt.resource, tt.action)
			if result != tt.expected {
				t.Errorf("HasPermission(%s, %s) = %v, want %v",
					tt.resource, tt.action, result, tt.expected)
			}
		})
	}
}
