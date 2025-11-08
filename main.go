package main

import (
	"database/sql"
	"log"
	"net/http"
	"user_management_service/cofig"
	"user_management_service/handlers"
	"user_management_service/middleware"
	"user_management_service/repository/repositoryImpl"
	"user_management_service/services/serviceImpl"

	_ "github.com/lib/pq"

	"github.com/gorilla/mux"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Connect to database
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	log.Printf(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}
	log.Println("Connected to database successfully")

	// Initialize repositories
	userRepo := repositoryImpl.NewUserRepository(db)
	sessionRepo := repositoryImpl.NewSessionRepository(db)
	permissionRepo := repositoryImpl.NewPermissionRepository(db)
	roleRepo := repositoryImpl.NewRoleRepository(db, permissionRepo)

	// Initialize services
	userService := serviceImpl.NewUserService(userRepo)
	authService := serviceImpl.NewAuthService(userRepo, sessionRepo, roleRepo, permissionRepo, cfg.JWTSecret, cfg.AccessTokenDuration, cfg.RefreshTokenDuration, cfg.BCryptCost)
	roleService := serviceImpl.NewRoleService(roleRepo, permissionRepo)
	permissionService := serviceImpl.NewPermissionService(permissionRepo)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	authHandler := handlers.NewAuthHandler(authService)
	roleHandler := handlers.NewRoleHandler(roleService)
	permissionHandler := handlers.NewPermissionHandler(permissionService)

	// Setup middleware
	authMiddleware := middleware.NewAuthMiddleware(authService)

	// Setup routes - All routes under /authapi/*
	r := mux.NewRouter()
	api := r.PathPrefix("/authapi").Subrouter()

	// Public routes (no authentication required)
	api.HandleFunc("/register", authHandler.Register).Methods("POST")
	api.HandleFunc("/login", authHandler.Login).Methods("POST")
	api.HandleFunc("/refresh", authHandler.RefreshToken).Methods("POST")
	api.HandleFunc("/health", healthCheck).Methods("GET")

	// Protected routes (authentication required)
	api.Handle("/logout", authMiddleware.Authenticate(http.HandlerFunc(authHandler.Logout))).Methods("POST")
	api.Handle("/introspect", authMiddleware.Authenticate(http.HandlerFunc(authHandler.Introspect))).Methods("GET")

	// User management protected routes
	api.Handle("/users", authMiddleware.Authenticate(http.HandlerFunc(userHandler.GetAllUsers))).Methods("GET")
	api.Handle("/users/username/{username:[a-zA-Z0-9._-]+}", authMiddleware.Authenticate(http.HandlerFunc(userHandler.GetUserByUsername))).Methods("GET")
	api.Handle("/users/email/{email:[a-zA-Z0-9._%+-@]+}", authMiddleware.Authenticate(http.HandlerFunc(userHandler.GetUserByEmail))).Methods("GET")
	api.Handle("/users/id/{id:[0-9]+}", authMiddleware.Authenticate(http.HandlerFunc(userHandler.GetUserByUserID))).Methods("GET")
	api.Handle("/users/{id:[0-9]+}", authMiddleware.Authenticate(http.HandlerFunc(userHandler.UpdateUser))).Methods("PUT")
	api.Handle("/users/{id:[0-9]+}/deactivate", authMiddleware.Authenticate(http.HandlerFunc(userHandler.DeactivateUser))).Methods("PUT")
	api.Handle("/users/{id:[0-9]+}/toggle", authMiddleware.Authenticate(http.HandlerFunc(userHandler.ToggleUserStatus))).Methods("PUT")

	// Role management protected routes
	api.Handle("/roles", authMiddleware.Authenticate(http.HandlerFunc(roleHandler.GetAllRoles))).Methods("GET")
	api.Handle("/roles", authMiddleware.Authenticate(http.HandlerFunc(roleHandler.CreateRole))).Methods("POST")
	api.Handle("/roles/{id:[0-9]+}", authMiddleware.Authenticate(http.HandlerFunc(roleHandler.UpdateRole))).Methods("PUT")

	// Permission management protected routes
	api.Handle("/permissions", authMiddleware.Authenticate(http.HandlerFunc(permissionHandler.GetAllPermissions))).Methods("GET")

	// Start server
	cors := config.CorsConfig{cfg.AllowedOrigins}
	log.Printf("Server starting on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, cors.WithCORS(r)))
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "healthy", "service": "user-management"}`))
}
