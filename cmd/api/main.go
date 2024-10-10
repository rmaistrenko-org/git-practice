package main

import (
	"example.com/m/config"
	"example.com/m/internal/database"
	"example.com/m/pkg/router"
	"log"
	"net/http"
)

func main() {
	// Initialize configuration
	cfg := config.LoadConfig()

	// Establish connection to the database
	database.ConnectToDB(cfg) // Ensure this is called to initialize the database

	// Check if the database connection is properly initialized
	if database.DB == nil {
		log.Fatal("Database connection failed")
	}

	// Initialize the router for API endpoints
	r := router.SetupRouter()

	// Start the server
	log.Printf("Starting server on port %s...", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(cfg.ServerPort, r))
}
