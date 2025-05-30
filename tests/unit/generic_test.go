package unit

import (
	"encoding/json"
	"github.com/kaushiksamanta/vayu"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestJSONResponse tests the generic JSONResponse function
func TestJSONResponse(t *testing.T) {
	app := vayu.New()
	
	// Define a test struct with different types of fields
	type TestResponse struct {
		StringField string   `json:"string_field"`
		IntField    int      `json:"int_field"`
		BoolField   bool     `json:"bool_field"`
		FloatField  float64  `json:"float_field"`
		ArrayField  []string `json:"array_field"`
		NestedField struct {
			Key string `json:"key"`
		} `json:"nested_field"`
	}
	
	app.GET("/generic-json", func(c *vayu.Context, next vayu.NextFunc) {
		response := TestResponse{
			StringField: "test string",
			IntField:    42,
			BoolField:   true,
			FloatField:  3.14159,
			ArrayField:  []string{"one", "two", "three"},
		}
		response.NestedField.Key = "nested value"
		
		err := vayu.JSONResponse(c, vayu.StatusOK, response)
		if err != nil {
			t.Fatalf("Failed to send JSON response: %v", err)
		}
	})
	
	// Test the endpoint
	req := httptest.NewRequest("GET", "/generic-json", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)
	
	// Verify response status
	if w.Code != vayu.StatusOK {
		t.Errorf("Expected status code %d, got %d", vayu.StatusOK, w.Code)
	}
	
	// Verify content type
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
	}
	
	// Verify JSON content
	var result TestResponse
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	
	// Verify individual fields
	if result.StringField != "test string" {
		t.Errorf("Expected StringField 'test string', got '%s'", result.StringField)
	}
	if result.IntField != 42 {
		t.Errorf("Expected IntField 42, got %d", result.IntField)
	}
	if result.BoolField != true {
		t.Errorf("Expected BoolField true, got %v", result.BoolField)
	}
	if result.FloatField != 3.14159 {
		t.Errorf("Expected FloatField 3.14159, got %v", result.FloatField)
	}
	if len(result.ArrayField) != 3 || result.ArrayField[0] != "one" {
		t.Errorf("ArrayField does not match expected value: %v", result.ArrayField)
	}
	if result.NestedField.Key != "nested value" {
		t.Errorf("Expected NestedField.Key 'nested value', got '%s'", result.NestedField.Key)
	}
}

// TestBindJSONRequest tests the generic BindJSONRequest function
func TestBindJSONBody(t *testing.T) {
	app := vayu.New()
	
	// Define a test struct for binding
	type TestUser struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Age      int    `json:"age"`
		IsAdmin  bool   `json:"is_admin"`
	}
	
	app.POST("/bind-test", func(c *vayu.Context, next vayu.NextFunc) {
		user, err := vayu.BindJSONBody[TestUser](c)
		if err != nil {
			c.JSON(vayu.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		
		// Return the bound user to verify
		err = vayu.JSONResponse(c, vayu.StatusOK, user)
		if err != nil {
			t.Fatalf("Failed to send response: %v", err)
		}
	})
	
	// Test valid JSON binding
	t.Run("Valid JSON", func(t *testing.T) {
		validJSON := `{
			"username": "testuser",
			"email": "test@example.com",
			"age": 25,
			"is_admin": true
		}`
		
		req := httptest.NewRequest("POST", "/bind-test", strings.NewReader(validJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)
		
		// Check status
		if w.Code != vayu.StatusOK {
			t.Errorf("Expected status code %d, got %d", vayu.StatusOK, w.Code)
		}
		
		// Verify response content
		var response TestUser
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}
		
		if response.Username != "testuser" {
			t.Errorf("Expected Username 'testuser', got '%s'", response.Username)
		}
		if response.Email != "test@example.com" {
			t.Errorf("Expected Email 'test@example.com', got '%s'", response.Email)
		}
		if response.Age != 25 {
			t.Errorf("Expected Age 25, got %d", response.Age)
		}
		if response.IsAdmin != true {
			t.Errorf("Expected IsAdmin true, got %v", response.IsAdmin)
		}
	})
	
	// Test invalid JSON binding
	t.Run("Invalid JSON", func(t *testing.T) {
		invalidJSON := `{"username": "testuser", "email": "test@example.com", "age": "not-a-number"}`
		
		req := httptest.NewRequest("POST", "/bind-test", strings.NewReader(invalidJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)
		
		// Check for bad request status
		if w.Code != vayu.StatusBadRequest {
			t.Errorf("Expected status code %d for invalid JSON, got %d", vayu.StatusBadRequest, w.Code)
		}
	})
}

// TestContextStoreGenerics tests the generic context store functions
func TestContextStoreGenerics(t *testing.T) {
	app := vayu.New()
	
	type CustomType struct {
		Name  string
		Value int
	}
	
	app.GET("/context-store", func(c *vayu.Context, next vayu.NextFunc) {
		// Test with string
		vayu.SetValue(c, "string_key", "string value")
		stringVal, stringOk := vayu.GetValue[string](c, "string_key")
		if !stringOk || stringVal != "string value" {
			t.Errorf("Failed to get string value: got %v, ok=%v", stringVal, stringOk)
		}
		
		// Test with int
		vayu.SetValue(c, "int_key", 42)
		intVal, intOk := vayu.GetValue[int](c, "int_key")
		if !intOk || intVal != 42 {
			t.Errorf("Failed to get int value: got %v, ok=%v", intVal, intOk)
		}
		
		// Test with bool
		vayu.SetValue(c, "bool_key", true)
		boolVal, boolOk := vayu.GetValue[bool](c, "bool_key")
		if !boolOk || !boolVal {
			t.Errorf("Failed to get bool value: got %v, ok=%v", boolVal, boolOk)
		}
		
		// Test with custom struct
		custom := CustomType{Name: "test", Value: 100}
		vayu.SetValue(c, "custom_key", custom)
		customVal, customOk := vayu.GetValue[CustomType](c, "custom_key")
		if !customOk || customVal.Name != "test" || customVal.Value != 100 {
			t.Errorf("Failed to get custom value: got %+v, ok=%v", customVal, customOk)
		}
		
		// Test with non-existent key
		_, nonExistOk := vayu.GetValue[string](c, "non_existent")
		if nonExistOk {
			t.Errorf("GetValue with non-existent key should return ok=false")
		}
		
		// Test with type mismatch
		vayu.SetValue(c, "type_mismatch", 123)
		_, typeMismatchOk := vayu.GetValue[string](c, "type_mismatch")
		if typeMismatchOk {
			t.Errorf("GetValue with type mismatch should return ok=false")
		}
		
		// Verify test completion
		c.JSON(vayu.StatusOK, map[string]string{"status": "ok"})
	})
	
	// Run the test handler
	req := httptest.NewRequest("GET", "/context-store", nil)
	w := httptest.NewRecorder()
	app.ServeHTTP(w, req)
	
	// Verify the handler ran successfully
	if w.Code != vayu.StatusOK {
		t.Errorf("Expected status code %d, got %d", vayu.StatusOK, w.Code)
	}
}

// TestMustBindJSONRequest tests the MustBindJSONRequest function
func TestMustBindJSONBody(t *testing.T) {
	// This test will verify that MustBindJSONRequest panics on invalid input
	
	// Define a recover function to catch panics
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("MustBindJSONRequest should have panicked with invalid JSON")
		}
	}()
	
	// Create context with invalid JSON
	req := httptest.NewRequest("POST", "/test", strings.NewReader(`{"invalid": json`))
	w := httptest.NewRecorder()
	c := &vayu.Context{
		Writer:  &vayu.ResponseWriter{ResponseWriter: w},
		Request: req,
	}
	
	// This should panic
	_ = vayu.MustBindJSONBody[struct{}](c)
}
