package user

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// CreateUser inserts a new user into the database
func CreateUser(db *sql.DB, user *User) error {
	query := "INSERT INTO users (name, email, created_at) VALUES (?, ?, ?)"
	createdAt := time.Now().Format("2006-01-02 15:04:05")

	result, err := db.Exec(query, user.Name, user.Email, createdAt)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = int(id)
	user.CreatedAt = createdAt

	return nil
}

// GetUsers retrieves all users from the database
func GetUsers(db *sql.DB) ([]User, error) {
	query := "SELECT id, name, email, created_at FROM users"
	rows, err := db.Query(query)
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
func GetUserByID(db *sql.DB, id int) (*User, error) {
	query := "SELECT id, name, email, created_at FROM users WHERE id = ?"
	var user User
	err := db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates an existing user in the database
func UpdateUser(db *sql.DB, id int, user *User) error {
	var updates []string
	var args []interface{}

	if user.Name != "" {
		updates = append(updates, "name = ?")
		args = append(args, user.Name)
	}
	if user.Email != "" {
		updates = append(updates, "email = ?")
		args = append(args, user.Email)
	}

	if len(updates) == 0 {
		return fmt.Errorf("no fields to update")
	}

	args = append(args, id)
	query := fmt.Sprintf("UPDATE users SET %s WHERE id = ?", strings.Join(updates, ", "))

	_, err := db.Exec(query, args...)
	return err
}

// DeleteUser deletes a user by ID
func DeleteUser(db *sql.DB, id int) error {
	query := "DELETE FROM users WHERE id = ?"
	_, err := db.Exec(query, id)
	return err
}
