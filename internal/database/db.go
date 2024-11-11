package database

import (
	"database/sql"
	"example.com/m/config"
	"fmt"
	"log"

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
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error connecting to the database: %s", err.Error()) // Лог ошибки подключения
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("Database ping failed: %s", err.Error()) // Лог ошибки ping
	}

	log.Println("Connected to the database successfully")
}
