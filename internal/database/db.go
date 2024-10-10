package database

import (
	"database/sql"
	"example.com/m/config"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
)

var DB *sql.DB

// ConnectToDB establishes a connection to the MySQL database
func ConnectToDB(cfg *config.Config) {
	// Build Data Source Name (DSN) for the database connection
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	var err error
	// Open the database connection
	DB, err = sql.Open("mysql", dsn) // Ensure that DB is correctly initialized
	if err != nil {
		log.Fatalf("Error connecting to the database: %s", err.Error())
	}

	// Verify the connection with Ping
	err = DB.Ping()
	if err != nil {
		log.Fatalf("Database ping failed: %s", err.Error())
	}

	log.Println("Connected to the database successfully")
}
