package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-auth-app/database"
	"go-auth-app/handlers"
	"go-auth-app/models"
	"go-auth-app/repository"
	"go-auth-app/utils"
	"log"
	"net/http"
	"net/http/httptest"
)

// 🔹 Fixture: Create an Authenticated Test User
func CreateAuthenticatedUser(email, password string) (*models.User, string, error) {
	// ✅ Hash password before storing
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, "", err
	}

	// ✅ Create test user
	user := &models.User{
		Name:     "Test User",
		Email:    email,
		Password: hashedPassword,
	}

	// ✅ Insert into database
	userRepo := repository.UserRepository{DB: database.DB}
	err = userRepo.CreateUser(user)
	if err != nil {
		return nil, "", err
	}


	// ✅ Simulate login request to obtain access token
	loginPayload := map[string]string{"email": email, "password": password}
	loginBody, _ := json.Marshal(loginPayload)

	reqLogin, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(loginBody))
	reqLogin.Header.Set("Content-Type", "application/json")
	rrLogin := httptest.NewRecorder()

	handlers.LoginUser(rrLogin, reqLogin)

	// ✅ Parse login response
	var loginResponse map[string]string
	json.Unmarshal(rrLogin.Body.Bytes(), &loginResponse)

	accessToken, exists := loginResponse["access_token"]
	if !exists {
		return user, "", fmt.Errorf("access token not received")
	}

	log.Println("✅ Access token obtained:", accessToken)

	return user, accessToken, nil
}
