package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Hilina-t/microservice-authenticator/config"
	"github.com/Hilina-t/microservice-authenticator/models"
	"golang.org/x/oauth2"
)

// OAuthService handles OAuth/OIDC authentication
type OAuthService struct {
	config      *config.Config
	oauthConfig *oauth2.Config
}

// NewOAuthService creates a new OAuth service
func NewOAuthService(cfg *config.Config) *OAuthService {
	oauthConfig := &oauth2.Config{
		ClientID:     cfg.OAuthClientID,
		ClientSecret: cfg.OAuthClientSecret,
		RedirectURL:  cfg.OAuthRedirectURL,
		Scopes:       cfg.OAuthScopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  cfg.OAuthAuthURL,
			TokenURL: cfg.OAuthTokenURL,
		},
	}

	return &OAuthService{
		config:      cfg,
		oauthConfig: oauthConfig,
	}
}

// GetAuthURL returns the authorization URL for OAuth flow
func (s *OAuthService) GetAuthURL(state string) string {
	return s.oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// ExchangeCode exchanges the authorization code for tokens
func (s *OAuthService) ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := s.oauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}
	return token, nil
}

// GetUserInfo fetches user information from the provider
func (s *OAuthService) GetUserInfo(ctx context.Context, token *oauth2.Token) (*models.User, error) {
	client := s.oauthConfig.Client(ctx, token)
	resp, err := client.Get(s.config.OAuthUserInfoURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get user info: status %d, body: %s", resp.StatusCode, string(body))
	}

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	user := &models.User{
		Provider: s.config.OAuthProvider,
		Created:  time.Now(),
	}

	// Extract user information based on provider
	switch s.config.OAuthProvider {
	case "google":
		user.ID = getStringField(userInfo, "id")
		user.Email = getStringField(userInfo, "email")
		user.Name = getStringField(userInfo, "name")
		user.Picture = getStringField(userInfo, "picture")
	case "okta":
		user.ID = getStringField(userInfo, "sub")
		user.Email = getStringField(userInfo, "email")
		user.Name = getStringField(userInfo, "name")
		user.Picture = getStringField(userInfo, "picture")
	case "azure":
		user.ID = getStringField(userInfo, "id")
		user.Email = getStringField(userInfo, "mail")
		if user.Email == "" {
			user.Email = getStringField(userInfo, "userPrincipalName")
		}
		user.Name = getStringField(userInfo, "displayName")
	}

	// Assign default role if no roles assigned
	if len(user.Roles) == 0 {
		user.Roles = []string{string(models.RoleUser)}
	}

	return user, nil
}

func getStringField(data map[string]interface{}, field string) string {
	if val, ok := data[field]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}
