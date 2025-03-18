package handlers

import (
	"encoding/json"
	"fmt"
	"go-auth-app/database"
	"go-auth-app/middleware"
	"go-auth-app/repository"
	"go-auth-app/utils"
	"net/http"
	"strconv"
)



// GetUserDetails retrieves the authenticated user's details
func GetUserDetails(w http.ResponseWriter, r *http.Request) {
	// Get user ID from middleware
	userID := r.Context().Value(middleware.UserIDKey).(int)

	// Fetch user from database
	userRepo := repository.UserRepository{DB: database.DB}
	user, err := userRepo.GetUserByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Return user details (without password)
	response := UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		IsDeleted: user.IsDeleted,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// UpdateUser updates the authenticated user's details
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Get user ID from JWT middleware
	userID := r.Context().Value(middleware.UserIDKey).(int)

	// Parse request body
	var updatedData struct {
		Name *string `json:"name"`
	}
	json.NewDecoder(r.Body).Decode(&updatedData)

	// Ensure name is provided
	if updatedData.Name == nil || *updatedData.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	// Fetch user from DB
	userRepo := repository.UserRepository{DB: database.DB}
	user, err := userRepo.GetUserByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Update name
	user.Name = *updatedData.Name

	// Save changes to DB
	err = userRepo.UpdateUser(user)
	if err != nil {
		http.Error(w, "Failed to update name", http.StatusInternalServerError)
		return
	}

	// Return updated user details
	response := UserResponse{
		ID:   user.ID,
		Name: user.Name,
		Email: user.Email,
		IsDeleted: user.IsDeleted,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	// 🔹 Extract userID safely from context
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	fmt.Println("🗑️ DeleteUser: Request to delete user ID:", userID)

	// Call repository to mark user as deleted
	userRepo := repository.UserRepository{DB: database.DB}
	err := userRepo.SoftDeleteUser(userID)

	if err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	fmt.Println("✅ DeleteUser: User ID marked as deleted successfully", userID)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully"})
}


type UserListResponse struct {
	Users      []UserResponse `json:"users"`
	TotalUsers int           `json:"total_users"`
	Page       int           `json:"page"`
	Limit      int           `json:"limit"`
}

// GetAllUsers retrieves all users with pagination
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	// Parse pagination query params (default: page=1, limit=10)
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	// Fetch users from the database
	userRepo := repository.UserRepository{DB: database.DB}
	users, totalUsers, err := userRepo.GetUsersWithPagination(limit, offset)
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}

	// Convert users to response format
	var userResponses []UserResponse
	for _, user := range users {
		userResponses = append(userResponses, UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		})
	}

	// Return paginated user list
	response := UserListResponse{
		Users:      userResponses,
		TotalUsers: totalUsers,
		Page:       page,
		Limit:      limit,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}



type ResetPasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// ResetPassword allows a user to change their password
func ResetPassword(w http.ResponseWriter, r *http.Request) {
	// Get user ID from JWT middleware
	userID := r.Context().Value(middleware.UserIDKey).(int)

	// Parse request body
	var req ResetPasswordRequest
	json.NewDecoder(r.Body).Decode(&req)

	// Validate input
	if req.OldPassword == "" || req.NewPassword == "" {
		http.Error(w, "Both old and new passwords are required", http.StatusBadRequest)
		return
	}
	if len(req.NewPassword) < 6 {
		http.Error(w, "Password must be at least 6 characters long", http.StatusBadRequest)
		return
	}

	// Fetch only the hashed password
	userRepo := repository.UserRepository{DB: database.DB}
	hashedPassword, err := userRepo.GetUserPasswordByID(userID)
	if err != nil {
		http.Error(w, "User not found or password retrieval failed", http.StatusNotFound)
		return
	}

	fmt.Println("🔑 Stored Password Hash:", hashedPassword) // Debugging

	// Verify old password
	if !utils.CheckPasswordHash(req.OldPassword, hashedPassword) {
		http.Error(w, "Incorrect old password", http.StatusUnauthorized)
		return
	}

	// Hash new password
	newHashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		http.Error(w, "Failed to hash new password", http.StatusInternalServerError)
		return
	}

	// Update password in DB
	err = userRepo.UpdateUserPassword(userID, newHashedPassword)
	if err != nil {
		http.Error(w, "Failed to update password", http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Password updated successfully",
	})
}
