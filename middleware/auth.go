package middleware

import (
	"context"
	"fmt"
	"go-auth-app/utils"

	// "log"
	"net/http"
	"strings"
)

// Context key for storing user ID
type contextKey string

const UserIDKey contextKey = "user_id"

// JWTMiddleware ensures that only authenticated users can access protected routes
func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		// Extract token (Format: "Bearer <token>")
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := tokenParts[1]

		// Validate JWT token
		userID, err := utils.ValidateToken(tokenString, false) // false = access token
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		fmt.Println("✅ JWTMiddleware: User ID extracted from token\n", userID)

		// Store user ID in request context
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

