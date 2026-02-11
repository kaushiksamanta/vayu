# Vayu

Vayu is a lightweight web framework for Go inspired by Express.js. The name "Vayu" (वायु in Sanskrit, বায়ু in Bangla) means "air" or "wind", reflecting a thin, minimal framework.

## Features

- Fluent API for routes and middleware
- Context-aware request handling (`context.Context`)
- Type-safe generic helpers for JSON/query/context operations
- Built-in HTTP status constants and response helpers
- Middleware chaining, route groups, and static file serving
- Unit and integration test coverage

## Documentation

- Usage guide: [docs/USAGE.md](docs/USAGE.md)
- Architecture and internals: [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)

## Quick Start

```bash
go get github.com/kaushiksamanta/vayu
```

```go
package main

import "github.com/kaushiksamanta/vayu"

func main() {
	app := vayu.New()
	app.GET("/", func(c *vayu.Context, next vayu.NextFunc) {
		c.Send(vayu.StatusOK, "Hello, Vayu!")
	})
	app.Listen(":8080")
}
```

```bash
go run main.go
```

## Development

```bash
make check
```

## Contributing

See the contributing section in [docs/USAGE.md](docs/USAGE.md#contributing).

## License

MIT License. See [LICENSE](LICENSE).
