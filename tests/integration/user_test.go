package integration

import (
	"bytes"
	"encoding/json"
	"example.com/m/config"
	"example.com/m/internal/database"
	"example.com/m/internal/user"
	"example.com/m/pkg/router"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"
)

// TestMain will run before any other tests to set up the environment
func TestMain(m *testing.M) {
	cfg := config.LoadConfig()

	database.ConnectToDB(cfg)

	if database.DB == nil {
		log.Fatal("Database is not initialized")
	}

	code := m.Run()

	os.Exit(code)
}

// TestCreateUser tests the POST /user endpoint to create a new user
func TestCreateUser(t *testing.T) {
	r := router.SetupRouter()

	// Generate a unique email to avoid duplicates
	uniqueEmail := "test" + strconv.FormatInt(time.Now().UnixNano(), 10) + "@example.com"

	// Create a new user with a unique email
	newUser := &user.User{Name: "Test User", Email: uniqueEmail}
	jsonUser, _ := json.Marshal(newUser)

	// Create a POST request with the new user data
	req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(jsonUser))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Log the response body for debugging
	responseBody, _ := io.ReadAll(w.Body)
	log.Printf("POST Response body: %s", responseBody)

	// Check if the status code is 201 Created
	if status := w.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	// Decode the response body to check if the user was created correctly
	var createdUser user.User
	err := json.NewDecoder(bytes.NewReader(responseBody)).Decode(&createdUser)
	if err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}

	// Ensure the returned user has the correct name and email
	if createdUser.Name != newUser.Name {
		t.Errorf("Expected Name to be '%v'. Got '%v'", newUser.Name, createdUser.Name)
	}

	if createdUser.Email != newUser.Email {
		t.Errorf("Expected Email to be '%v'. Got '%v'", uniqueEmail, createdUser.Email)
	}

	// Check that the ID and CreatedAt fields are not zero/empty
	if createdUser.ID == 0 {
		t.Errorf("Expected user ID to be set. Got '%v'", createdUser.ID)
	}

	if createdUser.CreatedAt == "" {
		t.Errorf("Expected CreatedAt to be set. Got '%v'", createdUser.CreatedAt)
	}

	// Check that user was added to the database
	retrievedUser, err := user.GetUserByID(createdUser.ID)
	if err != nil || retrievedUser == nil {
		t.Errorf("User was not found in the database after creation")
	}
}

// TestGetUsers tests the GET /users endpoint to fetch all users
func TestGetUsers(t *testing.T) {
	r := router.SetupRouter()

	// Create a GET request to fetch all users
	req, _ := http.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Check if the status code is 200 OK
	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Decode the response body to check if it returns users
	var users []user.User
	err := json.NewDecoder(w.Body).Decode(&users)
	if err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}

	// Optionally check that users list is not empty (if there should be some users)
	if len(users) == 0 {
		t.Errorf("Expected users list to not be empty") // Убрано %v
	}
}

// TestGetUserByID tests the GET /user/{id} endpoint to retrieve a user by ID
func TestGetUserByID(t *testing.T) {
	r := router.SetupRouter()

	// Generate a unique email to avoid duplicates
	uniqueEmail := "test" + strconv.FormatInt(time.Now().UnixNano(), 10) + "@example.com"

	// First, create a new user with a unique email
	newUser := &user.User{Name: "Test User", Email: uniqueEmail}
	jsonUser, _ := json.Marshal(newUser)

	req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(jsonUser))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Log the status code and body of the POST response
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code 201. Got %d", w.Code)
	}
	responseBody, _ := io.ReadAll(w.Body)
	log.Printf("POST Response body: %s", responseBody)

	// Decode the created user
	var createdUser user.User
	err := json.NewDecoder(bytes.NewReader(responseBody)).Decode(&createdUser)
	if err != nil {
		t.Errorf("Error decoding response body after user creation: %v", err)
	}

	// Log the created user's details for debugging
	log.Printf("Created user: ID=%d, Name=%s, Email=%s", createdUser.ID, createdUser.Name, createdUser.Email)

	// Then, retrieve the user by ID
	req, _ = http.NewRequest("GET", "/user/"+strconv.Itoa(createdUser.ID), nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Check if the status code is 200 OK
	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Log the GET response body
	responseBody, _ = io.ReadAll(w.Body)
	log.Printf("GET Response body: %s", responseBody)

	// Decode the response to ensure the correct user was returned
	var retrievedUser user.User
	err = json.NewDecoder(bytes.NewReader(responseBody)).Decode(&retrievedUser)
	if err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}

	// Check if the retrieved user's ID matches the created user's ID
	if retrievedUser.ID != createdUser.ID {
		t.Errorf("Expected ID to be '%v'. Got '%v'", createdUser.ID, retrievedUser.ID)
	}

	// Log the retrieved user's details for debugging
	log.Printf("Retrieved user: ID=%d, Name=%s, Email=%s", retrievedUser.ID, retrievedUser.Name, retrievedUser.Email)
}

// TestUpdateUser tests the PUT /user/{id} endpoint to update a user's details
func TestUpdateUser(t *testing.T) {
	r := router.SetupRouter()

	// Generate a unique email to avoid duplicates
	uniqueEmail := "test" + strconv.FormatInt(time.Now().UnixNano(), 10) + "@example.com"

	// First, create a new user with a unique email
	newUser := &user.User{Name: "Test User", Email: uniqueEmail}
	jsonUser, _ := json.Marshal(newUser)

	req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(jsonUser))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Log the POST response body
	responseBody, _ := io.ReadAll(w.Body)
	log.Printf("POST Response body: %s", responseBody)

	// Decode the created user
	var createdUser user.User
	err := json.NewDecoder(bytes.NewReader(responseBody)).Decode(&createdUser)
	if err != nil {
		t.Errorf("Error decoding response body after user creation: %v", err)
	}

	// Log the created user's details
	log.Printf("Created user: ID=%d, Name=%s, Email=%s", createdUser.ID, createdUser.Name, createdUser.Email)

	// Update the user with new data
	updatedUser := &user.User{Name: "Updated User", Email: "updated" + strconv.FormatInt(time.Now().UnixNano(), 10) + "@example.com"}
	jsonUpdatedUser, _ := json.Marshal(updatedUser)

	req, _ = http.NewRequest("PUT", "/user/"+strconv.Itoa(createdUser.ID), bytes.NewBuffer(jsonUpdatedUser))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Log the PUT response body
	responseBody, _ = io.ReadAll(w.Body)
	log.Printf("PUT Response body: %s", responseBody)

	// Check if the status code is 200 OK
	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Decode the response to ensure the user was updated correctly
	var resultUser user.User
	err = json.NewDecoder(bytes.NewReader(responseBody)).Decode(&resultUser)
	if err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}

	// Ensure the updated fields are correct
	if resultUser.Name != updatedUser.Name {
		t.Errorf("Expected Name to be '%v'. Got '%v'", updatedUser.Name, resultUser.Name)
	}

	if resultUser.Email != updatedUser.Email {
		t.Errorf("Expected Email to be '%v'. Got '%v'", updatedUser.Email, resultUser.Email)
	}

	// Log the updated user's details
	log.Printf("Updated user: ID=%d, Name=%s, Email=%s", resultUser.ID, resultUser.Name, resultUser.Email)
}

// TestDeleteUser tests the DELETE /user/{id} endpoint to delete a user by ID
func TestDeleteUser(t *testing.T) {
	r := router.SetupRouter()

	// Generate a unique email to avoid duplicates
	uniqueEmail := "test" + strconv.FormatInt(time.Now().UnixNano(), 10) + "@example.com"

	// First, create a new user with a unique email
	newUser := &user.User{Name: "Test User", Email: uniqueEmail}
	jsonUser, _ := json.Marshal(newUser)

	req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(jsonUser))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Log the POST response body
	responseBody, _ := io.ReadAll(w.Body)
	log.Printf("POST Response body: %s", responseBody)

	// Decode the created user
	var createdUser user.User
	err := json.NewDecoder(bytes.NewReader(responseBody)).Decode(&createdUser)
	if err != nil {
		t.Errorf("Error decoding response body after user creation: %v", err)
	}

	// Log the created user's details
	log.Printf("Created user: ID=%d, Name=%s, Email=%s", createdUser.ID, createdUser.Name, createdUser.Email)

	// Then, delete the user by ID
	req, _ = http.NewRequest("DELETE", "/user/"+strconv.Itoa(createdUser.ID), nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Log the DELETE response body
	responseBody, _ = io.ReadAll(w.Body)
	log.Printf("DELETE Response body: %s", responseBody)

	// Check if the status code is 204 No Content
	if status := w.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}

	// Try to retrieve the deleted user
	req, _ = http.NewRequest("GET", "/user/"+strconv.Itoa(createdUser.ID), nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Log the GET response body for the deleted user
	responseBody, _ = io.ReadAll(w.Body)
	log.Printf("GET Response body: %s", responseBody)

	// Ensure the user no longer exists (404 Not Found)
	if status := w.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}
