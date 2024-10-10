package main

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"log"
	"net/http"
)

func main() {
	cfg := mysql.Config{
		User:   "root",
		Passwd: "1111",
		Net:    "tcp",
		Addr:   "localhost:3306",
		DBName: "go_crud_api",
	}

	// MySQL Connection: We establish a connection to the MySQL database
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Routing: The Gorilla Mux router is used to define routes and handle API requests.
	router := mux.NewRouter()

	// Define API routes
	router.HandleFunc("/user", createUser).Methods("POST")        // Create a new user
	router.HandleFunc("/users", getUsers).Methods("GET")          // Fetch all users
	router.HandleFunc("/user/{id}", getUser).Methods("GET")       // Fetch a user by ID
	router.HandleFunc("/user/{id}", updateUser).Methods("PUT")    // Update a user by ID
	router.HandleFunc("/user/{id}", deleteUser).Methods("DELETE") // Delete a user by ID

	// Start server on port 8000
	log.Fatal(http.ListenAndServe(":8000", router))
}
