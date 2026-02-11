# Repository Guidelines

## Project Structure & Module Organization
- Core framework code lives at the repository root (`vayu.go`, `context.go`, `route.go`, `response.go`, etc.) and is published as `github.com/kaushiksamanta/vayu`.
- `example/main.go` is the runnable reference app for local development and manual API checks.
- Automated tests are split by scope: `tests/unit` for package-level behavior and `tests/integration` for end-to-end HTTP flows.
- CI config is in `.github/workflows/go.yml`; lint settings are in `.golangci.yml`.
- `internal/` is currently reserved (placeholder), and `bin/` contains local tooling binaries (do **not** delete `bin/` — use `make clean` which only removes generated artifacts).
- Documentation lives in `docs/` (`USAGE.md`, `ARCHITECTURE.md`).

## Build, Test, and Development Commands
- `make build`: compile the framework into `./vayu`.
- `make run`: run the sample app from `example/main.go`.
- `make test`: run all tests (`go test ./...`).
- `make test-unit` / `make test-integration`: run only one test tier.
- `make test-race`: run all tests with the Go race detector (`go test -race ./...`). Also runs in CI.
- `make test-cover` or `make test-cover-html`: generate coverage metrics/report.
- `make lint`: run formatting, vet, and static analysis.
- `make check`: CI-equivalent local gate (`test` + `lint`).
- `make clean`: remove generated artifacts (`vayu` binary, `coverage.out`, `coverage.html`). Does **not** remove `bin/`.

## Coding Style & Naming Conventions
- Use Go `1.24.x` conventions from `go.mod`.
- Format with `make fmt` (`gofmt -w -s .`) before committing.
- Follow idiomatic Go naming: exported identifiers in `PascalCase`, internal helpers in `camelCase`, and `_test.go` suffix for tests.
- Prefer small, composable middleware/handler functions and existing response helpers (`OK`, `BadRequest`, `JSONResponse`, etc.).

## Key API & Design Notes
- **Middleware chain** is copied per-request (`make`+`copy`) to avoid race conditions from shared slice backing arrays.
- **Group/Static prefix matching** uses exact-or-slash logic (`path == prefix || HasPrefix(path, prefix+"/")`) to prevent `/api` from matching `/api2`.
- **`Context.WithTimeout`** returns a `context.CancelFunc`; callers must `defer cancel()`.
- **`Recovery()` middleware** checks `c.Writer.Written()` before writing a 500 response to avoid double-write corruption.
- **`ResponseWriter`** tracks implicit 200 status on `Write`, guards against double `WriteHeader`, and forwards `http.Flusher` and `http.Hijacker` interfaces.
- **`BindQueryParams[T]`** enforces `T` is a struct and uses bit-size-aware integer parsing to catch overflows.

## Testing Guidelines
- Primary framework: Go's built-in `testing` package; `testify` is available when richer assertions are helpful.
- Keep unit tests in `tests/unit` (`package unit`) and integration tests in `tests/integration` (`package integration`).
- Name tests `TestXxx` with clear behavior focus (example: `TestRouterMethods`).
- No enforced coverage threshold exists; maintain or improve coverage for changed code and include regression tests for bug fixes.
- Run `make test-race` before submitting changes that touch middleware, routing, or shared state.

## CI Pipeline
- CI is defined in `.github/workflows/go.yml`.
- `golangci-lint` is pinned to `v1.64.8` via `golangci/golangci-lint-action@v6` for deterministic builds.
- CI steps: **Build → Test → Test Race → Lint**.

## Commit & Pull Request Guidelines
- Prefer Conventional Commit-style subjects used in history: `feat:`, `fix:`, `docs:`, `chore:`, `refactor:`, `style:`, `ci:`.
- Keep commits scoped to one logical change and use imperative summaries.
- PRs should include: problem/solution summary, linked issue (if any), and exact verification commands run (for example, `make check`).
- For API behavior changes, include example request/response snippets or routes exercised from `example/main.go`.
