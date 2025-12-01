package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Example client demonstrating how to interact with the IAG service

const (
	defaultBaseURL = "http://localhost:8080"
)

type LoginResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func main() {
	baseURL := os.Getenv("IAG_BASE_URL")
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	fmt.Println("IAG Client Example")
	fmt.Println("==================")
	fmt.Println()

	// Test 1: Health Check
	fmt.Println("1. Testing health endpoint...")
	testHealthEndpoint(baseURL)
	fmt.Println()

	// Test 2: Root Endpoint
	fmt.Println("2. Testing root endpoint...")
	testRootEndpoint(baseURL)
	fmt.Println()

	// Test 3: Protected Endpoint Without Auth
	fmt.Println("3. Testing protected endpoint without authentication...")
	testProtectedWithoutAuth(baseURL)
	fmt.Println()

	// Test 4: OAuth Flow
	fmt.Println("4. OAuth Flow:")
	fmt.Printf("   To complete the OAuth flow, open your browser to:\n")
	fmt.Printf("   %s/auth/login\n", baseURL)
	fmt.Println()
	fmt.Println("   After authentication, you'll receive a JWT token.")
	fmt.Println("   Set the token and test protected endpoints:")
	fmt.Println()
	fmt.Println("   export IAG_TOKEN='your-jwt-token-here'")
	fmt.Printf("   go run %s\n", os.Args[0])
	fmt.Println()

	// Test 5: Protected Endpoints with Token (if available)
	token := os.Getenv("IAG_TOKEN")
	if token != "" {
		fmt.Println("5. Testing with JWT token...")
		testProtectedEndpoints(baseURL, token)
	} else {
		fmt.Println("5. Skipping protected endpoint tests (no IAG_TOKEN set)")
	}
}

func testHealthEndpoint(baseURL string) {
	resp, err := http.Get(baseURL + "/health")
	if err != nil {
		fmt.Printf("   ✗ Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode == 200 {
		fmt.Printf("   ✓ Health check passed\n")
		fmt.Printf("   Response: %s\n", string(body))
	} else {
		fmt.Printf("   ✗ Health check failed (HTTP %d)\n", resp.StatusCode)
		fmt.Printf("   Response: %s\n", string(body))
	}
}

func testRootEndpoint(baseURL string) {
	resp, err := http.Get(baseURL + "/")
	if err != nil {
		fmt.Printf("   ✗ Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode == 200 {
		fmt.Printf("   ✓ Root endpoint passed\n")
		fmt.Printf("   Response: %s\n", string(body))
	} else {
		fmt.Printf("   ✗ Root endpoint failed (HTTP %d)\n", resp.StatusCode)
		fmt.Printf("   Response: %s\n", string(body))
	}
}

func testProtectedWithoutAuth(baseURL string) {
	resp, err := http.Get(baseURL + "/auth/profile")
	if err != nil {
		fmt.Printf("   ✗ Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode == 401 {
		fmt.Printf("   ✓ Correctly rejected (401 Unauthorized)\n")
		fmt.Printf("   Response: %s\n", string(body))
	} else {
		fmt.Printf("   ✗ Unexpected status code: %d (expected 401)\n", resp.StatusCode)
		fmt.Printf("   Response: %s\n", string(body))
	}
}

func testProtectedEndpoints(baseURL, token string) {
	endpoints := []struct {
		name          string
		path          string
		expectSuccess bool
	}{
		{"Profile", "/auth/profile", true},
		{"User Data", "/api/user/data", true},
		{"Viewer Data", "/api/viewer/data", true},
		{"Admin", "/api/admin", false},            // Likely to fail unless admin role
		{"Data Create", "/api/data/create", true}, // If user or admin role
	}

	for _, endpoint := range endpoints {
		fmt.Printf("\n   Testing %s endpoint...\n", endpoint.name)
		testProtectedEndpoint(baseURL+endpoint.path, token, endpoint.expectSuccess)
	}
}

func testProtectedEndpoint(url, token string, expectSuccess bool) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("   ✗ Error creating request: %v\n", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("   ✗ Error: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == 200 {
		fmt.Printf("   ✓ Success (200 OK)\n")
		var prettyJSON bytes.Buffer
		if err := json.Indent(&prettyJSON, body, "   ", "  "); err == nil {
			fmt.Printf("   Response:\n%s\n", prettyJSON.String())
		} else {
			fmt.Printf("   Response: %s\n", string(body))
		}
	} else if resp.StatusCode == 403 {
		fmt.Printf("   ⚠ Forbidden (403) - Insufficient permissions\n")
		fmt.Printf("   Response: %s\n", string(body))
		if !expectSuccess {
			fmt.Printf("   (This is expected behavior)\n")
		}
	} else if resp.StatusCode == 401 {
		fmt.Printf("   ✗ Unauthorized (401) - Token may be invalid or expired\n")
		fmt.Printf("   Response: %s\n", string(body))
	} else {
		fmt.Printf("   ✗ Unexpected status code: %d\n", resp.StatusCode)
		fmt.Printf("   Response: %s\n", string(body))
	}
}
