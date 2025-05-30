# Vayu - A Lightweight Go Web Framework

<p align="center">
  <img src="vayu.jpeg" alt="Vayu Framework Mascot" width="300"/>
</p>

Vayu is a lightweight web framework for Go inspired by Express.js. The name "Vayu" is derived from the Sanskrit word for "air", symbolizing a thin and light framework. It focuses on simplicity, minimalism, and flexibility, allowing you to build web applications quickly while offering a variety of features like middleware support, routing, query parameter handling, file uploads, and more.

## ⚡ Features

- **Idiomatic Go Structure**: Follows Go best practices for package layout and API design
- **Fluent API**: Method chaining for route registration and middleware
- **Context Support**: Integrates with Go's `context.Context` for cancellation and timeouts
- **Type-Safe Generics**: Generic functions for type-safe JSON and context operations
- **Complete HTTP Methods**: Support for GET, POST, PUT, DELETE, PATCH, OPTIONS, HEAD
- **Built-in Status Codes**: No need to import net/http just for status codes
- **Response Helpers**: Convenient methods like `OK()`, `Created()`, `BadRequest()` for common responses
- **Middleware System**: Express-style middleware with next() function
- **Error Handling**: Robust error handling middleware with panic recovery
- **Route Groups**: For API versioning and modular routing
- **Static File Serving**: Serve assets from a directory
- **JSON Handling**: Parse and respond with JSON with type safety
- **Form Processing**: Handle form data and file uploads
- **Request Store**: Context-based key-value store with type-safe access
- **TLS Support**: HTTPS server capability
- **Comprehensive Testing**: Structured unit and integration tests

## 🚀 Installation

To use Vayu, install it using Go modules:

```bash
go get github.com/kaushiksamanta/vayu
```

Then import it in your Go project:

```go
import "github.com/kaushiksamanta/vayu"
```
## 💡 Basic Usage

### Creating a New App

```go
package main

import (
	"github.com/kaushiksamanta/vayu"
)

func main() {
	// Create a new Vayu app
	app := vayu.New()

	// Global middleware
	app.Use(vayu.Logger())
	   .Use(vayu.Recovery())

	// Simple GET route
	app.GET("/", func(c *vayu.Context, next vayu.NextFunc) {
		c.Send(vayu.StatusOK, "Hello, Vayu!")
	})

	// Start the app on port 8080
	app.Listen(":8080")
}
```

```bash
# Run the server
go run main.go
```

Visit http://localhost:8080 to see "Hello, Vayu!".

## 🌍 Advanced Features

### Middleware Chaining

Vayu supports middleware chaining similar to Express.js:

```go
app.Use(func(c *vayu.Context, next vayu.NextFunc) {
    fmt.Println("Before Handler")
    next() // Call next middleware or handler
    fmt.Println("After Handler")
})

app.GET("/hello", func(c *vayu.Context, next vayu.NextFunc) {
    c.Send(vayu.StatusOK, "Hello, Middleware!")
})

// The output when accessing the /hello route would be:
// Before Handler
// After Handler
```

### Query Parameters

Easily access query parameters with `c.Query("param")`:

```go
app.GET("/search", func(c *vayu.Context, next vayu.NextFunc) {
	term := c.Query("term")
	if term == "" {
		c.JSON(vayu.StatusBadRequest, map[string]string{"error": "Missing 'term'"})
		return
	}
	c.JSON(vayu.StatusOK, map[string]string{"search": term})
})

// Test with: curl "http://localhost:8080/search?term=go"
// Response: {"search":"go"}
```

### Static File Serving
Serve static files (e.g., images, CSS, JavaScript) from a directory:

```go
app.Static("/assets", "./public")

// This serves files from the `public` folder under the `/assets` route:
// /public/logo.png  →  http://localhost:8080/assets/logo.png
```
### Route Groups

Use route groups to organize API routes (e.g., for API versioning):

```go
api := app.Group("/api/v1")
api.GET("/users", func(c *vayu.Context, next vayu.NextFunc) {
    c.JSON(vayu.StatusOK, map[string]string{"message": "User List"})
})
```

This endpoint will be accessible at `http://localhost:8080/api/v1/users`.

### Error Handling Middleware

Catch panics globally and prevent server crashes:

```go
app.Use(vayu.Recovery())
```

This middleware catches panics and returns a 500 Internal Server Error response.

### Type-Safe Generics

Vayu provides generic functions for improved type safety using Go's generics feature. These functions help catch type errors at compile time rather than at runtime:

#### Type-Safe JSON Responses

```go
// Define your response type
type ApiResponse struct {
    Status  string `json:"status"`
    Message string `json:"message"`
    Data    []Item `json:"data"`
}

// Create a typed response
response := ApiResponse{
    Status:  "success",
    Message: "Items retrieved successfully",
    Data:    items,
}

// Send with type safety
err := vayu.JSONResponse(c, vayu.StatusOK, response)
```

#### Type-Safe JSON Binding

```go
// Define your request type
type CreateUserRequest struct {
    Username string `json:"username"`
    Email    string `json:"email"`
    Age      int    `json:"age"`
}

// Bind with type safety
userRequest, err := vayu.BindJSONBody[CreateUserRequest](c)
if err != nil {
    // Handle error
    return
}

// Use the strongly typed fields
fmt.Println(userRequest.Username) // Type-safe field access

// For cases when you know binding will succeed
user := vayu.MustBindJSONBody[UserProfile](c) // Panics on error
```

#### Type-Safe Context Store

```go
// Store values with type safety
vayu.SetValue(c, "user_id", 42)
vayu.SetValue(c, "is_admin", true)

// Define a custom type
type User struct {
    Name  string
    Email string
}

// Store a struct
user := User{Name: "John", Email: "john@example.com"}
vayu.SetValue(c, "current_user", user)

// Retrieve values with type safety
userID, ok := vayu.GetValue[int](c, "user_id")         // 42, true
isAdmin, ok := vayu.GetValue[bool](c, "is_admin")     // true, true
currentUser, ok := vayu.GetValue[User](c, "current_user") // User{...}, true

// Type mismatches are caught
str, ok := vayu.GetValue[string](c, "user_id") // "", false (incorrect type)
```

#### Type-Safe Query Parameter JSON Binding

Bind JSON data from query parameters with type safety:

```go
// Define your type
type Filter struct {
    Category string  `json:"category"`
    MinPrice float64 `json:"minPrice"`
    InStock  bool    `json:"inStock"`
}

// URL: /products?filter={"category":"books","minPrice":19.99,"inStock":true}
filter, err := vayu.BindQueryJSON[Filter](c, "filter")
if err != nil {
    // Handle error
    return
}

// Access type-safe fields
fmt.Println(filter.Category)  // "books"
fmt.Println(filter.MinPrice)  // 19.99
fmt.Println(filter.InStock)   // true

// For cases when binding will always succeed
filter := vayu.MustBindQueryJSON[Filter](c, "filter") // Panics on error
```

#### Type-Safe Path Parameter JSON Binding

Bind JSON data from path parameters with type safety:

```go
// Define your type
type Config struct {
    View       string `json:"view"`
    ShowPrices bool   `json:"showPrices"`
}

// Route: /products/:config
// URL: /products/%7B%22view%22%3A%22grid%22%2C%22showPrices%22%3Atrue%7D (URL-encoded JSON)
config, err := vayu.BindParamJSON[Config](c, "config")
if err != nil {
    // Handle error
    return
}

// Access type-safe fields
fmt.Println(config.View)        // "grid"
fmt.Println(config.ShowPrices)  // true
```

#### Type-Safe Query Parameters Binding

Bind multiple query parameters to a struct using tags:

```go
// Define struct with query tags
type SearchParams struct {
    Term     string   `query:"q" required:"true"`  // required:"true" marks as mandatory
    Page     int      `query:"page"`
    PerPage  int      `query:"per_page"`
    Tags     []string `query:"tags"`              // Will parse comma-separated values
    SortBy   string   `query:"sort"`
    Descending bool    `query:"desc"`
}

// URL: /search?q=golang&page=2&per_page=20&tags=web,api&sort=relevance&desc=true
params, err := vayu.BindQueryParams[SearchParams](c)
if err != nil {
    // Handle error (including required param missing)
    return
}

// Access type-safe fields
fmt.Println(params.Term)      // "golang"
fmt.Println(params.Page)      // 2
fmt.Println(params.Tags)      // ["web", "api"]
fmt.Println(params.Descending) // true
```

### File Uploads

Handle file uploads via `multipart/form-data`:

```go
app.POST("/upload", func(c *vayu.Context, next vayu.NextFunc) {
	file, header, err := c.FormFile("myfile")
	if err != nil {
		c.JSON(vayu.StatusBadRequest, map[string]string{"error": "File upload failed"})
		return
	}
	defer file.Close()

	dst, _ := os.Create("./uploads/" + header.Filename)
	defer dst.Close()
	io.Copy(dst, file)

	c.JSON(vayu.StatusOK, map[string]string{"message": "File uploaded!"})
})
```

This route accepts file uploads, stores them in the `uploads/` directory, and responds with a success message.

### Custom Middleware

Create and use custom middleware functions:

```go
// Define authentication middleware
auth := func(c *vayu.Context, next vayu.NextFunc) {
    if c.Request.Header.Get("X-Token") != "secret" {
        c.JSON(vayu.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
        return
    }
    next() // Continue to handler if authorized
}

// Apply middleware to specific route
app.GET("/secure", vayu.WithMiddleware(func(c *vayu.Context, next vayu.NextFunc) {
    c.Send(vayu.StatusOK, "Access granted")
}, auth))
```

This example protects the `/secure` route with token-based authentication.

## Project Structure

```
vayu/
├── context.go           # Request context implementation
├── context_extension.go # Additional context methods
├── error_handler.go     # Error handling middleware with SilentMode
├── group.go             # Route group implementation
├── logger.go            # Logging middleware
├── middleware.go        # Middleware utilities
├── response.go          # Response helper methods
├── response_writer.go   # Custom ResponseWriter implementation
├── route.go             # Router implementation
├── status.go            # HTTP status code constants
├── vayu.go              # Core application code
├── Makefile             # Build/test automation
├── .gitignore           # Git ignore file
├── go.mod               # Go module definition
├── README.md            # Project documentation
├── vayu.jpeg            # Framework mascot image
├── .github/             # GitHub configurations
│   └── workflows/       # GitHub Actions workflows
├── example/             # Example applications
│   └── main.go          # Simple demo app
└── tests/               # Test suite
    ├── unit/            # Unit tests
    │   ├── context_test.go         # Context tests
    │   ├── log_test.go             # SilentMode tests
    │   ├── middleware_test.go      # Middleware tests
    │   ├── response_helpers_test.go # Response helper tests
    │   └── router_test.go          # Router tests
    └── integration/     # Integration tests
        └── api_test.go  # API endpoint tests
```

## Development

Vayu uses standard Go tooling. Common tasks are available through the Makefile:

```bash
# Build the project
make build

# Run the example application
make run

# Clean build artifacts
make clean
```

### Testing

Vayu has a comprehensive test suite organized in a `tests/` directory with both unit and integration tests. The framework automatically detects test environments and enables `SilentMode`, which suppresses panic recovery logs during test execution, resulting in cleaner test output.

```bash
# Run all tests
make test

# Run only unit tests
make test-unit

# Run only integration tests
make test-integration

# Run all tests with verbose output
make test-v

# Generate test coverage report
make test-cover-html

# Run benchmarks
make bench

# Run tests with race detection
make test-race

# Run tests with panic logs explicitly enabled (override SilentMode)
make test-verbose-logs

# Run tests with panic logs explicitly disabled (force SilentMode)
make test-silent
```

### Code Quality

Vayu maintains high code quality standards with these commands:

```bash
# Format code
make fmt

# Run Go's static analysis (vet)
make vet

# Run staticcheck advanced linter
make staticcheck

# Run all linters
make lint

# Run tests and linting
make check
```

## HTTP Status Codes

Vayu provides all standard HTTP status codes as constants, so you don't need to import `net/http` just for status codes:

```go
// Use Vayu's built-in status codes
c.Send(vayu.StatusOK, "Success!")
c.JSON(vayu.StatusNotFound, map[string]string{"error": "Resource not found"})
```

Available status codes include `StatusOK` (200), `StatusCreated` (201), `StatusBadRequest` (400), `StatusNotFound` (404), and all other standard HTTP status codes.

## Response Helpers

Vayu provides convenient response helper methods to simplify common response patterns:

```go
// Success responses
c.OK(map[string]string{"message": "Success"})        // 200 OK
c.Created(map[string]string{"id": "123"})           // 201 Created
c.NoContent()                                        // 204 No Content

// Error responses
c.BadRequest("Invalid parameters")                   // 400 Bad Request
c.Unauthorized("Authentication required")           // 401 Unauthorized
c.Forbidden("Access denied")                        // 403 Forbidden
c.NotFound("Resource not found")                    // 404 Not Found
c.InternalServerError("Something went wrong")       // 500 Internal Server Error
```

## Error Handling

Vayu includes robust error handling middleware that can catch and process errors and panics. The framework provides a `SilentMode` flag that automatically detects test environments to suppress panic recovery logs during tests, resulting in cleaner output while still maintaining full error handling capabilities:

```go
// Custom error handler
customErrorHandler := func(c *vayu.Context, err error) {
    // Log the error
    fmt.Printf("Error caught: %v\n", err)
    
    // Send appropriate response
    c.JSON(vayu.StatusInternalServerError, map[string]string{
        "error": err.Error(),
    })
}

// Add error handling middleware
app.Use(vayu.ErrorHandlerMiddleware(customErrorHandler))

// You can also manually control silent mode if needed
// vayu.SilentMode = true // Suppress panic logs
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

MIT License
