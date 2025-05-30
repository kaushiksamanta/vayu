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

	// Generics demonstration routes
	gen := app.Group("/generics")

	// Demonstrate JSONT - Type-safe JSON response with generics
	gen.GET("/json", func(c *vayu.Context, next vayu.NextFunc) {
		// Define a typed response structure
		type ApiResponse struct {
			Message string   `json:"message"`
			Status  int      `json:"status"`
			Items   []string `json:"items"`
		}

		// Create a response with the correct type
		response := ApiResponse{
			Message: "This response is type-safe with generics",
			Status:  vayu.StatusOK,
			Items:   []string{"item1", "item2", "item3"},
		}

		// Use the generic JSONResponse function
		err := vayu.JSONResponse(c, vayu.StatusOK, response)
		if err != nil {
			log.Printf("Error sending JSON response: %v", err)
		}
	})

	// Demonstrate SetT and GetT - Type-safe context store with generics
	gen.GET("/store", func(c *vayu.Context, next vayu.NextFunc) {
		// Define a custom type
		type User struct {
			Name  string
			Email string
			Age   int
		}

		// Store a typed value using generics
		user := User{Name: "John", Email: "john@example.com", Age: 30}
		vayu.SetValue(c, "current_user", user)

		// Retrieve with type safety
		retrievedUser, ok := vayu.GetValue[User](c, "current_user")

		if !ok {
			c.JSON(vayu.StatusInternalServerError, map[string]string{"error": "Failed to retrieve user"})
			return
		}

		// We can access the fields with full type safety
		err := vayu.JSONResponse(c, vayu.StatusOK, map[string]interface{}{
			"user":             retrievedUser,
			"name_from_store":  retrievedUser.Name,
			"email_from_store": retrievedUser.Email,
			"age_from_store":   retrievedUser.Age,
		})

		if err != nil {
			log.Printf("Error sending response: %v", err)
		}
	})

	// Demonstrate BindJSONT - Type-safe binding with generics
	gen.POST("/bind", func(c *vayu.Context, next vayu.NextFunc) {
		// Define the expected request structure
		type CreateUserRequest struct {
			Username string `json:"username"`
			Email    string `json:"email"`
			Age      int    `json:"age"`
		}

		// Use generic binding
		userRequest, err := vayu.BindJSONBody[CreateUserRequest](c)
		if err != nil {
			c.JSON(vayu.StatusBadRequest, map[string]string{
				"error": "Invalid request body: " + err.Error(),
			})
			return
		}

		// Create a response using the bound data
		response := map[string]interface{}{
			"message": "User created successfully",
			"user": map[string]interface{}{
				"username": userRequest.Username,
				"email":    userRequest.Email,
				"age":      userRequest.Age,
			},
		}

		err = vayu.JSONResponse(c, vayu.StatusCreated, response)
		if err != nil {
			log.Printf("Error sending response: %v", err)
		}
	})

	// Static files
	app.Static("/public", "./public")

	// Set up query and path parameter binding demonstrations
	// Example: /search-filters?filter={"category":"books","minPrice":10.5,"inStock":true}
	app.GET("/search-filters", func(c *vayu.Context, next vayu.NextFunc) {
		// Define a filter type for the query parameter
		type Filter struct {
			Category string  `json:"category"`
			MinPrice float64 `json:"minPrice"`
			InStock  bool    `json:"inStock"`
		}

		// Bind the JSON from the query parameter
		filter, err := vayu.BindQueryJSON[Filter](c, "filter")
		if err != nil {
			c.JSON(vayu.StatusBadRequest, map[string]string{
				"error": "Invalid filter format: " + err.Error(),
			})
			return
		}

		// Use the typed filter to send a response
		vayu.JSONResponse(c, vayu.StatusOK, map[string]interface{}{
			"message": "Filter applied successfully",
			"filter":  filter,
		})
	})

	// Example: /products/:config
	// Where config is URL-encoded JSON like %7B%22view%22%3A%22grid%22%2C%22showPrices%22%3Atrue%7D
	app.GET("/products/:config", func(c *vayu.Context, next vayu.NextFunc) {
		// Define a configuration type for the path parameter
		type ProductConfig struct {
			View       string `json:"view"`
			ShowPrices bool   `json:"showPrices"`
		}

		// Bind the JSON from the path parameter
		config, err := vayu.BindParamJSON[ProductConfig](c, "config")
		if err != nil {
			c.JSON(vayu.StatusBadRequest, map[string]string{
				"error": "Invalid configuration format: " + err.Error(),
			})
			return
		}

		// Use the typed configuration to send a response
		vayu.JSONResponse(c, vayu.StatusOK, map[string]interface{}{
			"message": "Products configured successfully",
			"config":  config,
		})
	})

	// Example: /search?q=golang&page=2&per_page=20&sort=relevance&desc=true
	app.GET("/search", func(c *vayu.Context, next vayu.NextFunc) {
		// Define a search parameters type with query tags
		type SearchParams struct {
			Term    string `query:"q" required:"true"`
			Page    int    `query:"page"`
			PerPage int    `query:"per_page"`
			SortBy  string `query:"sort"`
			Desc    bool   `query:"desc"`
		}

		// Bind the query parameters to the struct
		params, err := vayu.BindQueryParams[SearchParams](c)
		if err != nil {
			c.JSON(vayu.StatusBadRequest, map[string]string{
				"error": "Invalid search parameters: " + err.Error(),
			})
			return
		}

		// Use the typed parameters to send a response
		vayu.JSONResponse(c, vayu.StatusOK, map[string]interface{}{
			"message":    "Search completed successfully",
			"parameters": params,
			"results":    []string{"Result 1", "Result 2", "Result 3"},
		})
	})

	// Set up error handling demonstration routes
	setupErrorRoutes(app)

	// Print routes for testing
	log.Println("Available Routes:")
	log.Println("  GET  /hello                 - Basic text response")
	log.Println("  GET  /json                  - Basic JSON response")
	log.Println("  GET  /users/:id             - Route with parameter")
	log.Println("  GET  /api/status            - API status endpoint")
	log.Println("  GET  /generics/json         - Generic type-safe JSON response demo")
	log.Println("  GET  /generics/store        - Generic type-safe context store demo")
	log.Println(" POST  /generics/bind         - Generic type-safe JSON binding demo")
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
