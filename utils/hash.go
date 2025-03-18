package utils

import (

	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a given password
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err // Convert []byte to string
}

// CheckPasswordHash verifies if the hashed password matches the plain text password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		
		return false
	}
	return true
}
