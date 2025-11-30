# API Documentation

## Identity and Authorization Gateway (IAG)

This document describes the API endpoints provided by the IAG service.

## Public Endpoints

### GET /
Returns service information.

**Response:**
```json
{
  "message": "Identity and Authorization Gateway (IAG)",
  "version": "1.0.0"
}
```

### GET /health
Health check endpoint.

**Response:**
```json
{
  "status": "healthy"
}
```

## Authentication Endpoints

### GET /auth/login
Initiates the OAuth 2.0/OIDC authentication flow. Redirects to the configured Identity Provider.

**Response:** HTTP 307 redirect to IdP

### GET /auth/callback
OAuth callback endpoint. Handles the response from the Identity Provider, exchanges the authorization code for tokens, and returns a JWT.

**Query Parameters:**
- `code`: Authorization code from IdP
- `state`: State parameter for CSRF protection

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "123456",
    "email": "user@example.com",
    "name": "John Doe",
    "picture": "https://...",
    "provider": "google",
    "roles": ["user"],
    "created": "2024-01-01T00:00:00Z"
  }
}
```

### GET /auth/profile
Returns the authenticated user's profile.

**Headers:**
- `Authorization`: Bearer {jwt_token}

**Response:**
```json
{
  "id": "123456",
  "email": "user@example.com",
  "name": "John Doe",
  "picture": "https://...",
  "provider": "google",
  "roles": ["user"],
  "created": "2024-01-01T00:00:00Z"
}
```

### GET /auth/logout
Logs out the current user.

**Headers:**
- `Authorization`: Bearer {jwt_token}

**Response:**
```json
{
  "message": "Logged out successfully"
}
```

## RBAC Protected Endpoints

### GET /api/admin
Admin-only endpoint. Requires `admin` role.

**Headers:**
- `Authorization`: Bearer {jwt_token}

**Response:**
```json
{
  "message": "Welcome to admin area",
  "user": {...}
}
```

**Error (403):**
```json
{
  "error": "Insufficient permissions"
}
```

### GET /api/user/data
User data endpoint. Requires `user` or `admin` role.

**Headers:**
- `Authorization`: Bearer {jwt_token}

**Response:**
```json
{
  "message": "User data access granted",
  "user": {...},
  "data": {
    "resource": "user-data",
    "info": "This is protected user data"
  }
}
```

### GET /api/viewer/data
Viewer data endpoint. Requires `viewer`, `user`, or `admin` role.

**Headers:**
- `Authorization`: Bearer {jwt_token}

**Response:**
```json
{
  "message": "Viewer data access granted",
  "user": {...},
  "data": {
    "resource": "viewer-data",
    "info": "This is read-only data"
  }
}
```

### POST /api/data/create
Data creation endpoint. Requires `create` permission on `data` resource.

**Headers:**
- `Authorization`: Bearer {jwt_token}

**Response:**
```json
{
  "message": "Data creation allowed",
  "user": "user@example.com"
}
```

## Authentication Flow

1. Client initiates login by navigating to `/auth/login`
2. Server redirects to Identity Provider (Google/Okta/Azure AD)
3. User authenticates with Identity Provider
4. Identity Provider redirects back to `/auth/callback` with authorization code
5. Server exchanges code for access token
6. Server fetches user information from Identity Provider
7. Server generates JWT token with user info and roles
8. Server returns JWT to client
9. Client includes JWT in `Authorization: Bearer {token}` header for subsequent requests

## JWT Token Format

The JWT token contains the following claims:

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

## RBAC Roles

### Admin
- Full access to all resources
- Permissions: `*:*`

### User
- Can read and update profile
- Can read and create data
- Permissions:
  - `profile:read`
  - `profile:update`
  - `data:read`
  - `data:create`

### Viewer
- Read-only access
- Permissions:
  - `profile:read`
  - `data:read`

## Error Responses

### 401 Unauthorized
```json
{
  "error": "Authorization header required"
}
```

### 403 Forbidden
```json
{
  "error": "Insufficient permissions",
  "resource": "data",
  "action": "create"
}
```

### 500 Internal Server Error
```json
{
  "error": "Internal server error message"
}
```
