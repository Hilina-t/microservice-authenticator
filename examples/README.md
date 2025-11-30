# Examples

This directory contains example code demonstrating how to use the Identity and Authorization Gateway (IAG).

## Client Example

`client_example.go` demonstrates how to interact with the IAG service from a Go client application.

### Running the Example

1. Start the IAG service:
```bash
cd ..
go run main.go
```

2. In another terminal, run the example client:
```bash
cd examples
go run client_example.go
```

3. Follow the OAuth flow by visiting the login URL shown in the output.

4. After authentication, set the JWT token and run again:
```bash
export IAG_TOKEN='your-jwt-token-here'
go run client_example.go
```

### Example Output

```
IAG Client Example
==================

1. Testing health endpoint...
   ✓ Health check passed
   Response: {"status":"healthy"}

2. Testing root endpoint...
   ✓ Root endpoint passed
   Response: {"message":"Identity and Authorization Gateway (IAG)","version":"1.0.0"}

3. Testing protected endpoint without authentication...
   ✓ Correctly rejected (401 Unauthorized)
   Response: Authorization header required

4. OAuth Flow:
   To complete the OAuth flow, open your browser to:
   http://localhost:8080/auth/login
   ...

5. Testing with JWT token...
   Testing Profile endpoint...
   ✓ Success (200 OK)
   Response:
   {
     "id": "123456",
     "email": "user@example.com",
     "name": "John Doe",
     ...
   }
```

## Integration with Other Services

### Using IAG as an API Gateway

The IAG can be deployed as a gateway for multiple backend services:

```
Client → IAG (Auth + RBAC) → Backend Service A
                           → Backend Service B
                           → Backend Service C
```

### Example Architecture

1. **Client Application**: Web/Mobile app
2. **IAG Service**: Handles authentication and authorization
3. **Backend Services**: Microservices that trust IAG's JWT tokens

### Backend Service Example (Go)

```go
package main

import (
    "fmt"
    "net/http"
    "strings"
    
    "github.com/golang-jwt/jwt/v5"
)

func validateJWT(tokenString, secret string) (*jwt.Token, error) {
    return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        return []byte(secret), nil
    })
}

func protectedHandler(w http.ResponseWriter, r *http.Request) {
    authHeader := r.Header.Get("Authorization")
    if authHeader == "" {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    
    tokenString := strings.TrimPrefix(authHeader, "Bearer ")
    token, err := validateJWT(tokenString, "your-secret-key")
    
    if err != nil || !token.Valid {
        http.Error(w, "Invalid token", http.StatusUnauthorized)
        return
    }
    
    // Extract claims and process request
    if claims, ok := token.Claims.(jwt.MapClaims); ok {
        userID := claims["user_id"].(string)
        roles := claims["roles"].([]interface{})
        
        // Check roles and permissions
        // Process business logic
        fmt.Fprintf(w, "Hello, user %s with roles %v", userID, roles)
    }
}

func main() {
    http.HandleFunc("/api/resource", protectedHandler)
    http.ListenAndServe(":8081", nil)
}
```

### Frontend Integration Example (JavaScript)

```javascript
// Login flow
async function login() {
    // Redirect to IAG login
    window.location.href = 'http://localhost:8080/auth/login';
}

// Handle OAuth callback
async function handleCallback() {
    // Get JWT token from callback response
    const response = await fetch(window.location.href);
    const data = await response.json();
    
    // Store token
    localStorage.setItem('jwt_token', data.token);
    localStorage.setItem('user', JSON.stringify(data.user));
}

// Make authenticated requests
async function fetchProtectedResource() {
    const token = localStorage.getItem('jwt_token');
    
    const response = await fetch('http://localhost:8080/api/user/data', {
        headers: {
            'Authorization': `Bearer ${token}`
        }
    });
    
    if (response.status === 401) {
        // Token expired or invalid, redirect to login
        login();
        return;
    }
    
    if (response.status === 403) {
        // Insufficient permissions
        console.error('Access denied');
        return;
    }
    
    const data = await response.json();
    return data;
}

// Check user roles
function hasRole(role) {
    const user = JSON.parse(localStorage.getItem('user'));
    return user && user.roles.includes(role);
}

// Conditional UI rendering based on roles
if (hasRole('admin')) {
    showAdminPanel();
}
```

## Proxy Configuration

### Nginx Reverse Proxy

```nginx
server {
    listen 80;
    server_name example.com;

    # OAuth endpoints
    location /auth/ {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    # Protected API endpoints
    location /api/ {
        proxy_pass http://localhost:8080;
        proxy_set_header Authorization $http_authorization;
        proxy_pass_header Authorization;
    }

    # Backend services (after IAG validation)
    location /service-a/ {
        # First validate with IAG
        auth_request /auth/validate;
        
        # Then proxy to backend
        proxy_pass http://service-a:8081/;
    }
}
```

## Testing Examples

### cURL Examples

```bash
# Login (manual OAuth flow)
curl http://localhost:8080/auth/login

# After getting token from callback
export TOKEN="your-jwt-token"

# Get profile
curl -H "Authorization: Bearer $TOKEN" \
     http://localhost:8080/auth/profile

# Access protected resource
curl -H "Authorization: Bearer $TOKEN" \
     http://localhost:8080/api/user/data

# Test RBAC
curl -H "Authorization: Bearer $TOKEN" \
     http://localhost:8080/api/admin

# Logout
curl -H "Authorization: Bearer $TOKEN" \
     http://localhost:8080/auth/logout
```

### Python Example

```python
import requests
import webbrowser

BASE_URL = "http://localhost:8080"

# Step 1: Initiate OAuth flow
webbrowser.open(f"{BASE_URL}/auth/login")

# Step 2: After OAuth, you'll get a JWT token
# (In a real app, handle the callback automatically)
token = input("Enter JWT token: ")

# Step 3: Make authenticated requests
headers = {"Authorization": f"Bearer {token}"}

# Get profile
response = requests.get(f"{BASE_URL}/auth/profile", headers=headers)
print(f"Profile: {response.json()}")

# Access protected resources
response = requests.get(f"{BASE_URL}/api/user/data", headers=headers)
if response.status_code == 200:
    print(f"User data: {response.json()}")
elif response.status_code == 403:
    print("Access denied: Insufficient permissions")
```

## Best Practices

1. **Token Storage**: Store JWT tokens securely (httpOnly cookies or secure storage)
2. **Token Refresh**: Implement token refresh mechanism for long-lived sessions
3. **Error Handling**: Handle 401 (unauthorized) and 403 (forbidden) appropriately
4. **HTTPS**: Always use HTTPS in production
5. **Token Validation**: Backend services should validate JWT tokens
6. **Role Checking**: Check roles/permissions before rendering UI or processing requests

## Additional Resources

- [API Documentation](../docs/API.md)
- [Setup Guide](../docs/SETUP.md)
- [Testing Guide](../docs/TESTING.md)
