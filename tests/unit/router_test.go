package unit

import (
	"github.com/kaushiksamanta/vayu"
	"net/http/httptest"
	"testing"
)

func TestRouterBasic(t *testing.T) {
	app := vayu.New()

	// Test route handling
	app.GET("/test", func(c *vayu.Context, next vayu.NextFunc) {
		_, err := c.Send(vayu.StatusOK, "test route")
		if err != nil {
			t.Fatalf("Failed to send response: %v", err)
		}
	})

	// Create a test request
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// Serve the request
	app.ServeHTTP(w, req)

	// Check response
	if w.Code != vayu.StatusOK {
		t.Errorf("Expected status code %d, got %d", vayu.StatusOK, w.Code)
	}

	if w.Body.String() != "test route" {
		t.Errorf("Expected body 'test route', got '%s'", w.Body.String())
	}
}

func TestRouterParams(t *testing.T) {
	app := vayu.New()

	// Test route with parameters
	app.GET("/users/:id", func(c *vayu.Context, next vayu.NextFunc) {
		id := c.Params["id"]
		_, err := c.Send(vayu.StatusOK, id)
		if err != nil {
			t.Fatalf("Failed to send response: %v", err)
		}
	})

	// Create a test request
	req := httptest.NewRequest("GET", "/users/123", nil)
	w := httptest.NewRecorder()

	// Serve the request
	app.ServeHTTP(w, req)

	// Check response
	if w.Code != vayu.StatusOK {
		t.Errorf("Expected status code %d, got %d", vayu.StatusOK, w.Code)
	}

	if w.Body.String() != "123" {
		t.Errorf("Expected body '123', got '%s'", w.Body.String())
	}
}

func TestRouterMethods(t *testing.T) {
	app := vayu.New()

	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD"}
	path := "/methods-test"

	for _, method := range methods {
		// Register handler for each method dynamically
		switch method {
		case "GET":
			app.GET(path, func(c *vayu.Context, next vayu.NextFunc) {
				_, err := c.Send(vayu.StatusOK, method)
				if err != nil {
					t.Fatalf("Failed to send response: %v", err)
				}
			})
		case "POST":
			app.POST(path, func(c *vayu.Context, next vayu.NextFunc) {
				_, err := c.Send(vayu.StatusOK, method)
				if err != nil {
					t.Fatalf("Failed to send response: %v", err)
				}
			})
		case "PUT":
			app.PUT(path, func(c *vayu.Context, next vayu.NextFunc) {
				_, err := c.Send(vayu.StatusOK, method)
				if err != nil {
					t.Fatalf("Failed to send response: %v", err)
				}
			})
		case "DELETE":
			app.DELETE(path, func(c *vayu.Context, next vayu.NextFunc) {
				_, err := c.Send(vayu.StatusOK, method)
				if err != nil {
					t.Fatalf("Failed to send response: %v", err)
				}
			})
		case "PATCH":
			app.PATCH(path, func(c *vayu.Context, next vayu.NextFunc) {
				_, err := c.Send(vayu.StatusOK, method)
				if err != nil {
					t.Fatalf("Failed to send response: %v", err)
				}
			})
		case "OPTIONS":
			app.OPTIONS(path, func(c *vayu.Context, next vayu.NextFunc) {
				_, err := c.Send(vayu.StatusOK, method)
				if err != nil {
					t.Fatalf("Failed to send response: %v", err)
				}
			})
		case "HEAD":
			app.HEAD(path, func(c *vayu.Context, next vayu.NextFunc) {
				_, err := c.Send(vayu.StatusOK, method)
				if err != nil {
					t.Fatalf("Failed to send response: %v", err)
				}
			})
		}

		// Create test request for each method
		req := httptest.NewRequest(method, path, nil)
		w := httptest.NewRecorder()

		// Serve the request
		app.ServeHTTP(w, req)

		// Check response
		if w.Code != vayu.StatusOK {
			t.Errorf("Method %s: Expected status code %d, got %d", method, vayu.StatusOK, w.Code)
		}

		// For HEAD method, body is not included in response
		if method != "HEAD" && w.Body.String() != method {
			t.Errorf("Method %s: Expected body '%s', got '%s'", method, method, w.Body.String())
		}
	}
}

func TestRouteNotFound(t *testing.T) {
	app := vayu.New()

	// Create a test request to a non-existent route
	req := httptest.NewRequest("GET", "/not-found", nil)
	w := httptest.NewRecorder()

	// Serve the request
	app.ServeHTTP(w, req)

	// Check response
	if w.Code != vayu.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", vayu.StatusNotFound, w.Code)
	}
}
