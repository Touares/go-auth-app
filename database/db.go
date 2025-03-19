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

// Register transactional test DB (runs once)
func init() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		fmt.Println("⚠️ DATABASE_URL not set, using default connection settings.")
	} else if !strings.Contains(databaseURL, "sslmode") {
		databaseURL += " sslmode=disable" // Ensure sslmode is set
	}

	fmt.Println("🔹 DATABASE_URL Used for txdb:", databaseURL)
	txdb.Register("txdb", "postgres", databaseURL)
}

// ConnectDB initializes the database connection
func ConnectDB() {
	// ✅ Use `txdb` for tests if `TEST_MODE` is set
	if os.Getenv("TEST_MODE") == "true" {
		var err error
		DB, err = sql.Open("txdb", "test_db")
		if err != nil {
			log.Fatal("❌ Failed to start test DB:", err)
		}
		fmt.Println("✅ Test database (txdb) initialized successfully!")
		return
	}

	// ✅ Prioritize `DATABASE_URL`
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		// Fallback to manual connection string
		host := os.Getenv("DB_HOST")
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		dbName := os.Getenv("DB_NAME")
		port := os.Getenv("DB_PORT")

		dbURL = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbName, port)
		fmt.Println("⚠️ DATABASE_URL not set, using manually constructed DSN")
	}

	// ✅ Connect to PostgreSQL
	var err error
	DB, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}

	// ✅ Ensure database connection is live
	if err := DB.Ping(); err != nil {
		log.Fatal("❌ Database connection failed:", err)
	}

	fmt.Println("✅ Database connected successfully!")
}
