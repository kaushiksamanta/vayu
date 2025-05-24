package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/kaushiksamanta/vayu"
)

// setupErrorRoutes adds routes that demonstrate error handling
func setupErrorRoutes(app *vayu.App) {
	// Custom error handler
	customErrorHandler := func(c *vayu.Context, err error) {
		// Log the error
		fmt.Printf("Custom error handler caught: %v\n", err)

		// Send appropriate response based on error type
		c.JSON(vayu.StatusInternalServerError, map[string]string{
			"error":   "Something went wrong",
			"message": err.Error(),
		})
	}

	// APPROACH 1: ErrorHandlerMiddleware
	// This registers error handling as middleware for an entire group of routes
	errorDemo := app.Group("/errors")
	errorDemo.Use(vayu.ErrorHandlerMiddleware(customErrorHandler)) // Applied to ALL routes in this group

	// Route that returns a 400 Bad Request using helper
	errorDemo.GET("/bad-request", func(c *vayu.Context, next vayu.NextFunc) {
		c.BadRequest("This is a demonstration of a bad request")
	})

	// Route that returns a 404 Not Found using helper
	errorDemo.GET("/not-found", func(c *vayu.Context, next vayu.NextFunc) {
		c.NotFound("Resource not found")
	})

	// Route that deliberately panics (will be caught by middleware)
	errorDemo.GET("/panic", func(c *vayu.Context, next vayu.NextFunc) {
		panic("This is a deliberate panic for demonstration")
	})

	// Route that returns an error (will be caught by middleware)
	errorDemo.GET("/error", func(c *vayu.Context, next vayu.NextFunc) {
		// Simulate a database error or other runtime error
		err := errors.New("simulated error for demonstration")
		panic(err)
	})

	// APPROACH 2: WithErrorHandling
	// This wraps a specific handler with error handling, without affecting other routes
	// Create a separate error handler that shows it's from the wrapped handler
	wrappedErrorHandler := func(c *vayu.Context, err error) {
		fmt.Printf("Wrapped handler caught: %v\n", err)
		c.JSON(vayu.StatusInternalServerError, map[string]string{
			"error":   "Error in wrapped handler",
			"message": err.Error(),
			"handler": "This error was caught by WithErrorHandling",
		})
	}

	// A handler that will panic
	panicHandler := func(c *vayu.Context, next vayu.NextFunc) {
		panic("This panic is caught by WithErrorHandling wrapper")
	}

	// Register the route WITH the handler wrapped in error handling
	app.GET("/wrapped-error", vayu.WithErrorHandling(panicHandler, wrappedErrorHandler))

	// For comparison, register a similar route directly without middleware
	// WARNING: This would crash the server without global middleware!
	app.GET("/global-middleware-demo", func(c *vayu.Context, next vayu.NextFunc) {
		panic("This panic is caught by global middleware")
	})

	// Update the log message to include the new routes
	fmt.Println("  GET  /wrapped-error       - WithErrorHandling demo")
	fmt.Println("  GET  /global-middleware-demo - Global middleware demo")
}

func main() {
	app := vayu.New()

	// Global middleware
	app.Use(vayu.Logger())
	app.Use(vayu.Recovery())

	// Basic routes
	app.GET("/hello", func(c *vayu.Context, next vayu.NextFunc) {
		_, err := c.Send(vayu.StatusOK, "Hello, World!")
		if err != nil {
			log.Printf("Error sending response: %v", err)
		}
	})

	app.GET("/json", func(c *vayu.Context, next vayu.NextFunc) {
		err := c.OK(map[string]interface{}{
			"message": "Hello, JSON!",
			"status":  "success",
		})
		if err != nil {
			log.Printf("Error sending JSON response: %v", err)
		}
	})

	// Route with parameter
	app.GET("/users/:id", func(c *vayu.Context, next vayu.NextFunc) {
		id := c.Params["id"]
		err := c.OK(map[string]string{
			"userId": id,
		})
		if err != nil {
			log.Printf("Error sending JSON response: %v", err)
		}
	})

	// Route group
	api := app.Group("/api")
	api.GET("/status", func(c *vayu.Context, next vayu.NextFunc) {
		err := c.OK(map[string]string{
			"status": "API is running",
		})
		if err != nil {
			log.Printf("Error sending JSON response: %v", err)
		}
	})

	// Static files
	app.Static("/public", "./public")

	// Set up error handling demonstration routes
	setupErrorRoutes(app)

	// Print routes for testing
	log.Println("Available Routes:")
	log.Println("  GET  /hello                 - Basic text response")
	log.Println("  GET  /json                  - Basic JSON response")
	log.Println("  GET  /users/:id             - Route with parameter")
	log.Println("  GET  /api/status            - API status endpoint")
	log.Println("  GET  /errors/bad-request    - Error handling demo (group middleware)")
	log.Println("  GET  /errors/not-found      - Not found demo (group middleware)")
	log.Println("  GET  /errors/panic          - Panic recovery demo (group middleware)")
	log.Println("  GET  /errors/error          - Error handling demo (group middleware)")
	log.Println("  GET  /wrapped-error         - WithErrorHandling demo (handler wrapper)")
	log.Println("  GET  /global-middleware-demo - Global middleware demo")

	// Start server
	log.Println("Server starting on http://localhost:8080")
	if err := app.Listen(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
