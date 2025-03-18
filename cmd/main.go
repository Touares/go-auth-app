package main

import (
	"fmt"
	"go-auth-app/config"
	"go-auth-app/database"
	"go-auth-app/routes"
	"net/http"
)

func main() {
	config.LoadConfig()
	database.ConnectDB()
	routes.SetupRoutes()
	fmt.Println("Server running on port 8080...")
	http.ListenAndServe(":8080", nil)
}
