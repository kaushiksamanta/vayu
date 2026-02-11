# Vayu Architecture

This document describes how Vayu is structured internally and how requests flow through the framework.

For API examples and day-to-day usage, see [USAGE.md](USAGE.md).

## Design Goals

- Keep the core framework minimal and composable.
- Use a straightforward middleware/handler execution model.
- Provide ergonomic APIs without hiding Go primitives.
- Add type-safe helpers with generics where they improve correctness.

## Core Building Blocks

### `App` (`vayu.go`)

`App` is the runtime container for:

- `router`: method-indexed route table (`Router`)
- `middleware`: ordered global middleware chain
- `NotFoundHandler`: configurable fallback handler

Key responsibilities:

- Route registration (`GET`, `POST`, etc.)
- Middleware registration (`Use`)
- HTTP server entrypoint (`ServeHTTP`)
- Static files (`Static`)
- Network listeners (`Listen`, `ListenTLS`)

### `Router` (`route.go`)

The router stores routes per method (`map[string][]route`) and performs linear matching:

- path is split by `/`
- static segments must match exactly
- `:param` segments are captured into `Context.Params`

### `Context` (`context.go`, `context_extension.go`)

`Context` wraps request/response state for each handler call:

- `Request`: underlying `*http.Request`
- `Writer`: wrapped `ResponseWriter`
- `Params`: route params extracted by router
- `Ctx`: request-scoped `context.Context`
- `store`: request-local key/value map

It also provides convenience helpers for query access, JSON send/bind, form files, HTML/text responses, and request-local storage.

### `ResponseWriter` (`response_writer.go`)

`ResponseWriter` wraps `http.ResponseWriter` and tracks:

- whether any response bytes/headers were written
- the status code written (implicit `200` on first `Write` without `WriteHeader`)

It guards against double `WriteHeader` calls and forwards optional interfaces (`http.Flusher`, `http.Hijacker`) from the underlying writer for streaming and WebSocket upgrade compatibility.

### Middleware and Error Pipeline (`middleware.go`, `logger.go`, `error.go`, `error_handler.go`)

- `Use` appends global middleware to the app chain.
- `WithMiddleware` composes route-local middleware + one handler.
- `Recovery` provides panic recovery with a default 500 response (skipped if response was already written).
- `ErrorHandlerMiddleware` supports custom panic-to-error translation.
- `SilentMode` suppresses panic logs in test contexts when enabled.

### Generic Type-Safe Helpers (`generic.go`)

Generic helpers add compile-time typing over common operations:

- `JSONResponse[T]`
- `BindJSONBody[T]`, `MustBindJSONBody[T]`
- `SetValue[T]`, `GetValue[T]`
- `BindQueryJSON[T]`, `MustBindQueryJSON[T]`
- `BindParamJSON[T]`, `MustBindParamJSON[T]`
- `BindQueryParams[T]`, `MustBindQueryParams[T]`

`BindQueryParams[T]` enforces that `T` is a struct, uses struct tags (`query:"..."`, optional `required:"true"`), and performs bit-size-aware integer parsing to catch overflows.

### Route Groups (`group.go`)

`Group` applies a path prefix and allows grouped routes/middleware to be attached with shared path space. Prefix matching uses exact-or-slash logic (`path == prefix || HasPrefix(path, prefix+"/")`) to prevent `/api` from matching `/api2`.

## Request Lifecycle

1. Incoming HTTP request hits `App.ServeHTTP`.
2. Framework creates a timeout-bound request context (`30s`) and Vayu `Context`.
3. Router matches method + path and extracts path params.
4. Global middleware slice is copied per-request (`make`+`copy`) and the route handler is appended.
5. `next()` advances through chain until completion or stop.
6. If context deadline is exceeded, framework writes `StatusGatewayTimeout`.
7. If no route matches, `NotFoundHandler` (or default 404) is returned.

## Repository Layout

```text
vayu/
├── vayu.go                 # App lifecycle, route registration, HTTP serving
├── route.go                # Router + path/param matching
├── group.go                # Route groups
├── context.go              # Request context core helpers
├── context_extension.go    # Extra context helpers and store access
├── response_writer.go      # Response writer wrapper/status tracking
├── response.go             # Response helper methods
├── status.go               # HTTP status constants
├── middleware.go           # Route-local middleware composition
├── logger.go               # Logger middleware
├── error.go                # Basic recovery middleware
├── error_handler.go        # Configurable panic/error middleware
├── generic.go              # Generic typed helpers (binding/store/response)
├── example/main.go         # Runnable sample app
├── tests/unit/             # Unit tests
├── tests/integration/      # Integration tests
├── docs/                   # Additional project docs
├── Makefile                # Build/test/lint tasks
└── README.md               # Lean project overview
```

## Testing Strategy

- Unit tests validate component behavior (`tests/unit`).
- Integration tests validate HTTP flow end-to-end (`tests/integration`).
- `make check` runs both tests and lint-equivalent checks.
- `make test-race` runs all tests with the Go race detector; also runs in CI.

## Notes on Tradeoffs

- Route matching is intentionally simple and predictable (linear scan).
- Middleware model favors clarity over highly abstract pipelines.
- Reflection in query binding is localized to optional generic helpers, while core path execution stays simple.
