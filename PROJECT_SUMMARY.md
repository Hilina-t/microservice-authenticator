# Project Summary: Identity and Authorization Gateway (IAG)

## Overview

This project implements a complete Identity and Authorization Gateway (IAG) in Go that handles user authentication via external Identity Providers (Google, Okta, Azure AD) using OAuth 2.0/OIDC and enforces Role-Based Access Control (RBAC) across backend services using JWT tokens.

## Key Features Implemented

### 1. OAuth 2.0 / OpenID Connect (OIDC) Integration ✓

**Location**: `auth/oauth.go`, `handlers/auth.go`

- **Identity Provider Support**:
  - Google OAuth 2.0
  - Okta OIDC
  - Azure AD OAuth 2.0
  
- **OAuth Flow Implementation**:
  - Authorization URL generation with CSRF protection (state parameter)
  - Authorization code exchange for access tokens
  - User information retrieval from IdP
  - Secure callback handling with state validation

### 2. JWT Token Management ✓

**Location**: `utils/jwt.go`

- **JWT Generation**:
  - Creates signed JWT tokens using HS256 algorithm
  - Includes user claims (ID, email, name, roles, provider)
  - Configurable expiration time
  - Standard JWT claims (exp, iat, nbf, iss, sub)

- **JWT Validation**:
  - Signature verification
  - Expiration checking
  - Claims extraction

### 3. Role-Based Access Control (RBAC) ✓

**Location**: `models/user.go`, `middleware/rbac.go`

- **Predefined Roles**:
  - **Admin**: Full access to all resources (`*:*`)
  - **User**: Read/update profile, read/create data
  - **Viewer**: Read-only access to profile and data

- **Permission System**:
  - Resource-action based permissions
  - Wildcard support for flexible access control
  - Per-role permission mappings
  - `HasPermission()` and `HasRole()` methods

### 4. Policy Enforcement ✓

**Location**: `middleware/auth.go`, `middleware/rbac.go`

- **Authentication Middleware**:
  - Extracts and validates JWT tokens from Authorization header
  - Rejects requests without valid tokens (401 Unauthorized)
  - Adds authenticated user to request context

- **Authorization Middleware**:
  - `RequireRole()`: Checks if user has specific role(s)
  - `RequirePermission()`: Checks if user has specific permission
  - Returns 403 Forbidden for insufficient permissions

### 5. Identity Gateway Service ✓

**Location**: `main.go`, `handlers/`

- **Public Endpoints**:
  - `/` - Service information
  - `/health` - Health check
  - `/auth/login` - Initiates OAuth flow
  - `/auth/callback` - OAuth callback handler

- **Protected Endpoints**:
  - `/auth/profile` - User profile (requires authentication)
  - `/auth/logout` - Logout handler
  - `/api/admin` - Admin-only endpoint
  - `/api/user/data` - User/admin endpoint
  - `/api/viewer/data` - Viewer/user/admin endpoint
  - `/api/data/create` - Permission-based endpoint

### 6. Configuration Management ✓

**Location**: `config/config.go`, `.env.example`

- Environment-based configuration
- Support for multiple OAuth providers
- Configurable JWT settings
- RBAC enable/disable option
- Provider-specific endpoint configuration

## Project Structure

```
.
├── auth/                   # OAuth/OIDC authentication
│   └── oauth.go           # OAuth service implementation
├── config/                # Configuration management
│   └── config.go          # Config loading and validation
├── handlers/              # HTTP request handlers
│   ├── auth.go           # Authentication endpoints
│   └── protected.go      # Protected resource endpoints
├── middleware/            # HTTP middleware
│   ├── auth.go           # JWT authentication middleware
│   └── rbac.go           # RBAC authorization middleware
├── models/                # Data models
│   ├── user.go           # User model with RBAC methods
│   └── user_test.go      # Unit tests for user model
├── utils/                 # Utility functions
│   ├── jwt.go            # JWT generation/validation
│   └── jwt_test.go       # Unit tests for JWT
├── docs/                  # Documentation
│   ├── API.md            # API reference
│   ├── SETUP.md          # Setup instructions
│   └── TESTING.md        # Testing guide
├── examples/              # Example code
│   ├── client_example.go # Go client example
│   └── README.md         # Examples documentation
├── tests/                 # Integration tests
│   └── integration_test.sh
├── main.go               # Application entry point
├── Dockerfile            # Docker container definition
├── docker-compose.yml    # Docker Compose configuration
├── .env.example          # Example environment variables
├── .gitignore           # Git ignore patterns
└── README.md            # Project overview
```

## Technical Implementation Details

### OAuth 2.0 Flow

1. **Login Request** (`/auth/login`):
   - Generates CSRF protection state token
   - Stores state in cookie
   - Redirects to IdP authorization URL

2. **OAuth Callback** (`/auth/callback`):
   - Validates state parameter
   - Exchanges authorization code for access token
   - Fetches user information from IdP
   - Generates JWT token with user data and roles
   - Returns JWT to client

3. **Token Usage**:
   - Client includes JWT in `Authorization: Bearer {token}` header
   - Middleware validates JWT on each request
   - User context available for authorization checks

### RBAC Implementation

1. **Role Assignment**:
   - Default role: `user`
   - Can be customized based on user attributes
   - Multiple roles per user supported

2. **Permission Checking**:
   - Hierarchical permission system
   - Wildcard support (`*:*`, `resource:*`)
   - Exact match for specific permissions

3. **Middleware Chain**:
   ```
   Request → AuthMiddleware → RequireRole/RequirePermission → Handler
   ```

### JWT Token Structure

```json
{
  "user_id": "123456",
  "email": "user@example.com",
  "name": "John Doe",
  "roles": ["user"],
  "provider": "google",
  "exp": 1234567890,
  "iat": 1234567890,
  "nbf": 1234567890,
  "iss": "microservice-authenticator",
  "sub": "123456"
}
```

## Testing

### Unit Tests
- Model tests: User role and permission checks
- JWT tests: Token generation, validation, expiration
- All tests passing ✓

### Integration Tests
- Health check endpoint
- Root endpoint
- Protected endpoints without authentication
- Protected endpoints with invalid tokens
- OAuth login redirect
- All tests passing ✓

### Security Analysis
- CodeQL analysis completed
- No security vulnerabilities found ✓

## Deployment Options

### Local Development
```bash
go run main.go
```

### Production Build
```bash
go build -o microservice-authenticator
./microservice-authenticator
```

### Docker
```bash
docker build -t microservice-authenticator .
docker run -p 8080:8080 --env-file .env microservice-authenticator
```

### Docker Compose
```bash
docker-compose up -d
```

## Configuration Requirements

### Minimum Required Environment Variables
- `OAUTH_PROVIDER`: google/okta/azure
- `OAUTH_CLIENT_ID`: OAuth client ID
- `OAUTH_CLIENT_SECRET`: OAuth client secret
- `JWT_SECRET`: Secret for signing JWTs (change in production!)

### Provider-Specific Variables
- **Okta**: `OKTA_DOMAIN`
- **Azure**: `AZURE_TENANT_ID`

## Security Features

1. **CSRF Protection**: State parameter in OAuth flow
2. **JWT Signing**: HS256 algorithm with secret key
3. **Token Expiration**: Configurable expiration time
4. **Secure Cookies**: HttpOnly, SameSite protection
5. **Role-Based Access**: Fine-grained permission control
6. **Input Validation**: Token and state validation
7. **No Hardcoded Secrets**: All secrets via environment variables

## Performance Characteristics

- Stateless JWT validation (no database lookup per request)
- Efficient middleware chain
- Minimal memory footprint
- Fast HTTP router (standard library)
- Concurrent request handling

## Extension Points

### Adding New OAuth Providers
1. Add case to `LoadConfig()` in `config/config.go`
2. Implement user info extraction in `GetUserInfo()` in `auth/oauth.go`

### Adding New Roles
1. Define role constant in `models/user.go`
2. Add permission mappings in `RolePermissions`

### Adding New Protected Endpoints
1. Create handler function
2. Add route in `main.go`
3. Apply appropriate middleware

### Custom Role Assignment
Modify `GetUserInfo()` in `auth/oauth.go` to assign roles based on:
- Email domain
- User attributes from IdP
- External database lookup
- Custom business logic

## Documentation

- **README.md**: Project overview and quick start
- **docs/API.md**: Complete API reference with examples
- **docs/SETUP.md**: Detailed setup instructions for each OAuth provider
- **docs/TESTING.md**: Comprehensive testing guide
- **examples/README.md**: Example code and integration patterns

## Dependencies

- `github.com/golang-jwt/jwt/v5`: JWT token handling
- `golang.org/x/oauth2`: OAuth 2.0 client implementation

## Metrics

- **Total Files**: 23 source files
- **Lines of Code**: ~1,700 lines (excluding tests and docs)
- **Test Coverage**: Unit tests for core functionality
- **Security Score**: No vulnerabilities detected
- **Documentation**: 4 comprehensive guides + inline comments

## Success Criteria Met ✓

- [x] OAuth 2.0/OIDC integration with multiple providers
- [x] JWT token generation and validation
- [x] Role-Based Access Control implementation
- [x] Policy enforcement middleware
- [x] Callback handler with token forwarding
- [x] Complete authentication flow
- [x] Comprehensive testing
- [x] Production-ready deployment options
- [x] Complete documentation
- [x] Security validated

## Future Enhancements (Optional)

1. Token refresh mechanism
2. Database integration for user/role persistence
3. Rate limiting
4. Audit logging
5. Multi-factor authentication
6. Session management
7. Token revocation
8. CORS configuration
9. Metrics and monitoring
10. API versioning

## Conclusion

This project provides a complete, production-ready Identity and Authorization Gateway that implements OAuth 2.0/OIDC authentication with RBAC using JWT tokens. The system is secure, well-tested, documented, and ready for deployment in microservices architectures.
