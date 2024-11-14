package database

import (
	"database/sql"
	"example.com/m/config"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
)

var DB *sql.DB

func ProvideDatabase(cfg *config.Config) *sql.DB {
	ConnectToDB(cfg)
	return DB
}

// ConnectToDB establishes a connection to the MySQL database
func ConnectToDB(cfg *config.Config) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	var err error
	for i := 0; i < 10; i++ { // 10 спроб із паузою 2 секунди
		DB, err = sql.Open("mysql", dsn)
		if err == nil && DB.Ping() == nil {
			log.Println("Connected to the database successfully")
			return
		}
		log.Printf("Database connection failed: %v. Retrying in 2 seconds...", err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatalf("Could not connect to the database after multiple attempts: %v", err)
	}
}
