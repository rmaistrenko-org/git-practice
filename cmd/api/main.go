package main

import (
	"example.com/m/internal/user"
	"example.com/m/pkg/router"
	"log"
	"net/http"
)

func main() {
	// Инициализация зависимостей через Wire
	handler := user.InitializeUserHandler()

	// Настройка маршрутизатора
	r := router.SetupRouter(handler)

	// Запуск сервера
	log.Println("Remove me")
	log.Println("Server running on port 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
