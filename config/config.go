package config

import (
	"fmt"
	"os"
)

// Config holds the application configuration
type Config struct {
	// Server settings
	ServerPort string

	// OAuth/OIDC settings
	OAuthProvider     string // google, okta, azure
	OAuthClientID     string
	OAuthClientSecret string
	OAuthRedirectURL  string
	OAuthAuthURL      string
	OAuthTokenURL     string
	OAuthUserInfoURL  string
	OAuthScopes       []string

	// JWT settings
	JWTSecret     string
	JWTExpiration int // in hours

	// RBAC settings
	EnableRBAC bool
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	config := &Config{
		ServerPort:        getEnv("SERVER_PORT", "8080"),
		OAuthProvider:     getEnv("OAUTH_PROVIDER", "google"),
		OAuthClientID:     getEnv("OAUTH_CLIENT_ID", ""),
		OAuthClientSecret: getEnv("OAUTH_CLIENT_SECRET", ""),
		OAuthRedirectURL:  getEnv("OAUTH_REDIRECT_URL", "http://localhost:8080/auth/callback"),
		JWTSecret:         getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		JWTExpiration:     getEnvAsInt("JWT_EXPIRATION_HOURS", 24),
		EnableRBAC:        getEnvAsBool("ENABLE_RBAC", true),
	}

	// Set provider-specific OAuth endpoints
	switch config.OAuthProvider {
	case "google":
		config.OAuthAuthURL = "https://accounts.google.com/o/oauth2/v2/auth"
		config.OAuthTokenURL = "https://oauth2.googleapis.com/token"
		config.OAuthUserInfoURL = "https://www.googleapis.com/oauth2/v2/userinfo"
		config.OAuthScopes = []string{"openid", "profile", "email"}
	case "okta":
		oktaDomain := getEnv("OKTA_DOMAIN", "")
		if oktaDomain == "" {
			return nil, fmt.Errorf("OKTA_DOMAIN is required for Okta provider")
		}
		config.OAuthAuthURL = fmt.Sprintf("https://%s/oauth2/v1/authorize", oktaDomain)
		config.OAuthTokenURL = fmt.Sprintf("https://%s/oauth2/v1/token", oktaDomain)
		config.OAuthUserInfoURL = fmt.Sprintf("https://%s/oauth2/v1/userinfo", oktaDomain)
		config.OAuthScopes = []string{"openid", "profile", "email"}
	case "azure":
		tenantID := getEnv("AZURE_TENANT_ID", "common")
		config.OAuthAuthURL = fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/authorize", tenantID)
		config.OAuthTokenURL = fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", tenantID)
		config.OAuthUserInfoURL = "https://graph.microsoft.com/v1.0/me"
		config.OAuthScopes = []string{"openid", "profile", "email"}
	default:
		return nil, fmt.Errorf("unsupported OAuth provider: %s", config.OAuthProvider)
	}

	// Validate required config
	if config.OAuthClientID == "" {
		return nil, fmt.Errorf("OAUTH_CLIENT_ID is required")
	}
	if config.OAuthClientSecret == "" {
		return nil, fmt.Errorf("OAUTH_CLIENT_SECRET is required")
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	var value int
	fmt.Sscanf(valueStr, "%d", &value)
	return value
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	return valueStr == "true" || valueStr == "1"
}
