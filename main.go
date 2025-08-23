package main

import (
	"database/sql"
	"log"
	"net/http"
	"user_management_service/cofig"
	"user_management_service/handlers"
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
	roleRepo := repositoryImpl.NewRoleRepository(db)

	// Initialize services
	userService := serviceImpl.NewUserService(userRepo)
	authService := serviceImpl.NewAuthService(userRepo, sessionRepo, roleRepo, cfg.JWTSecret, cfg.TokenDuration, 12)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	authHandler := handlers.NewAuthHandler(authService)

	// Setup middleware
	//authMiddleware := handlers.NewAuthMiddleware(authService)

	// Setup routes
	r := mux.NewRouter()
	api := r.PathPrefix("/api").Subrouter()

	// Public routes
	//api.HandleFunc("/auth/register", authHandler.Register).Methods("POST")
	//api.HandleFunc("/auth/login", authHandler.Login).Methods("POST")
	api.HandleFunc("/health", healthCheck).Methods("GET")

	// Protected routes
	protected := api.PathPrefix("/user_management").Subrouter()
	//protected.Use(authMiddleware.Authenticate)

	protected.HandleFunc("/users/username/{username:[a-zA-Z0-9._-]+}", userHandler.GetUserByUsername).Methods("GET")
	protected.HandleFunc("/users/email/{email:[a-zA-Z0-9._%+-@]+}", userHandler.GetUserByEmail).Methods("GET")
	protected.HandleFunc("/users/id/{id:[0-9]+}", userHandler.GetUserByUserID).Methods("GET")
	protected.HandleFunc("/users/{id:[0-9]+}/deactivate", userHandler.DeactivateUser).Methods("PUT")
	//protected.HandleFunc("/auth/change-password", authHandler.ChangePassword).Methods("POST")
	//protected.HandleFunc("/auth/logout", authHandler.Logout).Methods("POST")

	protected.HandleFunc("/users/register", authHandler.Register).Methods("POST")
	protected.HandleFunc("/users/login", authHandler.Login).Methods("POST")

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "healthy", "service": "user-management"}`))
}
