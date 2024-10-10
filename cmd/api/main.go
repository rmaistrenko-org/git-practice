package main

import (
	"example.com/m/config"
	"example.com/m/internal/database"
	"example.com/m/pkg/router"
	"log"
	"net/http"
)

func main() {
	// Загрузка конфигурации
	cfg := config.LoadConfig()

	// Инициализация подключения к базе данных
	database.ConnectToDB(cfg)

	// Проверяем подключение к базе данных
	if database.DB == nil {
		log.Fatal("Database is not initialized")
	}

	// Настройка маршрутизатора
	r := router.SetupRouter()

	// Запуск сервера
	log.Printf("Starting server on port %s...", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(cfg.ServerPort, r))
}
