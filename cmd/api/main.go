package main

import (
	"example.com/m/config"
	"example.com/m/internal/database"
	"example.com/m/internal/user"
	"example.com/m/pkg/router"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Завантаження конфігурації
	cfg := config.ProvideConfig()

	// Ініціалізація підключення до бази даних
	db := database.ProvideDatabase(cfg)
	defer database.CloseDatabase() // Закриття підключення

	// Ініціалізація Handler
	service := user.NewService(db)
	handler := user.NewHandler(service)

	// Налаштування маршрутизатора
	r := router.SetupRouter(handler)

	// Канал для graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Запуск сервера
	server := &http.Server{
		Addr:    ":8000",
		Handler: r,
	}

	go func() {
		log.Println("Server is starting on port 8000...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Очікування сигналу для завершення роботи
	<-stop
	log.Println("Shutting down server...")

	if err := server.Close(); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
	log.Println("Server gracefully stopped")
}
