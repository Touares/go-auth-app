package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-auth-app/database"
	"go-auth-app/middleware"

	// "go-auth-app/utils"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	// "github.com/gorilla/mux"
	_ "github.com/lib/pq" // Import PostgreSQL driver
	"github.com/stretchr/testify/assert"
	"go-auth-app/handlers"
)

func TestMain(m *testing.M) {
	// ✅ Set TEST_MODE=true so the app uses `txdb` for testing
	os.Setenv("TEST_MODE", "true")

	// ✅ Initialize the database (ConnectDB will use txdb)
	database.ConnectDB()

	code := m.Run() // Run all tests
	os.Exit(code)
}

func TestRegisterAndLoginUser(t *testing.T) {
	// 1️⃣ Register a User
	registerPayload := map[string]string{
		"name":     "Test User",
		"email":    "test@example.com",
		"password": "securepassword",
	}
	registerBody, _ := json.Marshal(registerPayload)


	reqRegister, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(registerBody))
	reqRegister.Header.Set("Content-Type", "application/json")
	rrRegister := httptest.NewRecorder()

	handlers.RegisterUser(rrRegister, reqRegister) // Call the actual register handler

	assert.Equal(t, http.StatusCreated, rrRegister.Code, "Expected 201 Created, got %d", rrRegister.Code)

	// 2️⃣ Log in the Same User
	loginPayload := map[string]string{
		"email":    "test@example.com",
		"password": "securepassword",
	}
	loginBody, _ := json.Marshal(loginPayload)

	reqLogin, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(loginBody))
	reqLogin.Header.Set("Content-Type", "application/json")
	rrLogin := httptest.NewRecorder()

	handlers.LoginUser(rrLogin, reqLogin) // Call the actual login handler



	assert.Equal(t, http.StatusOK, rrLogin.Code, "Expected 200 OK, got %d", rrLogin.Code)

	// Parse login response
	var loginResponse map[string]string
	json.Unmarshal(rrLogin.Body.Bytes(), &loginResponse)

	assert.NotEmpty(t, loginResponse["access_token"], "Access token should not be empty")
	assert.NotEmpty(t, loginResponse["refresh_token"], "Refresh token should not be empty")
}


// ✅ Test: Attempt to register with an email that already exists (409 Conflict)
func TestRegisterUser_ExistingEmail(t *testing.T) {
	// Ensure test DB is initialized

	// Register a user first
	registerPayload := map[string]string{
		"name":     "Test User",
		"email":    "duplicate@example.com",
		"password": "securepassword",
	}
	registerBody, _ := json.Marshal(registerPayload)

	reqRegister, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(registerBody))
	reqRegister.Header.Set("Content-Type", "application/json")
	rrRegister := httptest.NewRecorder()

	handlers.RegisterUser(rrRegister, reqRegister)
	assert.Equal(t, http.StatusCreated, rrRegister.Code, "Initial registration should succeed")

	// Attempt to register with the same email again
	reqDuplicate, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(registerBody))
	reqDuplicate.Header.Set("Content-Type", "application/json")
	rrDuplicate := httptest.NewRecorder()

	handlers.RegisterUser(rrDuplicate, reqDuplicate)
	

	assert.Equal(t, http.StatusConflict, rrDuplicate.Code, "Expected 409 Conflict for duplicate email")
}

// ✅ Test: Provide an invalid email format (400 Bad Request)
func TestRegisterUser_InvalidEmailFormat(t *testing.T) {

	registerPayload := map[string]string{
		"name":     "Invalid Email",
		"email":    "invalid-email", // ❌ Invalid format
		"password": "securepassword",
	}
	registerBody, _ := json.Marshal(registerPayload)

	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(registerBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handlers.RegisterUser(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Expected 400 Bad Request for invalid email format")
}

// ✅ Test: Ensure passwords shorter than 6 characters fail (400 Bad Request)
func TestRegisterUser_ShortPassword(t *testing.T) {

	registerPayload := map[string]string{
		"name":     "Weak Password",
		"email":    "shortpass@example.com",
		"password": "123", // ❌ Too short
	}
	registerBody, _ := json.Marshal(registerPayload)

	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(registerBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handlers.RegisterUser(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Expected 400 Bad Request for short password")
}

// ✅ Test: Attempt to register without required fields (400 Bad Request)
func TestRegisterUser_MissingFields(t *testing.T) {

	registerPayload := map[string]string{
		"name":  "", // ❌ Missing name
		"email": "valid@example.com",
	}
	registerBody, _ := json.Marshal(registerPayload)

	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(registerBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handlers.RegisterUser(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Expected 400 Bad Request for missing fields")
}



// ✅ Test: Attempt login with incorrect password (401 Unauthorized)
func TestLoginUser_InvalidPassword(t *testing.T) {
	database.ConnectDB()

	// First, register a user to test against
	registerPayload := map[string]string{
		"name":     "Test User",
		"email":    "invalidpass@example.com",
		"password": "securepassword",
	}
	registerBody, _ := json.Marshal(registerPayload)

	reqRegister, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(registerBody))
	reqRegister.Header.Set("Content-Type", "application/json")
	rrRegister := httptest.NewRecorder()
	handlers.RegisterUser(rrRegister, reqRegister)
	assert.Equal(t, http.StatusCreated, rrRegister.Code, "User registration should succeed")

	// Now, attempt login with the wrong password
	loginPayload := map[string]string{
		"email":    "invalidpass@example.com",
		"password": "wrongpassword", // ❌ Incorrect password
	}
	loginBody, _ := json.Marshal(loginPayload)

	reqLogin, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(loginBody))
	reqLogin.Header.Set("Content-Type", "application/json")
	rrLogin := httptest.NewRecorder()

	handlers.LoginUser(rrLogin, reqLogin)

	assert.Equal(t, http.StatusUnauthorized, rrLogin.Code, "Expected 401 Unauthorized for incorrect password")
}

// ✅ Test: Attempt login with an email that does not exist (401 Unauthorized)
func TestLoginUser_NonExistentEmail(t *testing.T) {
	database.ConnectDB()

	loginPayload := map[string]string{
		"email":    "doesnotexist@example.com", // ❌ Non-existent email
		"password": "securepassword",
	}
	loginBody, _ := json.Marshal(loginPayload)

	reqLogin, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(loginBody))
	reqLogin.Header.Set("Content-Type", "application/json")
	rrLogin := httptest.NewRecorder()

	handlers.LoginUser(rrLogin, reqLogin)
	fmt.Println("📝 Login Non-Existent Email Response Body:", rrLogin.Body.String())

	assert.Equal(t, http.StatusUnauthorized, rrLogin.Code, "Expected 401 Unauthorized for non-existent email")
}

// ✅ Test: Attempt login with missing email or password (400 Bad Request)
func TestLoginUser_EmptyFields(t *testing.T) {
	database.ConnectDB()

	// Case 1: Missing password
	loginPayload1 := map[string]string{
		"email": "missingpass@example.com",
	}
	loginBody1, _ := json.Marshal(loginPayload1)

	reqLogin1, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(loginBody1))
	reqLogin1.Header.Set("Content-Type", "application/json")
	rrLogin1 := httptest.NewRecorder()

	handlers.LoginUser(rrLogin1, reqLogin1)
	fmt.Println("📝 Login Missing Password Response Body:", rrLogin1.Body.String())

	assert.Equal(t, http.StatusBadRequest, rrLogin1.Code, "Expected 400 Bad Request for missing password")

	// Case 2: Missing email
	loginPayload2 := map[string]string{
		"password": "securepassword",
	}
	loginBody2, _ := json.Marshal(loginPayload2)

	reqLogin2, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(loginBody2))
	reqLogin2.Header.Set("Content-Type", "application/json")
	rrLogin2 := httptest.NewRecorder()

	handlers.LoginUser(rrLogin2, reqLogin2)
	fmt.Println("📝 Login Missing Email Response Body:", rrLogin2.Body.String())

	assert.Equal(t, http.StatusBadRequest, rrLogin2.Code, "Expected 400 Bad Request for missing email")
}


// ✅ Test: Fetch the list of users with valid authentication (200 OK)
func TestGetAllUsers_ValidRequest(t *testing.T) {
	// 🛠 Setup Test DB
	// _ = utils.SetupTestDB(t)

	// ✅ Create an authenticated test user
	_, accessToken, err := CreateAuthenticatedUser("testuser@example.com", "securepassword")
	if err != nil {
		t.Fatalf("❌ Failed to create authenticated user: %v", err)
	}

	// 🚀 Make a request to fetch all users
	req, _ := http.NewRequest("GET", "/users", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	rr := httptest.NewRecorder()
	handlers.GetAllUsers(rr, req)

	// 📝 Check Response
	assert.Equal(t, http.StatusOK, rr.Code, "Expected 200 OK for valid request")
}


// ✅ Test: Fetch users without authentication (401 Unauthorized)
func TestGetAllUsers_Unauthorized(t *testing.T) {
	// database.ConnectDB()

	// ❌ Fetch users without providing an authentication token
	reqGetUsers, _ := http.NewRequest("GET", "/users", nil)
	rrGetUsers := httptest.NewRecorder()

	middleware.JWTMiddleware(http.HandlerFunc(handlers.GetAllUsers)).ServeHTTP(rrGetUsers, reqGetUsers)

	assert.Equal(t, http.StatusUnauthorized, rrGetUsers.Code, "Expected 401 Unauthorized for missing authentication")
}

// ✅ Test: Ensure pagination parameters work correctly
func TestGetAllUsers_Pagination(t *testing.T) {
	// _ = utils.SetupTestDB(t)

	// ✅ Create authenticated test user
	_, accessToken, err := CreateAuthenticatedUser("pagetest@example.com", "securepassword")
	if err != nil {
		t.Fatalf("❌ Failed to create authenticated user: %v", err)
	}

	// 🚀 Make a paginated request
	req, _ := http.NewRequest("GET", "/users?page=1&limit=5", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	rr := httptest.NewRecorder()
	handlers.GetAllUsers(rr, req)

	// 📝 Check Response
	assert.Equal(t, http.StatusOK, rr.Code, "Expected 200 OK for pagination")
}



