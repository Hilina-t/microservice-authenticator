package models

import "time"

// User represents an authenticated user
type User struct {
	ID       string    `json:"id"`
	Email    string    `json:"email"`
	Name     string    `json:"name"`
	Picture  string    `json:"picture,omitempty"`
	Provider string    `json:"provider"`
	Roles    []string  `json:"roles"`
	Created  time.Time `json:"created"`
}

// Role represents a role in the RBAC system
type Role string

const (
	RoleAdmin  Role = "admin"
	RoleUser   Role = "user"
	RoleViewer Role = "viewer"
)

// Permission represents an action that can be performed
type Permission struct {
	Resource string `json:"resource"`
	Action   string `json:"action"`
}

// RolePermissions maps roles to their permissions
var RolePermissions = map[Role][]Permission{
	RoleAdmin: {
		{Resource: "*", Action: "*"},
	},
	RoleUser: {
		{Resource: "profile", Action: "read"},
		{Resource: "profile", Action: "update"},
		{Resource: "data", Action: "read"},
		{Resource: "data", Action: "create"},
	},
	RoleViewer: {
		{Resource: "profile", Action: "read"},
		{Resource: "data", Action: "read"},
	},
}

// HasPermission checks if a user has permission for a resource and action
func (u *User) HasPermission(resource, action string) bool {
	for _, roleName := range u.Roles {
		role := Role(roleName)
		permissions, exists := RolePermissions[role]
		if !exists {
			continue
		}

		for _, perm := range permissions {
			// Check for wildcard permissions
			if perm.Resource == "*" && perm.Action == "*" {
				return true
			}
			// Check for specific resource with wildcard action
			if perm.Resource == resource && perm.Action == "*" {
				return true
			}
			// Check for exact match
			if perm.Resource == resource && perm.Action == action {
				return true
			}
		}
	}
	return false
}

// HasRole checks if a user has a specific role
func (u *User) HasRole(role string) bool {
	for _, r := range u.Roles {
		if r == role {
			return true
		}
	}
	return false
}
