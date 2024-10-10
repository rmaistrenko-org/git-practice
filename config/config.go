package config

import (
	"os"
)

type Config struct {
	DBUser     string
	DBPassword string
	DBName     string
	DBHost     string
	DBPort     string
	ServerPort string
}

// LoadConfig loads environment variables or defaults for database and server configuration
func LoadConfig() *Config {
	config := &Config{
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", "1111"),
		DBName:     getEnv("DB_NAME", "go_crud_api"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		ServerPort: getEnv("SERVER_PORT", ":8000"),
	}

	return config
}

// getEnv returns the value of an environment variable or a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
