package utils

import (
	"database/sql"
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/DATA-DOG/go-txdb"
	_ "github.com/lib/pq" // PostgreSQL driver
	"go-auth-app/database"
)

var once sync.Once

func registerTxDB() {
	once.Do(func() {
		fmt.Println("🚀 Registering txdb driver...")
		txdb.Register("txdb", "postgres", "user=your_user password=your_pass dbname=test_db sslmode=disable")
	})
}

func SetupTestDB(t *testing.T) *sql.DB {
	os.Setenv("TEST_MODE", "true") // Enable test mode

	// Ensure `txdb` is registered only once before connecting
	registerTxDB()

	// Open the test database
	testDB, err := sql.Open("txdb", "test_instance") // Must match txdb.Register()
	if err != nil {
		t.Fatalf("Failed to start test DB: %v", err)
	}

	fmt.Println("🚀 Test DB Initialized")

	// Assign test DB globally so all handlers use it
	database.DB = testDB

	// Ensure cleanup after test execution (rollback transactions)
	t.Cleanup(func() {
		fmt.Println("🧹 Closing Test DB Connection")
		testDB.Close()
	})

	return testDB
}
