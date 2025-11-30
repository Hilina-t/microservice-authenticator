# Testing Guide

This guide covers how to test the Identity and Authorization Gateway (IAG).

## Unit Tests

### Running All Tests

```bash
go test ./... -v
```

### Running Specific Package Tests

```bash
# Test models
go test ./models -v

# Test utils (JWT)
go test ./utils -v
```

### Test Coverage

```bash
go test ./... -cover
```

### Generate Coverage Report

```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

## Integration Tests

### Prerequisites

1. Build and start the application:
```bash
go build -o microservice-authenticator
./microservice-authenticator
```

Or with environment variables:
```bash
# Create .env file with minimal required values
cat > .env << EOF
OAUTH_PROVIDER=google
OAUTH_CLIENT_ID=test-client-id
OAUTH_CLIENT_SECRET=test-client-secret
EOF

# Run with dotenv or export variables
export $(cat .env | xargs)
./microservice-authenticator
```

2. Run the integration test script:
```bash
./tests/integration_test.sh
```

## Manual Testing

### Test Health Endpoint

```bash
curl http://localhost:8080/health
```

Expected response:
```json
{"status":"healthy"}
```

### Test Root Endpoint

```bash
curl http://localhost:8080/
```

Expected response:
```json
{"message":"Identity and Authorization Gateway (IAG)","version":"1.0.0"}
```

### Test Protected Endpoint (Without Auth)

```bash
curl http://localhost:8080/auth/profile
```

Expected response (401 Unauthorized):
```
Authorization header required
```

### Test OAuth Login Flow

1. Start the server:
```bash
go run main.go
```

2. Open browser to:
```
http://localhost:8080/auth/login
```

3. Complete OAuth authentication with your provider

4. Copy the JWT token from the response

5. Test authenticated endpoints:

```bash
# Set your JWT token
TOKEN="your-jwt-token-here"

# Get profile
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/auth/profile

# Test admin endpoint (requires admin role)
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/admin

# Test user data endpoint (requires user or admin role)
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/user/data

# Test viewer endpoint (requires viewer, user, or admin role)
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/viewer/data

# Test permission-based endpoint
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/data/create
```

## Testing with Different Roles

To test RBAC functionality with different roles, you need to modify the role assignment logic in `auth/oauth.go`:

### Example: Role Assignment Based on Email Domain

Edit `auth/oauth.go` in the `GetUserInfo` function:

```go
// Assign roles based on email domain
if strings.HasSuffix(user.Email, "@admin.example.com") {
    user.Roles = []string{string(models.RoleAdmin)}
} else if strings.HasSuffix(user.Email, "@user.example.com") {
    user.Roles = []string{string(models.RoleUser)}
} else if strings.HasSuffix(user.Email, "@viewer.example.com") {
    user.Roles = []string{string(models.RoleViewer)}
} else {
    user.Roles = []string{string(models.RoleUser)}
}
```

### Testing Admin Role

1. Login with an admin user
2. Get JWT token
3. Test admin-only endpoint:

```bash
curl -H "Authorization: Bearer $ADMIN_TOKEN" \
  http://localhost:8080/api/admin
```

Expected: 200 OK with success response

### Testing User Role

1. Login with a regular user
2. Get JWT token
3. Test admin endpoint (should fail):

```bash
curl -H "Authorization: Bearer $USER_TOKEN" \
  http://localhost:8080/api/admin
```

Expected: 403 Forbidden

4. Test user endpoint (should succeed):

```bash
curl -H "Authorization: Bearer $USER_TOKEN" \
  http://localhost:8080/api/user/data
```

Expected: 200 OK

### Testing Viewer Role

1. Login with a viewer user
2. Get JWT token
3. Test data creation (should fail):

```bash
curl -H "Authorization: Bearer $VIEWER_TOKEN" \
  http://localhost:8080/api/data/create
```

Expected: 403 Forbidden

4. Test viewer endpoint (should succeed):

```bash
curl -H "Authorization: Bearer $VIEWER_TOKEN" \
  http://localhost:8080/api/viewer/data
```

Expected: 200 OK

## Testing JWT Token Validation

### Test with Expired Token

Generate a short-lived token (set `JWT_EXPIRATION_HOURS=0` temporarily) and wait for it to expire, then try to use it:

```bash
curl -H "Authorization: Bearer $EXPIRED_TOKEN" \
  http://localhost:8080/auth/profile
```

Expected: 401 Unauthorized with "token is expired" error

### Test with Invalid Signature

Modify a valid JWT token slightly and try to use it:

```bash
curl -H "Authorization: Bearer invalid.token.here" \
  http://localhost:8080/auth/profile
```

Expected: 401 Unauthorized

### Test with Malformed Token

```bash
curl -H "Authorization: Bearer not-a-valid-jwt" \
  http://localhost:8080/auth/profile
```

Expected: 401 Unauthorized

## Performance Testing

### Using Apache Bench

Test health endpoint:
```bash
ab -n 1000 -c 10 http://localhost:8080/health
```

Test protected endpoint with auth:
```bash
ab -n 1000 -c 10 \
  -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/auth/profile
```

### Using wrk

```bash
wrk -t4 -c100 -d30s http://localhost:8080/health
```

With authentication:
```bash
wrk -t4 -c100 -d30s \
  -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/auth/profile
```

## Docker Testing

### Build Docker Image

```bash
docker build -t microservice-authenticator .
```

### Run with Docker

```bash
docker run -p 8080:8080 \
  -e OAUTH_PROVIDER=google \
  -e OAUTH_CLIENT_ID=your-client-id \
  -e OAUTH_CLIENT_SECRET=your-client-secret \
  -e JWT_SECRET=your-secret \
  microservice-authenticator
```

### Run with Docker Compose

```bash
# Create .env file with your credentials
cp .env.example .env
# Edit .env with your values

# Start service
docker-compose up -d

# View logs
docker-compose logs -f

# Stop service
docker-compose down
```

## Automated Testing with CI/CD

Example GitHub Actions workflow:

```yaml
name: Test

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - name: Run tests
        run: go test ./... -v -cover
      - name: Build
        run: go build -v ./...
```

## Troubleshooting Tests

### Test Failures

1. **JWT Test Failures**: Ensure system time is correct
2. **Integration Test Failures**: Verify server is running on correct port
3. **OAuth Test Failures**: Check OAuth credentials and redirect URLs

### Common Issues

1. **Port Already in Use**: Change `SERVER_PORT` in .env
2. **OAuth Configuration**: Verify credentials are correct
3. **Network Issues**: Check firewall settings

## Best Practices

1. Run unit tests before every commit
2. Run integration tests before pushing to main branch
3. Keep test coverage above 80%
4. Test both success and failure scenarios
5. Test edge cases (expired tokens, invalid roles, etc.)
6. Use table-driven tests for multiple scenarios
7. Mock external dependencies when possible

## Test Data

For testing purposes, you can use these test scenarios:

1. **Valid User**: User with `user` role
2. **Admin User**: User with `admin` role
3. **Viewer User**: User with `viewer` role
4. **No Role User**: User with no assigned roles
5. **Multi-Role User**: User with multiple roles

## Security Testing

1. Test SQL injection (if using database)
2. Test XSS vulnerabilities
3. Test CSRF protection
4. Test rate limiting
5. Test token leakage
6. Test secure headers
7. Verify HTTPS in production

## Monitoring Tests

Set up health checks in production:

```bash
# Add to monitoring system
curl -f http://localhost:8080/health || exit 1
```
