# Architecture Documentation

## System Overview

The Identity and Authorization Gateway (IAG) is a centralized authentication and authorization service that implements OAuth 2.0/OIDC for user authentication and JWT-based authorization with Role-Based Access Control (RBAC).

## Architecture Diagram

```
┌─────────────┐
│   Client    │
│  (Browser/  │
│   Mobile)   │
└──────┬──────┘
       │
       │ 1. GET /auth/login
       ▼
┌──────────────────────────────────────────────────────────┐
│         Identity & Authorization Gateway (IAG)           │
│                                                           │
│  ┌────────────────────────────────────────────────────┐  │
│  │  Authentication Flow                                │  │
│  │  1. /auth/login → Redirect to IdP                  │  │
│  │  2. User authenticates with IdP                    │  │
│  │  3. /auth/callback ← Authorization code            │  │
│  │  4. Exchange code for access token                 │  │
│  │  5. Fetch user info from IdP                       │  │
│  │  6. Generate JWT with user info & roles            │  │
│  │  7. Return JWT to client                           │  │
│  └────────────────────────────────────────────────────┘  │
│                                                           │
│  ┌────────────────────────────────────────────────────┐  │
│  │  Authorization Middleware                          │  │
│  │  • JWT Validation                                  │  │
│  │  • Role Checking                                   │  │
│  │  • Permission Verification                         │  │
│  └────────────────────────────────────────────────────┘  │
└──────┬────────────────────────────────────────┬──────────┘
       │                                        │
       │ 2. OAuth Flow                         │ 3. Protected
       │                                        │    Request
       ▼                                        ▼
┌──────────────┐                        ┌──────────────┐
│   Identity   │                        │   Backend    │
│   Provider   │                        │   Services   │
│              │                        │              │
│ • Google     │                        │ • Service A  │
│ • Okta       │                        │ • Service B  │
│ • Azure AD   │                        │ • Service C  │
└──────────────┘                        └──────────────┘
```

## Authentication Flow

### Detailed Step-by-Step Flow

```
1. User initiates login
   Client → IAG: GET /auth/login
   
2. IAG redirects to Identity Provider
   IAG → Client: HTTP 307 Redirect
   Location: https://idp.example.com/oauth/authorize?
             client_id=xxx&
             redirect_uri=http://iag/auth/callback&
             state=random-csrf-token&
             scope=openid profile email
   
3. User authenticates with IdP
   Client → IdP: User provides credentials
   IdP → IdP: Validates credentials
   
4. IdP redirects back to IAG with code
   IdP → Client: HTTP 302 Redirect
   Location: http://iag/auth/callback?
             code=auth-code&
             state=random-csrf-token
   
5. IAG receives callback
   Client → IAG: GET /auth/callback?code=auth-code&state=csrf-token
   
6. IAG validates state (CSRF protection)
   IAG: Verify state matches original
   
7. IAG exchanges code for access token
   IAG → IdP: POST /oauth/token
             code=auth-code&
             client_id=xxx&
             client_secret=yyy&
             grant_type=authorization_code
   
8. IdP returns access token
   IdP → IAG: {
             "access_token": "...",
             "token_type": "Bearer",
             "expires_in": 3600
           }
   
9. IAG fetches user information
   IAG → IdP: GET /userinfo
             Authorization: Bearer access-token
   
10. IdP returns user data
    IdP → IAG: {
              "id": "123456",
              "email": "user@example.com",
              "name": "John Doe"
            }
   
11. IAG generates JWT token
    IAG: Generate JWT with user info and roles
         Claims: {user_id, email, name, roles, exp, iat}
         Sign with secret key (HS256)
   
12. IAG returns JWT to client
    IAG → Client: {
                  "token": "eyJhbGciOiJIUzI1NiIs...",
                  "user": {user details}
                }
```

## Authorization Flow

### Protected Request Flow

```
1. Client makes request with JWT
   Client → IAG: GET /api/resource
                 Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
   
2. Auth Middleware validates JWT
   IAG: Extract token from header
        Verify signature
        Check expiration
        Extract claims
   
3. User context added to request
   IAG: Create user object from claims
        Add to request context
   
4. RBAC Middleware checks permissions
   IAG: Get required role/permission
        Check user has required access
        
5a. If authorized:
    IAG → Handler: Process request
    Handler → IAG: Response data
    IAG → Client: 200 OK + data
   
5b. If unauthorized:
    IAG → Client: 401 Unauthorized
   
5c. If forbidden:
    IAG → Client: 403 Forbidden
```

## Component Architecture

### Layer Structure

```
┌─────────────────────────────────────────────┐
│           HTTP Handler Layer                │
│  (handlers/auth.go, handlers/protected.go) │
└───────────────────┬─────────────────────────┘
                    │
┌───────────────────▼─────────────────────────┐
│          Middleware Layer                   │
│  • Authentication (middleware/auth.go)      │
│  • RBAC (middleware/rbac.go)               │
└───────────────────┬─────────────────────────┘
                    │
┌───────────────────▼─────────────────────────┐
│          Service Layer                      │
│  • OAuth Service (auth/oauth.go)           │
│  • JWT Utils (utils/jwt.go)                │
└───────────────────┬─────────────────────────┘
                    │
┌───────────────────▼─────────────────────────┐
│          Model Layer                        │
│  • User Model (models/user.go)             │
│  • Role & Permission Models                 │
└───────────────────┬─────────────────────────┘
                    │
┌───────────────────▼─────────────────────────┐
│        Configuration Layer                  │
│  (config/config.go, .env)                  │
└─────────────────────────────────────────────┘
```

## Data Flow

### JWT Token Lifecycle

```
Creation:
  User Data → Generate JWT → Sign with Secret → Token String
  
Validation:
  Token String → Parse JWT → Verify Signature → Extract Claims
  
Usage:
  Claims → User Object → Authorization Check → Access Granted/Denied
```

## Security Architecture

### Security Layers

```
┌─────────────────────────────────────────────┐
│  Layer 1: OAuth 2.0 Security                │
│  • CSRF Protection (state parameter)        │
│  • Secure redirect URIs                     │
│  • Client secret protection                 │
└───────────────────┬─────────────────────────┘
                    │
┌───────────────────▼─────────────────────────┐
│  Layer 2: JWT Security                      │
│  • HS256 signing algorithm                  │
│  • Secret key protection                    │
│  • Token expiration                         │
│  • Claims validation                        │
└───────────────────┬─────────────────────────┘
                    │
┌───────────────────▼─────────────────────────┐
│  Layer 3: RBAC Security                     │
│  • Role-based access control                │
│  • Permission checking                      │
│  • Least privilege principle                │
└───────────────────┬─────────────────────────┘
                    │
┌───────────────────▼─────────────────────────┐
│  Layer 4: Application Security              │
│  • Input validation                         │
│  • Secure headers                           │
│  • HTTPS enforcement (production)           │
└─────────────────────────────────────────────┘
```

## Role-Based Access Control (RBAC)

### Role Hierarchy

```
┌──────────────────┐
│      Admin       │ ← Full access (*.*)
│   (Superuser)    │
└────────┬─────────┘
         │
         ▼
┌──────────────────┐
│      User        │ ← Read/Write access
│   (Standard)     │   - profile: read, update
└────────┬─────────┘   - data: read, create
         │
         ▼
┌──────────────────┐
│     Viewer       │ ← Read-only access
│  (Read-only)     │   - profile: read
└──────────────────┘   - data: read
```

### Permission Matrix

| Role   | Profile (read) | Profile (update) | Data (read) | Data (create) | Data (delete) |
|--------|----------------|------------------|-------------|---------------|---------------|
| Admin  | ✓              | ✓                | ✓           | ✓             | ✓             |
| User   | ✓              | ✓                | ✓           | ✓             | ✗             |
| Viewer | ✓              | ✗                | ✓           | ✗             | ✗             |

## Deployment Architecture

### Single Instance Deployment

```
┌────────────────────────────────────────┐
│         Load Balancer / Proxy          │
│            (Nginx/HAProxy)             │
└───────────────┬────────────────────────┘
                │
                ▼
┌────────────────────────────────────────┐
│              IAG Service               │
│         (Port 8080)                    │
│  ┌──────────────────────────────────┐  │
│  │  • OAuth Endpoints               │  │
│  │  • JWT Validation                │  │
│  │  • RBAC Enforcement              │  │
│  └──────────────────────────────────┘  │
└───────────────┬────────────────────────┘
                │
    ┌───────────┴───────────┐
    │                       │
    ▼                       ▼
┌─────────┐           ┌─────────┐
│  IdP    │           │ Backend │
│ Services│           │ Services│
└─────────┘           └─────────┘
```

### High Availability Deployment

```
┌────────────────────────────────────────┐
│         Load Balancer                  │
└───────────┬────────────────────────────┘
            │
    ┌───────┼───────┐
    │       │       │
    ▼       ▼       ▼
┌─────┐ ┌─────┐ ┌─────┐
│IAG 1│ │IAG 2│ │IAG 3│ ← Stateless (JWT)
└──┬──┘ └──┬──┘ └──┬──┘
   │       │       │
   └───────┼───────┘
           │
    ┌──────┴──────┐
    │             │
    ▼             ▼
┌────────┐   ┌────────┐
│  IdP   │   │Backend │
│Services│   │Services│
└────────┘   └────────┘
```

## Configuration Management

### Environment-Based Configuration

```
Development:
  • HTTP (localhost)
  • Debug logging
  • Short token expiration
  • Test OAuth credentials

Staging:
  • HTTPS
  • Standard logging
  • Standard token expiration
  • Staging OAuth credentials

Production:
  • HTTPS (enforced)
  • Production logging
  • Secure token expiration
  • Production OAuth credentials
  • Rate limiting enabled
  • Monitoring enabled
```

## API Gateway Pattern

### IAG as API Gateway

```
External Clients
       │
       ▼
┌──────────────┐
│     IAG      │ ← Central authentication/authorization point
└──────┬───────┘
       │
   ┌───┴────┬────────┬────────┐
   │        │        │        │
   ▼        ▼        ▼        ▼
┌────┐  ┌────┐  ┌────┐  ┌────┐
│Svc1│  │Svc2│  │Svc3│  │Svc4│ ← Backend services trust IAG's JWTs
└────┘  └────┘  └────┘  └────┘
```

### Microservices Integration

Each backend service validates the JWT token:

```go
func validateRequest(r *http.Request) (*User, error) {
    token := extractToken(r)
    claims := validateJWT(token, sharedSecret)
    return claims.toUser(), nil
}
```

## Scalability Considerations

### Horizontal Scaling

- **Stateless Design**: JWT tokens eliminate need for session storage
- **No Database**: Authentication state in tokens
- **Load Balancing**: Any instance can handle any request
- **Auto-scaling**: Scale based on request rate

### Performance Optimization

- **Token Caching**: Cache user roles/permissions
- **Connection Pooling**: HTTP client pools for IdP
- **Async Operations**: Non-blocking I/O
- **Efficient Parsing**: Fast JWT validation

## Monitoring and Observability

### Key Metrics

```
Authentication Metrics:
  • Login requests/sec
  • OAuth failures
  • Token generation rate
  • Token validation rate

Authorization Metrics:
  • 401 Unauthorized rate
  • 403 Forbidden rate
  • Role distribution
  • Permission denials

Performance Metrics:
  • Response time (p50, p95, p99)
  • Throughput (req/sec)
  • Error rate
  • CPU/Memory usage
```

## Future Enhancements

### Planned Features

1. **Token Refresh**: Automatic token renewal
2. **Token Revocation**: Blacklist for revoked tokens
3. **Rate Limiting**: Per-user/IP rate limits
4. **Audit Logging**: Complete audit trail
5. **MFA Support**: Multi-factor authentication
6. **Session Management**: Optional session support
7. **Dynamic Roles**: Database-backed roles
8. **API Versioning**: Multiple API versions

## Best Practices

### Development

- Use environment variables for configuration
- Never commit secrets to version control
- Follow principle of least privilege
- Implement proper error handling
- Write comprehensive tests

### Deployment

- Use HTTPS in production
- Rotate JWT secrets regularly
- Monitor authentication patterns
- Set up alerting for failures
- Implement rate limiting

### Security

- Keep dependencies updated
- Regular security audits
- Follow OWASP guidelines
- Implement defense in depth
- Use secure communication channels
