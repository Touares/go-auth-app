package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/DATA-DOG/go-txdb" // Transactional DB for testing
	_ "github.com/lib/pq"         // PostgreSQL driver
)

// Global DB instance
var DB *sql.DB

// Register the transactional DB for testing (runs once)
func init() {
	databaseURL := os.Getenv("DATABASE_URL")

	// 🔹 Ensure `sslmode=disable` is included for tests
	if !strings.Contains(databaseURL, "sslmode") {
		databaseURL += " sslmode=disable"
	}

	txdb.Register("txdb", "postgres", databaseURL)
}

// ConnectDB initializes the database and creates it if missing
func ConnectDB() {
	// Use transactional test DB if running tests
	if os.Getenv("TEST_MODE") == "true" {
		var err error
		DB, err = sql.Open("txdb", "test_db") // Use txdb for testing
		if err != nil {
			log.Fatal("❌ Failed to start test DB:", err)
		}
		fmt.Println("✅ Test database (txdb) initialized successfully!")
		return
	}

	// ✅ Use DATABASE_URL if set (recommended for Docker)
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL != "" {
		var err error
		DB, err = sql.Open("postgres", dbURL)
		if err != nil {
			log.Fatal("❌ Failed to connect to database using DATABASE_URL:", err)
		}
	} else {
		// ⛔ Fallback to manually constructing the connection string
		host := os.Getenv("DB_HOST")
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		dbName := os.Getenv("DB_NAME")
		port := os.Getenv("DB_PORT")

		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbName, port)
		fmt.Println("⚠️ DATABASE_URL not set, using manually constructed DSN:", dsn)

		var err error
		DB, err = sql.Open("postgres", dsn)
		if err != nil {
			log.Fatal("❌ Failed to connect to database using DSN:", err)
		}
	}

	// ✅ Ping the database to confirm connection
	err := DB.Ping()
	if err != nil {
		log.Fatal("❌ Database connection failed:", err)
	}

	fmt.Println("✅ Database connected successfully!")
}

