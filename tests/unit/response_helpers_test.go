package unit

import (
	"encoding/json"
	"github.com/kaushiksamanta/vayu"
	"net/http/httptest"
	"testing"
)

func TestResponseHelperOK(t *testing.T) {
	app := vayu.New()

	app.GET("/ok", func(c *vayu.Context, next vayu.NextFunc) {
		err := c.OK(map[string]string{"message": "success"})
		if err != nil {
			t.Fatalf("Failed to send OK response: %v", err)
		}
	})

	req := httptest.NewRequest("GET", "/ok", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	if w.Code != vayu.StatusOK {
		t.Errorf("Expected status code %d, got %d", vayu.StatusOK, w.Code)
	}

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to decode JSON response: %v", err)
	}

	if response["message"] != "success" {
		t.Errorf("Expected message 'success', got '%s'", response["message"])
	}
}

func TestResponseHelperCreated(t *testing.T) {
	app := vayu.New()

	app.POST("/created", func(c *vayu.Context, next vayu.NextFunc) {
		err := c.Created(map[string]string{"id": "123", "status": "created"})
		if err != nil {
			t.Fatalf("Failed to send Created response: %v", err)
		}
	})

	req := httptest.NewRequest("POST", "/created", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	if w.Code != vayu.StatusCreated {
		t.Errorf("Expected status code %d, got %d", vayu.StatusCreated, w.Code)
	}

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to decode JSON response: %v", err)
	}

	if response["id"] != "123" || response["status"] != "created" {
		t.Errorf("Expected id '123' and status 'created', got id '%s' and status '%s'",
			response["id"], response["status"])
	}
}

func TestResponseHelperNoContent(t *testing.T) {
	app := vayu.New()

	app.DELETE("/no-content", func(c *vayu.Context, next vayu.NextFunc) {
		_, err := c.NoContent()
		if err != nil {
			t.Fatalf("Failed to send NoContent response: %v", err)
		}
	})

	req := httptest.NewRequest("DELETE", "/no-content", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	if w.Code != vayu.StatusNoContent {
		t.Errorf("Expected status code %d, got %d", vayu.StatusNoContent, w.Code)
	}

	if w.Body.Len() != 0 {
		t.Errorf("Expected empty body, got '%s'", w.Body.String())
	}
}

func TestResponseHelperBadRequest(t *testing.T) {
	app := vayu.New()

	app.GET("/bad-request", func(c *vayu.Context, next vayu.NextFunc) {
		err := c.BadRequest("Invalid parameters")
		if err != nil {
			t.Fatalf("Failed to send BadRequest response: %v", err)
		}
	})

	req := httptest.NewRequest("GET", "/bad-request", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	if w.Code != vayu.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", vayu.StatusBadRequest, w.Code)
	}

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to decode JSON response: %v", err)
	}

	if response["error"] != "Invalid parameters" {
		t.Errorf("Expected error 'Invalid parameters', got '%s'", response["error"])
	}
}

func TestResponseHelperNotFound(t *testing.T) {
	app := vayu.New()

	app.GET("/not-found", func(c *vayu.Context, next vayu.NextFunc) {
		err := c.NotFound("Resource not found")
		if err != nil {
			t.Fatalf("Failed to send NotFound response: %v", err)
		}
	})

	req := httptest.NewRequest("GET", "/not-found", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	if w.Code != vayu.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", vayu.StatusNotFound, w.Code)
	}

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to decode JSON response: %v", err)
	}

	if response["error"] != "Resource not found" {
		t.Errorf("Expected error 'Resource not found', got '%s'", response["error"])
	}
}

func TestResponseHelperInternalServerError(t *testing.T) {
	app := vayu.New()

	app.GET("/server-error", func(c *vayu.Context, next vayu.NextFunc) {
		err := c.InternalServerError("Something went wrong")
		if err != nil {
			t.Fatalf("Failed to send InternalServerError response: %v", err)
		}
	})

	req := httptest.NewRequest("GET", "/server-error", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	if w.Code != vayu.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", vayu.StatusInternalServerError, w.Code)
	}

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to decode JSON response: %v", err)
	}

	if response["error"] != "Something went wrong" {
		t.Errorf("Expected error 'Something went wrong', got '%s'", response["error"])
	}
}
