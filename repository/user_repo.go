package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"go-auth-app/models"
)

type UserRepository struct {
	DB *sql.DB
}

// CreateUser inserts a new user, but first checks if the email already exists
func (repo *UserRepository) CreateUser(user *models.User) error {
	// Check if email already exists
	var exists bool
	queryCheck := `SELECT EXISTS (SELECT 1 FROM users WHERE email = $1)`
	err := repo.DB.QueryRow(queryCheck, user.Email).Scan(&exists)
	if err != nil {
		return err
	}

	// If email exists, return an error
	if exists {
		return errors.New("email already registered")
	}

	// Insert new user if email does not exist
	queryInsert := `INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id`
	err = repo.DB.QueryRow(queryInsert, user.Name, user.Email, user.Password).Scan(&user.ID)
	if err != nil {
		fmt.Println("❌ SQL Error in CreateUser:", err) // 🛑 Debug SQL errors
		return err
	}

	return nil
}




// UserRepository handles database operations for users


// GetUserByEmail fetches a user by email
func (repo *UserRepository) GetUserByID(userID int) (models.User, error) {
	var user models.User
	query := `SELECT id, name, email, is_deleted FROM users WHERE id = $1`
	err := repo.DB.QueryRow(query, userID).Scan(&user.ID, &user.Name, &user.Email, &user.IsDeleted)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (repo *UserRepository) GetUserPasswordByID(userID int) (string, error) {
	var passwordHash string
	query := `SELECT password FROM users WHERE id = $1`
	err := repo.DB.QueryRow(query, userID).Scan(&passwordHash)

	if err != nil {
		return "", err
	}

	return passwordHash, nil
}


// GetUserByEmail fetches a user by email (for authentication)
func (repo *UserRepository) GetUserByEmail(email string) (models.User, error) {
	var user models.User
	query := `SELECT id, name, email, password, is_deleted FROM users WHERE email = $1`
	err := repo.DB.QueryRow(query, email).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.IsDeleted)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}


func (repo *UserRepository) UpdateUser(user models.User) error {
	query := `UPDATE users SET name = $1 WHERE id = $2`
	_, err := repo.DB.Exec(query, user.Name, user.ID)
	return err
}


func (repo *UserRepository) SoftDeleteUser(userID int) error {
	query := `UPDATE users SET is_deleted = TRUE WHERE id = $1`
	_, err := repo.DB.Exec(query, userID)
	return err
}


// GetUsersWithPagination retrieves users with pagination
func (repo *UserRepository) GetUsersWithPagination(limit, offset int) ([]models.User, int, error) {
	var users []models.User
	var totalUsers int

	// Query to get total users count (excluding deleted users)
	countQuery := `SELECT COUNT(*) FROM users WHERE is_deleted = FALSE`
	err := repo.DB.QueryRow(countQuery).Scan(&totalUsers)
	if err != nil {
		return nil, 0, err
	}

	// Query to get paginated users (excluding deleted users)
	query := `SELECT id, name, email FROM users WHERE is_deleted = FALSE ORDER BY id ASC LIMIT $1 OFFSET $2`
	rows, err := repo.DB.Query(query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// Parse users
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	return users, totalUsers, nil
}


func (repo *UserRepository) UpdateUserPassword(userID int, newPassword string) error {
	query := `UPDATE users SET password = $1 WHERE id = $2`
	_, err := repo.DB.Exec(query, newPassword, userID)
	return err
}
