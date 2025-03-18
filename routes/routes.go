package routes

import (
	"go-auth-app/handlers"
	"go-auth-app/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func SetupRoutes() {
	r := mux.NewRouter()

	// Public Routes (No Authentication Required)
	r.HandleFunc("/register", handlers.RegisterUser).Methods("POST")
	r.HandleFunc("/login", handlers.LoginUser).Methods("POST")
	r.HandleFunc("/refresh", handlers.RefreshToken).Methods("POST")

	// Protected Routes (Require JWT)
	protected := r.PathPrefix("/users").Subrouter()
	protected.Use(middleware.JWTMiddleware) // Apply JWT middleware to all /users routes
	protected.HandleFunc("", handlers.GetAllUsers).Methods("GET")
	// ✅ Separate Routes for Different Actions
	protected.HandleFunc("/me", handlers.GetUserDetails).Methods("GET")      // Fetch user details
	protected.HandleFunc("/me/update", handlers.UpdateUser).Methods("PATCH")       // Update user details
	protected.HandleFunc("/me/deactivate", handlers.DeleteUser).Methods("DELETE")   // Soft delete user
	protected.HandleFunc("/me/reset-password", handlers.ResetPassword).Methods("POST")
	// Start HTTP server
	http.Handle("/", r)
}

