# microservice-authenticator

Identity and Authorization Gateway (IAG) implemented in Go, enforcing JWT validation and Role-Based Access Control (RBAC) via OAuth 2.0/OIDC.

## Features

- **OAuth 2.0 / OpenID Connect (OIDC)** authentication with multiple Identity Providers:
  - Google
  - Okta
  - Azure AD
- **JWT Token Generation** with configurable expiration
- **Role-Based Access Control (RBAC)** with predefined roles (Admin, User, Viewer)
- **Permission-Based Authorization** for fine-grained access control
- **Token Validation Middleware** for protecting API endpoints
- **Policy Enforcement** across multiple backend services

## Architecture

The IAG service acts as a centralized authentication and authorization gateway:

1. **Authentication Flow**: Users authenticate via OAuth 2.0/OIDC with external Identity Providers
2. **Token Exchange**: Authorization codes are exchanged for access tokens
3. **JWT Generation**: User information is encoded into signed JWT tokens
4. **Token Forwarding**: JWTs are sent to clients and used for subsequent API calls
5. **RBAC Enforcement**: Middleware validates tokens and enforces role-based access policies

## Quick Start

### Prerequisites

- Go 1.21 or higher
- OAuth credentials from Google, Okta, or Azure AD

### Installation

1. Clone the repository:
```bash
git clone https://github.com/Hilina-t/microservice-authenticator.git
cd microservice-authenticator
```

2. Install dependencies:
```bash
go mod download
```

3. Configure environment variables:
```bash
cp .env.example .env
# Edit .env with your OAuth credentials
```

4. Run the service:
```bash
go run main.go
```

The service will start on `http://localhost:8080`

## Usage

### Authentication

1. Navigate to `http://localhost:8080/auth/login` to start OAuth flow
2. Authenticate with your Identity Provider
3. Receive a JWT token in the callback response
4. Use the JWT token for API requests:

```bash
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  http://localhost:8080/auth/profile
```

### Protected Endpoints

The service provides several protected endpoints demonstrating RBAC:

- `/api/admin` - Admin-only access
- `/api/user/data` - User and Admin access
- `/api/viewer/data` - Viewer, User, and Admin access
- `/api/data/create` - Permission-based (requires `data:create`)

## Documentation

- [Setup Guide](docs/SETUP.md) - Detailed setup instructions for each OAuth provider
- [API Documentation](docs/API.md) - Complete API reference

## Project Structure

```
.
├── auth/           # OAuth/OIDC authentication logic
├── config/         # Configuration management
├── handlers/       # HTTP request handlers
├── middleware/     # Authentication and RBAC middleware
├── models/         # Data models (User, Role, Permission)
├── utils/          # Utility functions (JWT)
├── docs/           # Documentation
├── main.go         # Application entry point
└── .env.example    # Example environment configuration
```

## Configuration

Key configuration options (via environment variables):

- `OAUTH_PROVIDER`: Identity provider (google, okta, azure)
- `OAUTH_CLIENT_ID`: OAuth client ID
- `OAUTH_CLIENT_SECRET`: OAuth client secret
- `JWT_SECRET`: Secret for signing JWT tokens
- `JWT_EXPIRATION_HOURS`: Token expiration time
- `ENABLE_RBAC`: Enable/disable RBAC (default: true)

See [.env.example](.env.example) for all options.

## Security

- JWT tokens are signed using HS256 algorithm
- CSRF protection via state parameter in OAuth flow
- Configurable token expiration
- Role-based and permission-based access control
- Secure cookie handling (HttpOnly, SameSite)

**Production Considerations:**
- Use HTTPS for all communication
- Rotate JWT signing keys regularly
- Implement rate limiting
- Use secure, random JWT secrets
- Enable secure cookies

## License

MIT

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
