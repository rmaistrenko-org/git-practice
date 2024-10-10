package main

import (
	"database/sql"
	"encoding/json"
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

func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := db.Exec("INSERT INTO users (name, email) VALUES (?, ?)", user.Name, user.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNetworkAuthenticationRequired)
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotExtended)

		//http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.ID = int(id)
	user.CreatedAt = "now" // Placeholder
	json.NewEncoder(w).Encode(user)
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	var users []User
	rows, err := db.Query("SELECT id, name, email, created_at FROM users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotImplemented)

		//http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}
	json.NewEncoder(w).Encode(users)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	var user User
	err := db.QueryRow("SELECT id, name, email, created_at FROM users WHERE id = ?", id).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	json.NewEncoder(w).Encode(user)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE users SET name = ?, email = ? WHERE id = ?", user.Name, user.Email, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.CreatedAt = "now" // Placeholder
	json.NewEncoder(w).Encode(user)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	_, err := db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
