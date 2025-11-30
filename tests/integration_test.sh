#!/bin/bash

# Integration test script for IAG
# This script tests the basic functionality of the Identity and Authorization Gateway

set -e

BASE_URL="${BASE_URL:-http://localhost:8080}"

echo "Testing Identity and Authorization Gateway (IAG)"
echo "Base URL: $BASE_URL"
echo ""

# Test 1: Health check
echo "Test 1: Health Check"
response=$(curl -s -w "\n%{http_code}" "$BASE_URL/health")
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | head -n1)

if [ "$http_code" -eq 200 ]; then
    echo "✓ Health check passed"
    echo "  Response: $body"
else
    echo "✗ Health check failed (HTTP $http_code)"
    exit 1
fi
echo ""

# Test 2: Root endpoint
echo "Test 2: Root Endpoint"
response=$(curl -s -w "\n%{http_code}" "$BASE_URL/")
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | head -n1)

if [ "$http_code" -eq 200 ]; then
    echo "✓ Root endpoint passed"
    echo "  Response: $body"
else
    echo "✗ Root endpoint failed (HTTP $http_code)"
    exit 1
fi
echo ""

# Test 3: Protected endpoint without auth (should fail)
echo "Test 3: Protected Endpoint Without Auth (should return 401)"
response=$(curl -s -w "\n%{http_code}" "$BASE_URL/auth/profile")
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | head -n1)

if [ "$http_code" -eq 401 ]; then
    echo "✓ Protected endpoint correctly rejected unauthorized request"
    echo "  Response: $body"
else
    echo "✗ Protected endpoint should return 401 (got HTTP $http_code)"
    exit 1
fi
echo ""

# Test 4: Protected endpoint with invalid token (should fail)
echo "Test 4: Protected Endpoint With Invalid Token (should return 401)"
response=$(curl -s -w "\n%{http_code}" -H "Authorization: Bearer invalid-token" "$BASE_URL/auth/profile")
http_code=$(echo "$response" | tail -n1)
body=$(echo "$response" | head -n1)

if [ "$http_code" -eq 401 ]; then
    echo "✓ Protected endpoint correctly rejected invalid token"
    echo "  Response: $body"
else
    echo "✗ Protected endpoint should return 401 for invalid token (got HTTP $http_code)"
    exit 1
fi
echo ""

# Test 5: Check login endpoint redirects
echo "Test 5: Login Endpoint (should redirect to OAuth provider)"
response=$(curl -s -w "\n%{http_code}" -I "$BASE_URL/auth/login")
http_code=$(echo "$response" | tail -n1)

if [ "$http_code" -eq 307 ] || [ "$http_code" -eq 302 ]; then
    echo "✓ Login endpoint redirects correctly"
else
    echo "✗ Login endpoint should redirect (got HTTP $http_code)"
    # Note: This might fail if OAuth is not configured, which is expected
    echo "  (This is expected if OAuth credentials are not configured)"
fi
echo ""

echo "═══════════════════════════════════════════════════════════"
echo "Basic integration tests completed!"
echo ""
echo "Note: Full OAuth flow testing requires:"
echo "  1. Valid OAuth credentials configured"
echo "  2. Manual authentication with OAuth provider"
echo "  3. Capturing the JWT token from callback"
echo ""
echo "To test the full flow manually:"
echo "  1. Set up OAuth credentials in .env file"
echo "  2. Start the server: go run main.go"
echo "  3. Visit: $BASE_URL/auth/login"
echo "  4. Complete OAuth flow"
echo "  5. Use returned JWT token to test protected endpoints"
echo "═══════════════════════════════════════════════════════════"
