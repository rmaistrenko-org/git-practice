package user

import (
	"database/sql"
	"example.com/m/internal/database"
	"fmt"
	"strings"
	"time"
)

// CreateUser inserts a new user into the database
func CreateUser(user *User) error {
	// Check if database is initialized
	if database.DB == nil {
		return fmt.Errorf("database is not initialized")
	}

	// Insert new user and capture the current timestamp for CreatedAt
	query := "INSERT INTO users (name, email, created_at) VALUES (?, ?, ?)"
	createdAt := time.Now().Format("2006-01-02 15:04:05")

	result, err := database.DB.Exec(query, user.Name, user.Email, createdAt)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	// Update the User struct with the new ID and CreatedAt timestamp
	user.ID = int(id)
	user.CreatedAt = createdAt

	return nil
}

// GetUsers retrieves all users from the database
func GetUsers() ([]User, error) {
	// Check if database is initialized
	if database.DB == nil {
		return nil, fmt.Errorf("database is not initialized")
	}

	query := "SELECT id, name, email, created_at FROM users"
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// GetUserByID retrieves a single user by ID
func GetUserByID(id int) (*User, error) {
	// Check if database is initialized
	if database.DB == nil {
		return nil, fmt.Errorf("database is not initialized")
	}

	query := "SELECT id, name, email, created_at FROM users WHERE id = ?"
	var user User
	err := database.DB.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates an existing user in the database
func UpdateUser(id int, user *User) error {
	// Check if database is initialized
	if database.DB == nil {
		return fmt.Errorf("database is not initialized")
	}

	// Create a slice to hold the columns to update
	var updates []string
	var args []interface{}

	// Check which fields are provided and build the query dynamically
	if user.Name != "" {
		updates = append(updates, "name = ?")
		args = append(args, user.Name)
	}
	if user.Email != "" {
		updates = append(updates, "email = ?")
		args = append(args, user.Email)
	}

	// If no updates, return an error
	if len(updates) == 0 {
		return fmt.Errorf("no fields to update")
	}

	// Add the user ID to the arguments (it will be used in the WHERE clause)
	args = append(args, id)

	// Build the final SQL query by joining the fields
	query := fmt.Sprintf("UPDATE users SET %s WHERE id = ?", strings.Join(updates, ", "))

	// Execute the update query
	_, err := database.DB.Exec(query, args...)
	if err != nil {
		return err
	}

	// Now fetch the updated user data from the database
	return database.DB.QueryRow("SELECT id, name, email, created_at FROM users WHERE id = ?", id).
		Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
}

// DeleteUser deletes a user by ID
func DeleteUser(id int) error {
	// Check if database is initialized
	if database.DB == nil {
		return fmt.Errorf("database is not initialized")
	}

	query := "DELETE FROM users WHERE id = ?"
	_, err := database.DB.Exec(query, id)
	return err
}
