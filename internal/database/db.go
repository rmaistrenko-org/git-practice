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

// ProvideDatabase initializes and returns the database connection
func ProvideDatabase(cfg *config.Config) *sql.DB {
	if DB == nil {
		ConnectToDB(cfg)
	}
	return DB
}

// ConnectToDB establishes a connection to the MySQL database with retry logic
func ConnectToDB(cfg *config.Config) {
	// Формування DSN з логуванням
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)
	log.Printf("Connecting to database with DSN: %s", dsn)

	var err error
	for i := 0; i < 10; i++ { // 10 спроб із паузою 2 секунди
		DB, err = sql.Open("mysql", dsn)
		if err != nil {
			log.Printf("Failed to open DB connection: %v. Retrying...", err)
			time.Sleep(2 * time.Second)
			continue
		}

		// Перевірка доступності бази даних
		if err = DB.Ping(); err == nil {
			log.Println("Connected to the database successfully")

			// Налаштування пулу з'єднань
			DB.SetMaxOpenConns(10)                  // Максимальна кількість одночасних з'єднань
			DB.SetMaxIdleConns(5)                   // Максимальна кількість "неактивних" з'єднань
			DB.SetConnMaxLifetime(30 * time.Minute) // Максимальний час життя з'єднання
			return
		}

		log.Printf("Database connection failed: %v. Retrying in 2 seconds...", err)
		time.Sleep(2 * time.Second)
	}

	// Якщо не вдалося підключитися після 10 спроб
	if err != nil {
		log.Fatalf("Could not connect to the database after multiple attempts: %v", err)
	}
}

// CloseDatabase закриває підключення до бази даних
func CloseDatabase() {
	if DB != nil {
		if err := DB.Close(); err != nil {
			log.Printf("Failed to close database connection: %v", err)
		} else {
			log.Println("Database connection closed successfully")
		}
	}
}
