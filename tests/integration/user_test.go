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

var testHandler *user.Handler

// TestMain initializes the test environment
func TestMain(m *testing.M) {
	cfg := config.LoadConfig()

	// Connect to the database
	database.ConnectToDB(cfg)

	if database.DB == nil {
		log.Fatal("Database is not initialized")
	}

	// Provide service and handler for tests
	service := user.ProvideUserService(database.DB)
	testHandler = user.ProvideUserHandler(service)

	code := m.Run()
	os.Exit(code)
}

// Helper function to create a unique email
func generateUniqueEmail() string {
	return "test" + strconv.FormatInt(time.Now().UnixNano(), 10) + "@example.com"
}

// Helper function to decode response body into a struct
func decodeResponseBody(t *testing.T, body []byte, target interface{}) {
	err := json.NewDecoder(bytes.NewReader(body)).Decode(target)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}
}

// TestCreateUser tests creating a user
func TestCreateUser(t *testing.T) {
	r := router.SetupRouter(testHandler)

	uniqueEmail := generateUniqueEmail()
	newUser := &user.User{Name: "Test User", Email: uniqueEmail}
	jsonUser, _ := json.Marshal(newUser)

	req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(jsonUser))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	responseBody, _ := io.ReadAll(w.Body)
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code 201. Got %d", w.Code)
	}

	var createdUser user.User
	decodeResponseBody(t, responseBody, &createdUser)

	if createdUser.Name != newUser.Name {
		t.Errorf("Expected Name '%v'. Got '%v'", newUser.Name, createdUser.Name)
	}
	if createdUser.Email != uniqueEmail {
		t.Errorf("Expected Email '%v'. Got '%v'", uniqueEmail, createdUser.Email)
	}
	if createdUser.ID == 0 {
		t.Error("Expected a valid ID")
	}
	if createdUser.CreatedAt == "" {
		t.Error("Expected a valid CreatedAt timestamp")
	}
}

// TestGetUsers tests fetching all users
func TestGetUsers(t *testing.T) {
	r := router.SetupRouter(testHandler)

	req, _ := http.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code 200. Got %d", w.Code)
	}

	var users []user.User
	decodeResponseBody(t, w.Body.Bytes(), &users)

	if len(users) == 0 {
		t.Error("Expected at least one user in the response")
	}
}

// TestGetUserByID tests fetching a user by ID
func TestGetUserByID(t *testing.T) {
	r := router.SetupRouter(testHandler)

	uniqueEmail := generateUniqueEmail()
	newUser := &user.User{Name: "Test User", Email: uniqueEmail}
	jsonUser, _ := json.Marshal(newUser)

	req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(jsonUser))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var createdUser user.User
	decodeResponseBody(t, w.Body.Bytes(), &createdUser)

	req, _ = http.NewRequest("GET", "/user/"+strconv.Itoa(createdUser.ID), nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code 200. Got %d", w.Code)
	}

	var retrievedUser user.User
	decodeResponseBody(t, w.Body.Bytes(), &retrievedUser)

	if retrievedUser.ID != createdUser.ID {
		t.Errorf("Expected ID '%v'. Got '%v'", createdUser.ID, retrievedUser.ID)
	}
}

// TestUpdateUser tests updating user details
func TestUpdateUser(t *testing.T) {
	r := router.SetupRouter(testHandler)

	uniqueEmail := generateUniqueEmail()
	newUser := &user.User{Name: "Test User", Email: uniqueEmail}
	jsonUser, _ := json.Marshal(newUser)

	req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(jsonUser))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var createdUser user.User
	decodeResponseBody(t, w.Body.Bytes(), &createdUser)

	updatedUser := &user.User{Name: "Updated User", Email: "updated" + uniqueEmail}
	jsonUpdatedUser, _ := json.Marshal(updatedUser)

	req, _ = http.NewRequest("PUT", "/user/"+strconv.Itoa(createdUser.ID), bytes.NewBuffer(jsonUpdatedUser))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code 200. Got %d", w.Code)
	}

	var resultUser user.User
	decodeResponseBody(t, w.Body.Bytes(), &resultUser)

	if resultUser.Name != updatedUser.Name {
		t.Errorf("Expected Name '%v'. Got '%v'", updatedUser.Name, resultUser.Name)
	}
	if resultUser.Email != updatedUser.Email {
		t.Errorf("Expected Email '%v'. Got '%v'", updatedUser.Email, resultUser.Email)
	}
}

// TestDeleteUser tests deleting a user by ID
func TestDeleteUser(t *testing.T) {
	r := router.SetupRouter(testHandler)

	uniqueEmail := generateUniqueEmail()
	newUser := &user.User{Name: "Test User", Email: uniqueEmail}
	jsonUser, _ := json.Marshal(newUser)

	req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(jsonUser))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var createdUser user.User
	decodeResponseBody(t, w.Body.Bytes(), &createdUser)

	req, _ = http.NewRequest("DELETE", "/user/"+strconv.Itoa(createdUser.ID), nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status code 204. Got %d", w.Code)
	}

	req, _ = http.NewRequest("GET", "/user/"+strconv.Itoa(createdUser.ID), nil)
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code 404. Got %d", w.Code)
	}
}
