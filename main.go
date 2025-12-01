package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Hilina-t/microservice-authenticator/auth"
	"github.com/Hilina-t/microservice-authenticator/config"
	"github.com/Hilina-t/microservice-authenticator/handlers"
	"github.com/Hilina-t/microservice-authenticator/middleware"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Printf("Starting Identity and Authorization Gateway (IAG)")
	log.Printf("OAuth Provider: %s", cfg.OAuthProvider)
	log.Printf("RBAC Enabled: %v", cfg.EnableRBAC)

	// Initialize services
	oauthService := auth.NewOAuthService(cfg)
	authHandler := handlers.NewAuthHandler(cfg, oauthService)
	protectedHandler := handlers.NewProtectedHandler()

	// Setup routes
	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"message": "Identity and Authorization Gateway (IAG)", "version": "1.0.0"}`)
	})
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status": "healthy"}`)
	})

	// Authentication routes
	mux.HandleFunc("/auth/login", authHandler.Login)
	mux.HandleFunc("/auth/callback", authHandler.Callback)

	// Protected routes (require authentication)
	mux.Handle("/auth/profile", middleware.AuthMiddleware(cfg)(http.HandlerFunc(authHandler.Profile)))
	mux.Handle("/auth/logout", middleware.AuthMiddleware(cfg)(http.HandlerFunc(authHandler.Logout)))

	// RBAC protected routes
	if cfg.EnableRBAC {
		// Admin-only endpoint
		mux.Handle("/api/admin",
			middleware.AuthMiddleware(cfg)(
				middleware.RequireRole("admin")(
					http.HandlerFunc(protectedHandler.AdminOnly),
				),
			),
		)

		// User endpoint (requires user or admin role)
		mux.Handle("/api/user/data",
			middleware.AuthMiddleware(cfg)(
				middleware.RequireRole("user", "admin")(
					http.HandlerFunc(protectedHandler.UserData),
				),
			),
		)

		// Viewer endpoint (requires viewer, user, or admin role)
		mux.Handle("/api/viewer/data",
			middleware.AuthMiddleware(cfg)(
				middleware.RequireRole("viewer", "user", "admin")(
					http.HandlerFunc(protectedHandler.ViewerData),
				),
			),
		)

		// Permission-based endpoint example
		mux.Handle("/api/data/create",
			middleware.AuthMiddleware(cfg)(
				middleware.RequirePermission("data", "create")(
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						user, _ := middleware.GetUserFromContext(r.Context())
						w.Header().Set("Content-Type", "application/json")
						fmt.Fprintf(w, `{"message": "Data creation allowed", "user": "%s"}`, user.Email)
					}),
				),
			),
		)
	}

	// Start server
	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Server starting on %s", addr)
	log.Printf("OAuth Login URL: http://localhost:%s/auth/login", cfg.ServerPort)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
