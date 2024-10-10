package router

import (
	"example.com/m/internal/user"
	"github.com/gorilla/mux"
)

// SetupRouter initializes the Gorilla Mux router and defines routes
func SetupRouter() *mux.Router {
	router := mux.NewRouter()

	// Define API routes
	router.HandleFunc("/user", user.CreateUserHandler).Methods("POST")
	router.HandleFunc("/users", user.GetUsersHandler).Methods("GET")
	router.HandleFunc("/user/{id}", user.GetUserHandler).Methods("GET")
	router.HandleFunc("/user/{id}", user.UpdateUserHandler).Methods("PUT")
	router.HandleFunc("/user/{id}", user.DeleteUserHandler).Methods("DELETE")

	return router
}
