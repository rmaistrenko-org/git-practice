package router

import (
	"example.com/m/internal/user"
	"github.com/gorilla/mux"
)

// SetupRouter инициализирует маршрутизатор Gorilla Mux
func SetupRouter(handler *user.Handler) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/user", handler.CreateUserHandler).Methods("POST")
	router.HandleFunc("/users", handler.GetUsersHandler).Methods("GET")
	router.HandleFunc("/user/{id}", handler.GetUserHandler).Methods("GET")
	router.HandleFunc("/user/{id}", handler.UpdateUserHandler).Methods("PUT")
	router.HandleFunc("/user/{id}", handler.DeleteUserHandler).Methods("DELETE")

	return router
}
