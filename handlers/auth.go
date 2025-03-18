package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-auth-app/database"
	"go-auth-app/models"
	"go-auth-app/repository"
	"go-auth-app/utils"
	"net/http"
	"regexp"
	"strings"
)

// UserResponse struct (without password)
type UserResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Password string `json:"-"`
	IsDeleted bool `json:"is_deleted"`

}

// Validate user input
func validateUserInput(user models.User) error {
	// Trim spaces from input
	user.Name = strings.TrimSpace(user.Name)
	user.Email = strings.TrimSpace(user.Email)
	user.Password = strings.TrimSpace(user.Password)

	// Validate name (at least 3 characters)
	if len(user.Name) < 3 {
		return errors.New("name must be at least 3 characters long")
	}

	// Validate email format
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	if matched, _ := regexp.MatchString(emailRegex, user.Email); !matched {
		return errors.New("invalid email format")
	}

	// Validate password (at least 6 characters)
	if len(user.Password) < 6 {
		return errors.New("password must be at least 6 characters long")
	}

	return nil
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)
	// Validate input
	if err := validateUserInput(user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Hash the password
	hashedPassword, _ := utils.HashPassword(user.Password)
	user.Password = hashedPassword

	// Create user repository
	userRepo := repository.UserRepository{DB: database.DB}
	err := userRepo.CreateUser(&user) // Pass user as pointer

	// Handle errors
	if err != nil {
		if err.Error() == "email already registered" {
			http.Error(w, "Email is already in use", http.StatusConflict)
			return
		}
		fmt.Println("❌ SQL Error in CreateUser:", err) // 🛑 Debug SQL errors
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	// Create response (without password)
	response := UserResponse{
		ID:    user.ID, // Now correctly retrieved
		Name:  user.Name,
		Email: user.Email,
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}


type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse struct
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// LoginUser handles user authentication and token issuance
func LoginUser(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	json.NewDecoder(r.Body).Decode(&req)

	// Validate input
	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	// Fetch user from DB
	userRepo := repository.UserRepository{DB: database.DB}
	user, err := userRepo.GetUserByEmail(req.Email)
	if err != nil || !utils.CheckPasswordHash(req.Password, user.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// ❌ Prevent login if the user is deactivated
	if user.IsDeleted {
		http.Error(w, "Account is deactivated. Contact support.", http.StatusForbidden)
		return
	}

	// Generate access & refresh tokens
	accessToken, err := utils.GenerateAccessToken(user.ID)
	if err != nil {
		http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		http.Error(w, "Failed to generate refresh token", http.StatusInternalServerError)
		return
	}

	// Send tokens to client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}



type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// RefreshToken generates a new access token using a valid refresh token
func RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	json.NewDecoder(r.Body).Decode(&req)

	// Validate the refresh token
	userID, err := utils.ValidateToken(req.RefreshToken, true)
	if err != nil {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	// Generate new access token
	accessToken, err := utils.GenerateAccessToken(userID)
	if err != nil {
		http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
		return
	}

	// Return new access token
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"access_token": accessToken,
	})
}
