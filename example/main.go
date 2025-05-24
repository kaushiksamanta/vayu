package main

import (
	"errors"
	"fmt"
	"github.com/kaushiksamanta/vayu"
	"log"
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

	// Add error handling middleware to a group
	errorDemo := app.Group("/errors")
	errorDemo.Use(vayu.ErrorHandlerMiddleware(customErrorHandler))

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
	log.Println("  GET  /hello            - Basic text response")
	log.Println("  GET  /json             - Basic JSON response")
	log.Println("  GET  /users/:id        - Route with parameter")
	log.Println("  GET  /api/status       - API status endpoint")
	log.Println("  GET  /errors/bad-request - Error handling demo")
	log.Println("  GET  /errors/not-found   - Not found demo")
	log.Println("  GET  /errors/panic       - Panic recovery demo")
	log.Println("  GET  /errors/error       - Error handling demo")

	// Start server
	log.Println("Server starting on http://localhost:8080")
	if err := app.Listen(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
