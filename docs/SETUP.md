# Setup Guide

This guide will help you set up the Identity and Authorization Gateway (IAG) with different OAuth providers.

## Prerequisites

- Go 1.21 or higher
- An OAuth 2.0 / OIDC provider account (Google, Okta, or Azure AD)

## Installation

1. Clone the repository:
```bash
git clone https://github.com/Hilina-t/microservice-authenticator.git
cd microservice-authenticator
```

2. Install dependencies:
```bash
go mod download
```

3. Copy the example environment file:
```bash
cp .env.example .env
```

4. Configure your OAuth provider (see sections below)

5. Run the application:
```bash
go run main.go
```

## OAuth Provider Configuration

### Google OAuth Setup

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select an existing one
3. Enable the Google+ API
4. Go to "Credentials" and create OAuth 2.0 credentials
5. Add authorized redirect URI: `http://localhost:8080/auth/callback`
6. Copy the Client ID and Client Secret to your `.env` file:

```env
OAUTH_PROVIDER=google
OAUTH_CLIENT_ID=your-google-client-id
OAUTH_CLIENT_SECRET=your-google-client-secret
OAUTH_REDIRECT_URL=http://localhost:8080/auth/callback
```

### Okta OAuth Setup

1. Sign up for [Okta Developer Account](https://developer.okta.com/)
2. Create a new application (Web Application)
3. Set the login redirect URI: `http://localhost:8080/auth/callback`
4. Copy the Client ID and Client Secret
5. Note your Okta domain (e.g., `dev-123456.okta.com`)
6. Configure your `.env` file:

```env
OAUTH_PROVIDER=okta
OKTA_DOMAIN=dev-123456.okta.com
OAUTH_CLIENT_ID=your-okta-client-id
OAUTH_CLIENT_SECRET=your-okta-client-secret
OAUTH_REDIRECT_URL=http://localhost:8080/auth/callback
```

### Azure AD OAuth Setup

1. Go to [Azure Portal](https://portal.azure.com/)
2. Navigate to Azure Active Directory > App registrations
3. Create a new registration
4. Add a redirect URI: `http://localhost:8080/auth/callback`
5. Create a client secret under "Certificates & secrets"
6. Note your Tenant ID and Application (client) ID
7. Configure your `.env` file:

```env
OAUTH_PROVIDER=azure
AZURE_TENANT_ID=your-tenant-id
OAUTH_CLIENT_ID=your-application-client-id
OAUTH_CLIENT_SECRET=your-client-secret
OAUTH_REDIRECT_URL=http://localhost:8080/auth/callback
```

## Configuration Options

### Server Configuration
- `SERVER_PORT`: Port to run the server on (default: 8080)

### JWT Configuration
- `JWT_SECRET`: Secret key for signing JWT tokens (change in production!)
- `JWT_EXPIRATION_HOURS`: Token expiration time in hours (default: 24)

### RBAC Configuration
- `ENABLE_RBAC`: Enable role-based access control (default: true)

## Running the Application

### Development
```bash
go run main.go
```

### Production Build
```bash
go build -o microservice-authenticator
./microservice-authenticator
```

### Using Docker
```bash
docker build -t microservice-authenticator .
docker run -p 8080:8080 --env-file .env microservice-authenticator
```

## Testing the Flow

1. Start the application:
```bash
go run main.go
```

2. Open your browser and navigate to:
```
http://localhost:8080/auth/login
```

3. You'll be redirected to your OAuth provider to authenticate

4. After successful authentication, you'll receive a JWT token

5. Use the JWT token to access protected endpoints:
```bash
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" http://localhost:8080/auth/profile
```

## User Roles

By default, all authenticated users are assigned the `user` role. To implement custom role assignment:

1. Modify the `GetUserInfo` function in `auth/oauth.go`
2. Add logic to assign roles based on user attributes or external database
3. Update the `RolePermissions` map in `models/user.go` to define custom roles

Example role assignment based on email domain:
```go
if strings.HasSuffix(user.Email, "@admin.com") {
    user.Roles = []string{string(models.RoleAdmin)}
} else {
    user.Roles = []string{string(models.RoleUser)}
}
```

## Security Considerations

### Production Deployment

1. **Use HTTPS**: Always use HTTPS in production
2. **Change JWT Secret**: Use a strong, random JWT secret
3. **Secure Cookies**: Enable secure cookies in production
4. **Environment Variables**: Never commit `.env` file to version control
5. **Token Expiration**: Set appropriate token expiration times
6. **Rate Limiting**: Implement rate limiting on authentication endpoints
7. **CORS**: Configure CORS appropriately for your frontend

### JWT Best Practices

1. Store tokens securely (e.g., httpOnly cookies or secure storage)
2. Implement token refresh mechanism
3. Validate tokens on every request
4. Include minimal information in JWT payload
5. Use short expiration times

## Troubleshooting

### "OAUTH_CLIENT_ID is required" error
Make sure you've set the OAuth credentials in your `.env` file.

### "Invalid state parameter" error
This indicates a CSRF protection issue. Make sure cookies are enabled in your browser.

### "Failed to exchange code" error
- Verify your OAuth credentials are correct
- Check that the redirect URI matches exactly what's configured in your OAuth provider
- Ensure your OAuth application is enabled

### Port already in use
Change the `SERVER_PORT` in your `.env` file to use a different port.

## Development

### Project Structure
```
.
├── auth/           # OAuth/OIDC authentication logic
├── config/         # Configuration management
├── handlers/       # HTTP handlers
├── middleware/     # Authentication and RBAC middleware
├── models/         # Data models
├── utils/          # Utility functions (JWT)
├── docs/           # Documentation
├── main.go         # Application entry point
└── .env.example    # Example environment configuration
```

### Adding New Endpoints

1. Create handler function in `handlers/`
2. Add route in `main.go`
3. Apply appropriate middleware (AuthMiddleware, RequireRole, RequirePermission)

Example:
```go
mux.Handle("/api/new-endpoint",
    middleware.AuthMiddleware(cfg)(
        middleware.RequireRole("admin")(
            http.HandlerFunc(handler.NewEndpoint),
        ),
    ),
)
```
