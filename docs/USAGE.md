# Vayu Usage Guide

For project overview and quick start, see [README.md](../README.md).
For framework internals and design, see [ARCHITECTURE.md](ARCHITECTURE.md).

## Features

- Idiomatic Go package and API design
- Fluent API for route and middleware registration
- `context.Context` integration for cancellation and timeouts
- Type-safe generic helpers for JSON/query/context operations
- HTTP methods: `GET`, `POST`, `PUT`, `DELETE`, `PATCH`, `OPTIONS`, `HEAD`
- Built-in HTTP status code constants (`StatusOK`, `StatusNotFound`, etc.)
- Response helpers (`OK`, `Created`, `BadRequest`, etc.)
- Middleware chaining with `next()`
- Panic/error handling middleware
- Route groups for API modularization
- Static file serving
- JSON/form/file upload handling
- Request-local key/value store helpers
- TLS support via `ListenTLS`
- Comprehensive unit and integration tests

## Installation

```bash
go get github.com/kaushiksamanta/vayu
```

```go
import "github.com/kaushiksamanta/vayu"
```

## Quick Start

```go
package main

import (
	"github.com/kaushiksamanta/vayu"
)

func main() {
	app := vayu.New()

	app.Use(vayu.Logger()).Use(vayu.Recovery())

	app.GET("/", func(c *vayu.Context, next vayu.NextFunc) {
		c.Send(vayu.StatusOK, "Hello, Vayu!")
	})

	app.Listen(":8080")
}
```

```bash
go run main.go
```

Visit `http://localhost:8080`.

## Middleware Chaining

```go
app.Use(func(c *vayu.Context, next vayu.NextFunc) {
	fmt.Println("Before handler")
	next()
	fmt.Println("After handler")
})

app.GET("/hello", func(c *vayu.Context, next vayu.NextFunc) {
	c.Send(vayu.StatusOK, "Hello, Middleware!")
})
```

## Query Parameters

```go
app.GET("/search", func(c *vayu.Context, next vayu.NextFunc) {
	term := c.Query("term")
	if term == "" {
		c.JSON(vayu.StatusBadRequest, map[string]string{"error": "Missing 'term'"})
		return
	}
	c.JSON(vayu.StatusOK, map[string]string{"search": term})
})

// curl "http://localhost:8080/search?term=go"
```

## Static File Serving

```go
app.Static("/assets", "./public")

// /public/logo.png -> http://localhost:8080/assets/logo.png
```

## Route Groups

```go
api := app.Group("/api/v1")
api.GET("/users", func(c *vayu.Context, next vayu.NextFunc) {
	c.JSON(vayu.StatusOK, map[string]string{"message": "User list"})
})
```

Endpoint: `http://localhost:8080/api/v1/users`.

## Error Handling Middleware

```go
app.Use(vayu.Recovery())
```

For custom panic handling:

```go
customErrorHandler := func(c *vayu.Context, err error) {
	fmt.Printf("Error caught: %v\n", err)
	c.JSON(vayu.StatusInternalServerError, map[string]string{
		"error": err.Error(),
	})
}

app.Use(vayu.ErrorHandlerMiddleware(customErrorHandler))
```

## Type-Safe Generics

### Type-Safe JSON Responses

```go
type APIResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    []string `json:"data"`
}

resp := APIResponse{
	Status:  "success",
	Message: "Items retrieved",
	Data:    []string{"a", "b"},
}

_ = vayu.JSONResponse(c, vayu.StatusOK, resp)
```

### Type-Safe JSON Body Binding

```go
type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
}

userRequest, err := vayu.BindJSONBody[CreateUserRequest](c)
if err != nil {
	return
}

fmt.Println(userRequest.Username)

// If failure should panic:
// user := vayu.MustBindJSONBody[CreateUserRequest](c)
```

### Type-Safe Context Store

```go
vayu.SetValue(c, "user_id", 42)
vayu.SetValue(c, "is_admin", true)

type User struct {
	Name  string
	Email string
}

vayu.SetValue(c, "current_user", User{Name: "John", Email: "john@example.com"})

userID, ok := vayu.GetValue[int](c, "user_id")
isAdmin, ok := vayu.GetValue[bool](c, "is_admin")
currentUser, ok := vayu.GetValue[User](c, "current_user")
_, _ = userID, ok
_, _ = isAdmin, currentUser
```

### Type-Safe Query Parameter JSON Binding

```go
type Filter struct {
	Category string  `json:"category"`
	MinPrice float64 `json:"minPrice"`
	InStock  bool    `json:"inStock"`
}

// URL: /products?filter={"category":"books","minPrice":19.99,"inStock":true}
filter, err := vayu.BindQueryJSON[Filter](c, "filter")
if err != nil {
	return
}

fmt.Println(filter.Category)
```

### Type-Safe Path Parameter JSON Binding

```go
type Config struct {
	View       string `json:"view"`
	ShowPrices bool   `json:"showPrices"`
}

// Route: /products/:config
// URL-encoded JSON in :config path param
config, err := vayu.BindParamJSON[Config](c, "config")
if err != nil {
	return
}

fmt.Println(config.View)
```

### Type-Safe Query Parameters Binding

```go
type SearchParams struct {
	Term       string   `query:"q" required:"true"`
	Page       int      `query:"page"`
	PerPage    int      `query:"per_page"`
	Tags       []string `query:"tags"`
	SortBy     string   `query:"sort"`
	Descending bool     `query:"desc"`
}

params, err := vayu.BindQueryParams[SearchParams](c)
if err != nil {
	return
}

fmt.Println(params.Term, params.Tags)
```

## File Uploads

```go
app.POST("/upload", func(c *vayu.Context, next vayu.NextFunc) {
	file, header, err := c.FormFile("myfile")
	if err != nil {
		c.JSON(vayu.StatusBadRequest, map[string]string{"error": "File upload failed"})
		return
	}
	defer file.Close()

	dst, err := os.Create("./uploads/" + header.Filename)
	if err != nil {
		c.JSON(vayu.StatusInternalServerError, map[string]string{"error": "Could not create destination file"})
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		c.JSON(vayu.StatusInternalServerError, map[string]string{"error": "Could not save file"})
		return
	}

	c.JSON(vayu.StatusOK, map[string]string{"message": "File uploaded"})
})
```

## Custom Middleware

```go
auth := func(c *vayu.Context, next vayu.NextFunc) {
	if c.Request.Header.Get("X-Token") != "secret" {
		c.JSON(vayu.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}
	next()
}

app.GET("/secure", vayu.WithMiddleware(func(c *vayu.Context, next vayu.NextFunc) {
	c.Send(vayu.StatusOK, "Access granted")
}, auth))
```

## HTTP Status Codes

Vayu ships HTTP status constants, so importing `net/http` only for status codes is optional:

```go
c.Send(vayu.StatusOK, "Success")
c.JSON(vayu.StatusNotFound, map[string]string{"error": "Resource not found"})
```

## Response Helpers

```go
// Success responses
c.OK(map[string]string{"message": "Success"})
c.Created(map[string]string{"id": "123"})
c.NoContent()

// Error responses
c.BadRequest("Invalid parameters")
c.Unauthorized("Authentication required")
c.Forbidden("Access denied")
c.NotFound("Resource not found")
c.InternalServerError("Something went wrong")
```

## Development

```bash
# Build framework binary
make build

# Run example application
make run

# Remove generated artifacts (preserves bin/)
make clean
```

## Testing

Vayu includes unit and integration tests in `tests/unit` and `tests/integration`.
When running tests, `SilentMode` can be controlled with linker flags, and by default test binaries suppress panic recovery logs for cleaner output.

```bash
# Run all tests
make test

# Run only unit tests
make test-unit

# Run only integration tests
make test-integration

# Verbose test run
make test-v

# Coverage (terminal)
make test-cover

# Coverage (HTML report)
make test-cover-html

# Benchmarks
make bench

# Race detector
make test-race

# Panic logs enabled
make test-verbose-logs

# Panic logs suppressed
make test-silent
```

## Code Quality

```bash
# Format code
make fmt

# Go vet
make vet

# staticcheck
make staticcheck

# Full lint pipeline
make lint

# CI-equivalent local check
make check
```

## Contributing

1. Fork the repository.
2. Create your feature branch (`git checkout -b feature/amazing-feature`).
3. Commit your changes (`git commit -m 'feat: your change'`).
4. Push the branch (`git push origin feature/amazing-feature`).
5. Open a pull request.

## License

MIT License. See [LICENSE](LICENSE).
