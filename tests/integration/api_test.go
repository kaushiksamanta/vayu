package integration

import (
	"bytes"
	"encoding/json"
	"github.com/kaushiksamanta/vayu"
	"net/http/httptest"
	"testing"
)

// Setup creates a new Vayu app with a sample API for testing
func setupAPITest() *vayu.App {
	app := vayu.New()

	// Add middleware for all routes
	app.Use(vayu.Logger())
	app.Use(vayu.Recovery())

	// Set up a sample API
	api := app.Group("/api")

	// GET endpoint
	api.GET("/users", func(c *vayu.Context, next vayu.NextFunc) {
		users := []map[string]interface{}{
			{"id": "1", "name": "Alice", "role": "admin"},
			{"id": "2", "name": "Bob", "role": "user"},
			{"id": "3", "name": "Charlie", "role": "user"},
		}

		err := c.OK(users)
		if err != nil {
			c.InternalServerError("Failed to send response")
		}
	})

	// GET with parameter
	api.GET("/users/:id", func(c *vayu.Context, next vayu.NextFunc) {
		id := c.Params["id"]

		// Sample user data
		users := map[string]map[string]string{
			"1": {"name": "Alice", "role": "admin"},
			"2": {"name": "Bob", "role": "user"},
			"3": {"name": "Charlie", "role": "user"},
		}

		if user, exists := users[id]; exists {
			err := c.OK(map[string]interface{}{
				"id":   id,
				"name": user["name"],
				"role": user["role"],
			})
			if err != nil {
				c.InternalServerError("Failed to send response")
			}
		} else {
			c.NotFound("User not found")
		}
	})

	// POST endpoint
	api.POST("/users", func(c *vayu.Context, next vayu.NextFunc) {
		type UserInput struct {
			Name string `json:"name"`
			Role string `json:"role"`
		}

		var input UserInput
		decoder := json.NewDecoder(c.Request.Body)
		if err := decoder.Decode(&input); err != nil {
			c.BadRequest("Invalid request body")
			return
		}

		// Validate input
		if input.Name == "" {
			c.BadRequest("Name is required")
			return
		}

		// In a real app, we would save to database here
		// For test, just return a success response with a mock ID
		err := c.Created(map[string]interface{}{
			"id":      "4", // Mock ID for new user
			"name":    input.Name,
			"role":    input.Role,
			"message": "User created successfully",
		})
		if err != nil {
			c.InternalServerError("Failed to send response")
		}
	})

	return app
}

func TestGetUsers(t *testing.T) {
	app := setupAPITest()

	// Create request to get all users
	req := httptest.NewRequest("GET", "/api/users", nil)
	w := httptest.NewRecorder()

	// Serve the request
	app.ServeHTTP(w, req)

	// Check response
	if w.Code != vayu.StatusOK {
		t.Errorf("Expected status code %d, got %d", vayu.StatusOK, w.Code)
	}

	// Decode response body
	var users []map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &users)
	if err != nil {
		t.Fatalf("Failed to decode JSON response: %v", err)
	}

	// Check that we got 3 users
	if len(users) != 3 {
		t.Errorf("Expected 3 users, got %d", len(users))
	}
}

func TestGetUserById(t *testing.T) {
	app := setupAPITest()

	// Create request to get a specific user
	req := httptest.NewRequest("GET", "/api/users/2", nil)
	w := httptest.NewRecorder()

	// Serve the request
	app.ServeHTTP(w, req)

	// Check response
	if w.Code != vayu.StatusOK {
		t.Errorf("Expected status code %d, got %d", vayu.StatusOK, w.Code)
	}

	// Decode response body
	var user map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &user)
	if err != nil {
		t.Fatalf("Failed to decode JSON response: %v", err)
	}

	// Check user details
	if user["id"] != "2" {
		t.Errorf("Expected user ID '2', got '%v'", user["id"])
	}

	if user["name"] != "Bob" {
		t.Errorf("Expected user name 'Bob', got '%v'", user["name"])
	}
}

func TestGetNonExistentUser(t *testing.T) {
	app := setupAPITest()

	// Create request to get a non-existent user
	req := httptest.NewRequest("GET", "/api/users/999", nil)
	w := httptest.NewRecorder()

	// Serve the request
	app.ServeHTTP(w, req)

	// Check response
	if w.Code != vayu.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", vayu.StatusNotFound, w.Code)
	}

	// Decode response body
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to decode JSON response: %v", err)
	}

	// Check error message
	if response["error"] != "User not found" {
		t.Errorf("Expected error 'User not found', got '%s'", response["error"])
	}
}

func TestCreateUser(t *testing.T) {
	app := setupAPITest()

	// Create request body
	input := map[string]string{
		"name": "Dave",
		"role": "user",
	}

	body, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	// Create request to create a new user
	req := httptest.NewRequest("POST", "/api/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Serve the request
	app.ServeHTTP(w, req)

	// Check response
	if w.Code != vayu.StatusCreated {
		t.Errorf("Expected status code %d, got %d", vayu.StatusCreated, w.Code)
	}

	// Decode response body
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to decode JSON response: %v", err)
	}

	// Check response fields
	if response["name"] != "Dave" {
		t.Errorf("Expected name 'Dave', got '%v'", response["name"])
	}

	if response["role"] != "user" {
		t.Errorf("Expected role 'user', got '%v'", response["role"])
	}

	if response["message"] != "User created successfully" {
		t.Errorf("Expected message 'User created successfully', got '%v'", response["message"])
	}
}
