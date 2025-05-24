package unit

import (
	"github.com/kaushiksamanta/vayu"
	"net/http/httptest"
	"testing"
)

func TestMiddlewareExecution(t *testing.T) {
	app := vayu.New()

	// Variables to track middleware execution
	var (
		middleware1Called bool
		middleware2Called bool
		handlerCalled     bool
	)

	// Add middleware
	app.Use(func(c *vayu.Context, next vayu.NextFunc) {
		middleware1Called = true
		next()
	})

	app.Use(func(c *vayu.Context, next vayu.NextFunc) {
		middleware2Called = true
		next()
	})

	// Add route handler
	app.GET("/middleware-test", func(c *vayu.Context, next vayu.NextFunc) {
		handlerCalled = true
		_, err := c.Send(vayu.StatusOK, "middleware test")
		if err != nil {
			t.Fatalf("Failed to send response: %v", err)
		}
	})

	// Create a test request
	req := httptest.NewRequest("GET", "/middleware-test", nil)
	w := httptest.NewRecorder()

	// Serve the request
	app.ServeHTTP(w, req)

	// Check that all middleware and the handler were called
	if !middleware1Called {
		t.Error("First middleware was not called")
	}

	if !middleware2Called {
		t.Error("Second middleware was not called")
	}

	if !handlerCalled {
		t.Error("Route handler was not called")
	}

	// Check response
	if w.Code != vayu.StatusOK {
		t.Errorf("Expected status code %d, got %d", vayu.StatusOK, w.Code)
	}

	if w.Body.String() != "middleware test" {
		t.Errorf("Expected body 'middleware test', got '%s'", w.Body.String())
	}
}

func TestMiddlewareOrder(t *testing.T) {
	app := vayu.New()

	// Track the order of execution
	executionOrder := []string{}

	// Add middleware
	app.Use(func(c *vayu.Context, next vayu.NextFunc) {
		executionOrder = append(executionOrder, "before_mw1")
		next()
		executionOrder = append(executionOrder, "after_mw1")
	})

	app.Use(func(c *vayu.Context, next vayu.NextFunc) {
		executionOrder = append(executionOrder, "before_mw2")
		next()
		executionOrder = append(executionOrder, "after_mw2")
	})

	// Add route handler
	app.GET("/order-test", func(c *vayu.Context, next vayu.NextFunc) {
		executionOrder = append(executionOrder, "handler")
		_, err := c.Send(vayu.StatusOK, "order test")
		if err != nil {
			t.Fatalf("Failed to send response: %v", err)
		}
	})

	// Create a test request
	req := httptest.NewRequest("GET", "/order-test", nil)
	w := httptest.NewRecorder()

	// Serve the request
	app.ServeHTTP(w, req)

	// Expected order of execution
	expected := []string{
		"before_mw1",
		"before_mw2",
		"handler",
		"after_mw2",
		"after_mw1",
	}

	// Check that the execution order is correct
	if len(executionOrder) != len(expected) {
		t.Errorf("Expected %d middleware executions, got %d", len(expected), len(executionOrder))
	}

	for i, step := range expected {
		if i >= len(executionOrder) || executionOrder[i] != step {
			t.Errorf("Expected execution step %d to be '%s', got '%s'", i, step, executionOrder[i])
		}
	}
}

func TestErrorHandlerMiddleware(t *testing.T) {
	app := vayu.New()

	// Add error handling middleware
	var errorHandled bool

	errorHandler := func(c *vayu.Context, err error) {
		errorHandled = true
		_ = c.JSON(vayu.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	app.Use(vayu.ErrorHandlerMiddleware(errorHandler))

	// Add route that panics
	app.GET("/panic-test", func(c *vayu.Context, next vayu.NextFunc) {
		panic("test panic")
	})

	// Create a test request
	req := httptest.NewRequest("GET", "/panic-test", nil)
	w := httptest.NewRecorder()

	// Serve the request
	app.ServeHTTP(w, req)

	// Check that the error was handled
	if !errorHandled {
		t.Error("Error was not handled by the error handler")
	}

	// Check response
	if w.Code != vayu.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", vayu.StatusInternalServerError, w.Code)
	}
}
